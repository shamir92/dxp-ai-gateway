package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/integration/newapi"
	"github.com/dxpglobal/ai-gateway/backend/internal/repository"
)

var (
	ErrUnsupportedSite        = errors.New("unsupported external site")
	ErrExternalAccountMissing = errors.New("external account is not provisioned")
	ErrExternalCredentials    = errors.New("external credentials could not be used")
	ErrExternalRejected       = errors.New("external platform rejected the request")
	ErrExternalUnavailable    = errors.New("external platform is unavailable")
)

const (
	// accessTokenName labels the gateway-managed access token on the external site.
	accessTokenName = "dxp-gateway"
	// accessTokenTTLDays matches the upstream default lifetime.
	accessTokenTTLDays = 30
	// maxKeyScanSize bounds the listing used to locate a freshly created key.
	maxKeyScanSize = 100
	// maxReprovisionAttempts caps how often an unusable remote account is
	// recreated before we stop and let an operator look at it.
	maxReprovisionAttempts = 5
)

// accessTokenScopes is the least privilege the gateway needs: manage its own API
// keys (read/write, plus a one-time reveal on create) and read subscription usage.
var accessTokenScopes = []string{
	newapi.ScopeAPIKeyRead,
	newapi.ScopeAPIKeyWrite,
	newapi.ScopeAPIKeyReveal,
	newapi.ScopeWalletRead,
}

// SiteClient is the external new-api surface this usecase depends on.
type SiteClient interface {
	Login(ctx context.Context, username string, password []byte) (newapi.SessionToken, error)
	Verify(ctx context.Context, bearer, scope string, proofContext any, password []byte) (newapi.SecurityProof, error)
	CreateAccessToken(ctx context.Context, bearer, proofToken string, request newapi.AccessTokenCreateRequest) (newapi.AccessTokenCreated, error)
	ListTokens(ctx context.Context, bearer string, page, pageSize int) (newapi.TokenPage, error)
	CreateToken(ctx context.Context, bearer string, input newapi.TokenInput) error
	UpdateToken(ctx context.Context, bearer string, input newapi.TokenInput) (newapi.Token, error)
	DeleteToken(ctx context.Context, bearer string, id int) error
	RevealTokenKey(ctx context.Context, bearer string, id int) (string, error)
	SubscriptionSelf(ctx context.Context, bearer string) (newapi.SubscriptionSelf, error)
	// Status is public: it supplies the platform's advertised API addresses.
	Status(ctx context.Context) (newapi.Status, error)
}

// SiteCredentialCipher decrypts the provisioned password and seals the minted
// access token at rest.
type SiteCredentialCipher interface {
	Decrypt(userID uuid.UUID, site entity.ExternalSite, ciphertext, nonce []byte, keyVersion string) ([]byte, error)
	EncryptAccessToken(userID uuid.UUID, site entity.ExternalSite, plaintext []byte) (ciphertext, nonce []byte, keyVersion string, err error)
	DecryptAccessToken(userID uuid.UUID, site entity.ExternalSite, ciphertext, nonce []byte, keyVersion string) ([]byte, error)
}

type APIKeyUsecase struct {
	users    repository.UserRepo
	accounts repository.ExternalSiteAccountRepo
	tokens   repository.AccessTokenRepo
	cipher   SiteCredentialCipher
	factory  ExternalAccountFactory
	clients  map[entity.ExternalSite]SiteClient
	now      func() time.Time
}

func NewAPIKeyUsecase(
	users repository.UserRepo,
	accounts repository.ExternalSiteAccountRepo,
	tokens repository.AccessTokenRepo,
	cipher SiteCredentialCipher,
	factory ExternalAccountFactory,
	clients map[entity.ExternalSite]SiteClient,
) *APIKeyUsecase {
	return &APIKeyUsecase{
		users:    users,
		accounts: accounts,
		tokens:   tokens,
		cipher:   cipher,
		factory:  factory,
		clients:  clients,
		now:      func() time.Time { return time.Now().UTC() },
	}
}

// ParseSite validates the `site` selector supplied by the caller.
func ParseSite(value string) (entity.ExternalSite, error) {
	site := entity.ExternalSite(strings.TrimSpace(strings.ToLower(value)))
	if site != entity.ExternalSiteTPAI && site != entity.ExternalSiteAPAI {
		return "", ErrUnsupportedSite
	}
	return site, nil
}

func (u *APIKeyUsecase) clientFor(site entity.ExternalSite) (SiteClient, error) {
	client := u.clients[site]
	if client == nil {
		return nil, ErrUnsupportedSite
	}
	return client, nil
}

// EnsureProvisioned creates the per-site account rows when they are missing, so
// the worker can provision a user who registered before the provisioning flow
// existed (or through a path that skipped it). Existing rows are never modified,
// so this is safe to call on every page load or profile save.
func (u *APIKeyUsecase) EnsureProvisioned(ctx context.Context, userID uuid.UUID) error {
	if u.users == nil || u.factory == nil {
		return nil
	}
	user, err := u.users.GetByID(ctx, userID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrExternalAccountMissing
	}
	if err != nil {
		return err
	}
	accounts, err := u.factory.Build(user)
	if err != nil {
		return err
	}
	return u.accounts.EnsurePending(ctx, accounts)
}

// WithBearer resolves a usable bearer for (user, site) and runs op, re-minting
// exactly once if the stored token was rejected. It never logs the bearer.
func withBearer[T any](
	u *APIKeyUsecase,
	ctx context.Context,
	userID uuid.UUID,
	site entity.ExternalSite,
	op func(bearer string) (T, error),
) (T, error) {
	var zero T

	bearer, err := u.ensureBearer(ctx, userID, site)
	if err != nil {
		return zero, err
	}
	result, err := op(bearer)
	if err == nil {
		return result, nil
	}

	var apiErr *newapi.APIError
	if !errors.As(err, &apiErr) || !apiErr.Unauthorized() {
		return zero, err
	}

	// The stored token is no longer accepted: discard it and mint once more.
	if clearErr := u.tokens.ClearAccessToken(ctx, userID, site); clearErr != nil {
		return zero, clearErr
	}
	bearer, err = u.ensureBearer(ctx, userID, site)
	if err != nil {
		return zero, err
	}
	return op(bearer)
}

// ensureBearer returns a usable access token, reusing the stored one and minting
// a replacement only when none is stored or it no longer decrypts.
func (u *APIKeyUsecase) ensureBearer(ctx context.Context, userID uuid.UUID, site entity.ExternalSite) (string, error) {
	if _, err := u.clientFor(site); err != nil {
		return "", err
	}

	stored, err := u.tokens.GetAccessToken(ctx, userID, site)
	switch {
	case err == nil:
		token, decErr := u.cipher.DecryptAccessToken(userID, site, stored.Ciphertext, stored.Nonce, stored.KeyVersion)
		if decErr == nil {
			defer clear(token)
			return string(token), nil
		}
		// Unreadable ciphertext: fall through and mint a fresh token.
	case errors.Is(err, repository.ErrNotFound):
		// No token yet: mint one.
	default:
		return "", err
	}

	return u.mintBearer(ctx, userID, site)
}

// mintBearer runs login → verify → access-token creation once and stores the
// result encrypted. A personal access token can never authorise this flow, so it
// must always start from the provisioned password.
func (u *APIKeyUsecase) mintBearer(ctx context.Context, userID uuid.UUID, site entity.ExternalSite) (string, error) {
	client, err := u.clientFor(site)
	if err != nil {
		return "", err
	}

	account, err := u.accounts.GetByUserAndSite(ctx, userID, site)
	if errors.Is(err, repository.ErrNotFound) {
		// No account row exists: enqueue both sites so the worker provisions them,
		// then report that setup is still in progress rather than failing outright.
		if ensureErr := u.EnsureProvisioned(ctx, userID); ensureErr != nil && !errors.Is(ensureErr, ErrExternalAccountMissing) {
			return "", ensureErr
		}
		return "", ErrExternalAccountMissing
	}
	if err != nil {
		return "", err
	}

	password, err := u.cipher.Decrypt(userID, site, account.PasswordCiphertext, account.PasswordNonce, account.EncryptionKeyVersion)
	if err != nil {
		return "", ErrExternalCredentials
	}
	defer clear(password)

	session, err := client.Login(ctx, account.Username, password)
	if err != nil {
		// A non-retryable rejection means the remote account is missing or was
		// changed even though our row says "provisioned" — new-api answers bad
		// credentials with HTTP 200 and success=false. Re-enqueue the row so the
		// worker recreates the account, and report "still provisioning" so the
		// caller retries instead of seeing a hard failure.
		var apiErr *newapi.APIError
		if errors.As(err, &apiErr) && !apiErr.Retryable {
			if requeued, requeueErr := u.accounts.Requeue(ctx, userID, site, maxReprovisionAttempts); requeueErr == nil && requeued {
				return "", ErrExternalAccountMissing
			}
		}
		return "", classifyMintError(err)
	}

	expiresAt := u.now().Add(accessTokenTTLDays * 24 * time.Hour).Unix()
	// The proof context and the create body must agree byte-for-byte or the
	// upstream answers SECURITY_PROOF_CONTEXT_MISMATCH.
	proofContext := newapi.AccessTokenGenerateContext{Scopes: accessTokenScopes, ExpiresAt: expiresAt}

	proof, err := client.Verify(ctx, session.AccessToken, newapi.ScopeAccessTokenGenerate, proofContext, password)
	if err != nil {
		return "", classifyMintError(err)
	}

	created, err := client.CreateAccessToken(ctx, session.AccessToken, proof.ProofToken, newapi.AccessTokenCreateRequest{
		Name:      accessTokenName,
		Scopes:    accessTokenScopes,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return "", classifyMintError(err)
	}

	ciphertext, nonce, keyVersion, err := u.cipher.EncryptAccessToken(userID, site, []byte(created.Token))
	if err != nil {
		return "", err
	}
	expires := time.Unix(expiresAt, 0).UTC()
	if err := u.tokens.SaveAccessToken(ctx, &entity.StoredAccessToken{
		UserID:     userID,
		Site:       site,
		Ciphertext: ciphertext,
		Nonce:      nonce,
		KeyVersion: keyVersion,
		ExpiresAt:  &expires,
		RemoteID:   strconv.Itoa(created.Item.ID),
		Status:     entity.AccessTokenActive,
	}); err != nil {
		return "", err
	}
	return created.Token, nil
}

// classifyMintError separates a permanent credential problem from a transient
// platform failure so callers can decide whether retrying is worthwhile.
func classifyMintError(err error) error {
	var apiErr *newapi.APIError
	switch {
	case errors.As(err, &apiErr) && apiErr.Unauthorized():
		return fmt.Errorf("%w: %s", ErrExternalCredentials, apiErr.Error())
	case errors.As(err, &apiErr) && apiErr.Retryable:
		return fmt.Errorf("%w: %s", ErrExternalUnavailable, apiErr.Error())
	case errors.As(err, &apiErr):
		return fmt.Errorf("%w: %s", ErrExternalRejected, apiErr.Error())
	default:
		return fmt.Errorf("%w: %v", ErrExternalUnavailable, err)
	}
}

// ListKeys returns one page of the user's API keys on one site.
func (u *APIKeyUsecase) ListKeys(ctx context.Context, userID uuid.UUID, site entity.ExternalSite, page, pageSize int) (newapi.TokenPage, error) {
	return withBearer(u, ctx, userID, site, func(bearer string) (newapi.TokenPage, error) {
		client, err := u.clientFor(site)
		if err != nil {
			return newapi.TokenPage{}, err
		}
		return client.ListTokens(ctx, bearer, page, pageSize)
	})
}

// CreatedKey carries the new key's metadata plus its plaintext value, which the
// upstream exposes exactly once.
type CreatedKey struct {
	Token newapi.Token
	Key   string
}

// CreateKey creates an API key and immediately reveals its value. The upstream
// does not return the created key, so the new row is located by diffing ids.
func (u *APIKeyUsecase) CreateKey(ctx context.Context, userID uuid.UUID, site entity.ExternalSite, input newapi.TokenInput) (CreatedKey, error) {
	return withBearer(u, ctx, userID, site, func(bearer string) (CreatedKey, error) {
		client, err := u.clientFor(site)
		if err != nil {
			return CreatedKey{}, err
		}

		existing, err := existingKeyIDs(ctx, client, bearer)
		if err != nil {
			return CreatedKey{}, err
		}
		if err := client.CreateToken(ctx, bearer, input); err != nil {
			return CreatedKey{}, err
		}

		page, err := client.ListTokens(ctx, bearer, 1, maxKeyScanSize)
		if err != nil {
			return CreatedKey{}, err
		}
		created, ok := newestKey(page.Items, existing)
		if !ok {
			// Created but not locatable: report success without the plaintext.
			return CreatedKey{}, nil
		}
		key, err := client.RevealTokenKey(ctx, bearer, created.ID)
		if err != nil {
			return CreatedKey{}, err
		}
		created.Key = ""
		return CreatedKey{Token: created, Key: key}, nil
	})
}

func existingKeyIDs(ctx context.Context, client SiteClient, bearer string) (map[int]struct{}, error) {
	page, err := client.ListTokens(ctx, bearer, 1, maxKeyScanSize)
	if err != nil {
		return nil, err
	}
	ids := make(map[int]struct{}, len(page.Items))
	for _, item := range page.Items {
		ids[item.ID] = struct{}{}
	}
	return ids, nil
}

func newestKey(items []newapi.Token, existing map[int]struct{}) (newapi.Token, bool) {
	var (
		found   newapi.Token
		hasFind bool
	)
	for _, item := range items {
		if _, seen := existing[item.ID]; seen {
			continue
		}
		if !hasFind || item.ID > found.ID {
			found, hasFind = item, true
		}
	}
	return found, hasFind
}

// UpdateKey updates an existing API key.
func (u *APIKeyUsecase) UpdateKey(ctx context.Context, userID uuid.UUID, site entity.ExternalSite, input newapi.TokenInput) (newapi.Token, error) {
	return withBearer(u, ctx, userID, site, func(bearer string) (newapi.Token, error) {
		client, err := u.clientFor(site)
		if err != nil {
			return newapi.Token{}, err
		}
		return client.UpdateToken(ctx, bearer, input)
	})
}

// DeleteKey removes an API key.
func (u *APIKeyUsecase) DeleteKey(ctx context.Context, userID uuid.UUID, site entity.ExternalSite, id int) error {
	_, err := withBearer(u, ctx, userID, site, func(bearer string) (struct{}, error) {
		client, err := u.clientFor(site)
		if err != nil {
			return struct{}{}, err
		}
		return struct{}{}, client.DeleteToken(ctx, bearer, id)
	})
	return err
}
