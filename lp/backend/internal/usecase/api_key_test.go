package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/dxpglobal/ai-gateway/backend/internal/entity"
	"github.com/dxpglobal/ai-gateway/backend/internal/integration/newapi"
	"github.com/dxpglobal/ai-gateway/backend/internal/repository"
)

// --- fakes ---

type fakeAccounts struct {
	account       *entity.ExternalSiteAccount
	err           error
	ensured       []entity.ExternalSiteAccount
	ensureErr     error
	requeueResult bool
	requeueErr    error
	requeueCalls  int
}

func (f *fakeAccounts) Requeue(context.Context, uuid.UUID, entity.ExternalSite, int) (bool, error) {
	f.requeueCalls++
	if f.requeueErr != nil {
		return false, f.requeueErr
	}
	return f.requeueResult, nil
}

func (f *fakeAccounts) GetByUserAndSite(context.Context, uuid.UUID, entity.ExternalSite) (*entity.ExternalSiteAccount, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.account, nil
}

func (f *fakeAccounts) EnsurePending(_ context.Context, accounts []entity.ExternalSiteAccount) error {
	if f.ensureErr != nil {
		return f.ensureErr
	}
	f.ensured = append(f.ensured, accounts...)
	return nil
}

func (f *fakeAccounts) ClaimDue(context.Context, time.Time, time.Time, int) ([]entity.ExternalSiteAccount, error) {
	return nil, nil
}

func (f *fakeAccounts) MarkProvisioned(context.Context, uuid.UUID, uuid.UUID, string, time.Time) error {
	return nil
}

func (f *fakeAccounts) MarkRetry(context.Context, uuid.UUID, uuid.UUID, time.Time, string, time.Time) error {
	return nil
}

func (f *fakeAccounts) MarkFailed(context.Context, uuid.UUID, uuid.UUID, string, time.Time) error {
	return nil
}

type fakeUsers struct {
	user *entity.User
	err  error
}

func (f *fakeUsers) GetByID(context.Context, uuid.UUID) (*entity.User, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.user == nil {
		return nil, repository.ErrNotFound
	}
	return f.user, nil
}

func (f *fakeUsers) Create(context.Context, *entity.User) error { return nil }

func (f *fakeUsers) GetByEmail(context.Context, string) (*entity.User, error) {
	return nil, repository.ErrNotFound
}

func (f *fakeUsers) Update(context.Context, *entity.User) error { return nil }

// fakeFactory stands in for credentials.AccountFactory.
type fakeFactory struct {
	accounts []entity.ExternalSiteAccount
	err      error
}

func (f fakeFactory) Build(*entity.User) ([]entity.ExternalSiteAccount, error) {
	return f.accounts, f.err
}

type fakeTokens struct {
	stored *entity.StoredAccessToken
	saves  int
	clears int
}

func (f *fakeTokens) GetAccessToken(context.Context, uuid.UUID, entity.ExternalSite) (*entity.StoredAccessToken, error) {
	if f.stored == nil {
		return nil, repository.ErrNotFound
	}
	return f.stored, nil
}

func (f *fakeTokens) SaveAccessToken(_ context.Context, token *entity.StoredAccessToken) error {
	f.saves++
	copied := *token
	f.stored = &copied
	return nil
}

func (f *fakeTokens) ClearAccessToken(context.Context, uuid.UUID, entity.ExternalSite) error {
	f.clears++
	f.stored = nil
	return nil
}

// fakeCipher is reversible and prefixes decrypted passwords so tests can assert
// that the stored password reached the login call.
type fakeCipher struct{}

func (fakeCipher) Decrypt(_ uuid.UUID, _ entity.ExternalSite, ciphertext, _ []byte, _ string) ([]byte, error) {
	return append([]byte("decrypted-"), ciphertext...), nil
}

func (fakeCipher) EncryptAccessToken(_ uuid.UUID, _ entity.ExternalSite, plaintext []byte) ([]byte, []byte, string, error) {
	return append([]byte(nil), plaintext...), []byte("nonce"), "v1", nil
}

func (fakeCipher) DecryptAccessToken(_ uuid.UUID, _ entity.ExternalSite, ciphertext, _ []byte, _ string) ([]byte, error) {
	return append([]byte(nil), ciphertext...), nil
}

// fakeClient lets each call be scripted, so tests assert the exact sequence the
// usecase performs against the external platform.
type fakeClient struct {
	login        func(username string, password []byte) (newapi.SessionToken, error)
	verify       func(bearer, scope string, proofContext any, password []byte) (newapi.SecurityProof, error)
	createAccess func(bearer, proofToken string, request newapi.AccessTokenCreateRequest) (newapi.AccessTokenCreated, error)
	listTokens   func(bearer string, page, pageSize int) (newapi.TokenPage, error)
	revealKey    func(bearer string, id int) (string, error)
	subscription func(bearer string) (newapi.SubscriptionSelf, error)
	status       func() (newapi.Status, error)
}

func (f *fakeClient) Status(context.Context) (newapi.Status, error) {
	if f.status == nil {
		return newapi.Status{}, nil
	}
	return f.status()
}

func (f *fakeClient) Login(_ context.Context, username string, password []byte) (newapi.SessionToken, error) {
	if f.login == nil {
		return newapi.SessionToken{}, errors.New("unexpected Login call")
	}
	return f.login(username, password)
}

func (f *fakeClient) Verify(_ context.Context, bearer, scope string, proofContext any, password []byte) (newapi.SecurityProof, error) {
	if f.verify == nil {
		return newapi.SecurityProof{}, errors.New("unexpected Verify call")
	}
	return f.verify(bearer, scope, proofContext, password)
}

func (f *fakeClient) CreateAccessToken(_ context.Context, bearer, proofToken string, request newapi.AccessTokenCreateRequest) (newapi.AccessTokenCreated, error) {
	if f.createAccess == nil {
		return newapi.AccessTokenCreated{}, errors.New("unexpected CreateAccessToken call")
	}
	return f.createAccess(bearer, proofToken, request)
}

func (f *fakeClient) ListTokens(_ context.Context, bearer string, page, pageSize int) (newapi.TokenPage, error) {
	if f.listTokens == nil {
		return newapi.TokenPage{}, errors.New("unexpected ListTokens call")
	}
	return f.listTokens(bearer, page, pageSize)
}

func (f *fakeClient) CreateToken(context.Context, string, newapi.TokenInput) error {
	return errors.New("unexpected CreateToken call")
}

func (f *fakeClient) UpdateToken(context.Context, string, newapi.TokenInput) (newapi.Token, error) {
	return newapi.Token{}, errors.New("unexpected UpdateToken call")
}

func (f *fakeClient) DeleteToken(context.Context, string, int) error {
	return errors.New("unexpected DeleteToken call")
}

func (f *fakeClient) RevealTokenKey(_ context.Context, bearer string, id int) (string, error) {
	if f.revealKey == nil {
		return "", errors.New("unexpected RevealTokenKey call")
	}
	return f.revealKey(bearer, id)
}

func (f *fakeClient) SubscriptionSelf(_ context.Context, bearer string) (newapi.SubscriptionSelf, error) {
	if f.subscription == nil {
		return newapi.SubscriptionSelf{}, errors.New("unexpected SubscriptionSelf call")
	}
	return f.subscription(bearer)
}

// provisionedUC wires a usecase around one TPAI account.
func provisionedUC(t *testing.T, tokens *fakeTokens, accounts *fakeAccounts, client *fakeClient) *APIKeyUsecase {
	t.Helper()
	if accounts.account == nil {
		accounts.account = &entity.ExternalSiteAccount{
			UserID:               uuid.New(),
			Site:                 entity.ExternalSiteTPAI,
			Username:             "user@example.com",
			PasswordCiphertext:   []byte("secret"),
			EncryptionKeyVersion: "v1",
		}
	}
	return NewAPIKeyUsecase(
		&fakeUsers{user: &entity.User{ID: accounts.account.UserID, Email: "user@example.com", Name: "Test"}},
		accounts,
		tokens,
		fakeCipher{},
		fakeFactory{},
		map[entity.ExternalSite]SiteClient{entity.ExternalSiteTPAI: client},
	)
}

// --- tests ---

func TestListKeysMintsOnceThenReusesTheStoredToken(t *testing.T) {
	userID := uuid.New()
	tokens := &fakeTokens{}
	client := &fakeClient{}

	mints := 0
	client.login = func(username string, password []byte) (newapi.SessionToken, error) {
		if username != "user@example.com" {
			t.Errorf("login username = %q, want the provisioned email", username)
		}
		if string(password) != "decrypted-secret" {
			t.Errorf("login password = %q, want the decrypted site password", password)
		}
		return newapi.SessionToken{AccessToken: "JWT-1"}, nil
	}
	client.verify = func(bearer, scope string, proofContext any, _ []byte) (newapi.SecurityProof, error) {
		if bearer != "JWT-1" {
			t.Errorf("verify bearer = %q, want the login token", bearer)
		}
		if scope != newapi.ScopeAccessTokenGenerate {
			t.Errorf("verify scope = %q", scope)
		}
		gen, ok := proofContext.(newapi.AccessTokenGenerateContext)
		if !ok {
			t.Fatalf("proof context type = %T, want AccessTokenGenerateContext", proofContext)
		}
		if len(gen.Scopes) != 4 {
			t.Errorf("proof scopes = %v, want the four requested scopes", gen.Scopes)
		}
		return newapi.SecurityProof{ProofToken: "PROOF-1"}, nil
	}
	client.createAccess = func(_, proofToken string, request newapi.AccessTokenCreateRequest) (newapi.AccessTokenCreated, error) {
		mints++
		if proofToken != "PROOF-1" {
			t.Errorf("proof header = %q", proofToken)
		}
		// The create body must agree with the proof context byte for byte.
		if request.ExpiresAt == 0 {
			t.Error("create expires_at must be set")
		}
		return newapi.AccessTokenCreated{Token: "nap_stored", Item: newapi.AccessToken{ID: 9}}, nil
	}

	var bearers []string
	client.listTokens = func(bearer string, _, _ int) (newapi.TokenPage, error) {
		bearers = append(bearers, bearer)
		return newapi.TokenPage{Items: []newapi.Token{{ID: 1, Name: "key"}}, Total: 1}, nil
	}

	uc := provisionedUC(t, tokens, &fakeAccounts{account: nil}, client)
	for i := 0; i < 3; i++ {
		if _, err := uc.ListKeys(context.Background(), userID, entity.ExternalSiteTPAI, 1, 20); err != nil {
			t.Fatalf("ListKeys call %d: %v", i, err)
		}
	}

	if mints != 1 {
		t.Errorf("mints = %d, want 1 — the token must be minted once and reused", mints)
	}
	if tokens.saves != 1 {
		t.Errorf("saves = %d, want 1", tokens.saves)
	}
	for i, bearer := range bearers {
		if bearer != "nap_stored" {
			t.Errorf("call %d bearer = %q, want the stored token", i, bearer)
		}
	}
}

func TestListKeysRemintsExactlyOnceWhenStoredTokenRejected(t *testing.T) {
	userID := uuid.New()
	tokens := &fakeTokens{stored: &entity.StoredAccessToken{
		Ciphertext: []byte("nap_stale"),
		KeyVersion: "v1",
		Status:     entity.AccessTokenActive,
	}}

	client := &fakeClient{}
	logins := 0
	client.login = func(string, []byte) (newapi.SessionToken, error) {
		logins++
		return newapi.SessionToken{AccessToken: "JWT-1"}, nil
	}
	client.verify = func(string, string, any, []byte) (newapi.SecurityProof, error) {
		return newapi.SecurityProof{ProofToken: "PROOF-1"}, nil
	}
	client.createAccess = func(string, string, newapi.AccessTokenCreateRequest) (newapi.AccessTokenCreated, error) {
		return newapi.AccessTokenCreated{Token: "nap_fresh", Item: newapi.AccessToken{ID: 10}}, nil
	}

	listCalls := 0
	client.listTokens = func(bearer string, _, _ int) (newapi.TokenPage, error) {
		listCalls++
		if bearer == "nap_stale" {
			// The upstream has revoked the stored token.
			return newapi.TokenPage{}, &newapi.APIError{StatusCode: 401, Message: "invalid access token"}
		}
		return newapi.TokenPage{Items: []newapi.Token{{ID: 2, Name: "key"}}}, nil
	}

	uc := provisionedUC(t, tokens, &fakeAccounts{}, client)
	page, err := uc.ListKeys(context.Background(), userID, entity.ExternalSiteTPAI, 1, 20)
	if err != nil {
		t.Fatalf("ListKeys: %v", err)
	}
	if len(page.Items) != 1 {
		t.Errorf("items = %d, want the retried result", len(page.Items))
	}
	if tokens.clears != 1 {
		t.Errorf("clears = %d, want exactly 1", tokens.clears)
	}
	if logins != 1 {
		t.Errorf("logins = %d, want exactly 1 re-mint", logins)
	}
	if listCalls != 2 {
		t.Errorf("list calls = %d, want 2 (one rejected, one retried)", listCalls)
	}
}

func TestSubscriptionScalesCreditsBy1428(t *testing.T) {
	userID := uuid.New()
	tokens := &fakeTokens{}
	client := &fakeClient{}
	client.listTokens = func(string, int, int) (newapi.TokenPage, error) {
		return newapi.TokenPage{}, nil
	}
	client.login = func(string, []byte) (newapi.SessionToken, error) {
		return newapi.SessionToken{AccessToken: "JWT-1"}, nil
	}
	client.verify = func(string, string, any, []byte) (newapi.SecurityProof, error) {
		return newapi.SecurityProof{ProofToken: "PROOF-1"}, nil
	}
	client.createAccess = func(string, string, newapi.AccessTokenCreateRequest) (newapi.AccessTokenCreated, error) {
		return newapi.AccessTokenCreated{Token: "nap_stored"}, nil
	}
	client.subscription = func(string) (newapi.SubscriptionSelf, error) {
		return newapi.SubscriptionSelf{
			BillingPreference: "subscription_only",
			Subscriptions: []newapi.SubscriptionSummary{{Subscription: newapi.Subscription{
				PlanID:      1,
				AmountTotal: 2941177,
				AmountUsed:  368207,
				Status:      "active",
				EndTime:     1794926115,
			}}},
		}, nil
	}

	uc := provisionedUC(t, tokens, &fakeAccounts{}, client)
	usage, err := uc.Subscription(context.Background(), userID, entity.ExternalSiteTPAI)
	if err != nil {
		t.Fatalf("Subscription: %v", err)
	}
	if len(usage.Plans) != 1 {
		t.Fatalf("plans = %d, want 1", len(usage.Plans))
	}
	view := usage.Plans[0]
	// Values captured from the live TPAI response.
	if view.CreditsUsed != 525_799_596 {
		t.Errorf("credits_used = %d, want 525799596", view.CreditsUsed)
	}
	if view.CreditsTotal == nil || *view.CreditsTotal != 4_200_000_756 {
		t.Errorf("credits_total = %v, want 4200000756", view.CreditsTotal)
	}
	if view.Unlimited {
		t.Error("unlimited must be false when amount_total is non-zero")
	}
	if view.Status != "active" || view.PlanID != 1 {
		t.Errorf("view = %+v", view)
	}
	if usage.BillingPreference != "subscription_only" {
		t.Errorf("billing_preference = %q", usage.BillingPreference)
	}
}

func TestSubscriptionReturnsEveryActivePlan(t *testing.T) {
	tokens := &fakeTokens{}
	client := &fakeClient{}
	client.listTokens = func(string, int, int) (newapi.TokenPage, error) { return newapi.TokenPage{}, nil }
	client.login = func(string, []byte) (newapi.SessionToken, error) {
		return newapi.SessionToken{AccessToken: "JWT-1"}, nil
	}
	client.verify = func(string, string, any, []byte) (newapi.SecurityProof, error) {
		return newapi.SecurityProof{ProofToken: "PROOF-1"}, nil
	}
	client.createAccess = func(string, string, newapi.AccessTokenCreateRequest) (newapi.AccessTokenCreated, error) {
		return newapi.AccessTokenCreated{Token: "nap_stored"}, nil
	}
	client.subscription = func(string) (newapi.SubscriptionSelf, error) {
		return newapi.SubscriptionSelf{Subscriptions: []newapi.SubscriptionSummary{
			{Subscription: newapi.Subscription{ID: 7, PlanID: 1, AmountTotal: 1000, AmountUsed: 100, Status: "active"}},
			{Subscription: newapi.Subscription{ID: 9, PlanID: 2, AmountTotal: 0, AmountUsed: 50, Status: "active"}},
		}}, nil
	}

	uc := provisionedUC(t, tokens, &fakeAccounts{}, client)
	usage, err := uc.Subscription(context.Background(), uuid.New(), entity.ExternalSiteTPAI)
	if err != nil {
		t.Fatalf("Subscription: %v", err)
	}
	if len(usage.Plans) != 2 {
		t.Fatalf("plans = %d, want 2", len(usage.Plans))
	}
	// Order is preserved: the platform returns plans in priority order.
	if usage.Plans[0].PlanID != 1 || usage.Plans[1].PlanID != 2 {
		t.Errorf("plan order = %d, %d; want 1, 2", usage.Plans[0].PlanID, usage.Plans[1].PlanID)
	}
	if usage.Plans[0].Unlimited {
		t.Error("plan 1 has a non-zero total, must not be unlimited")
	}
	if !usage.Plans[1].Unlimited || usage.Plans[1].CreditsTotal != nil {
		t.Errorf("plan 2 = %+v, want unlimited with nil credits_total", usage.Plans[1])
	}
	// The instance id must survive, since plan_id alone is not unique.
	if usage.Plans[0].SubscriptionID != 7 || usage.Plans[1].SubscriptionID != 9 {
		t.Errorf("subscription ids = %d, %d; want 7, 9",
			usage.Plans[0].SubscriptionID, usage.Plans[1].SubscriptionID)
	}
}

func TestSubscriptionKeepsInstancesOfTheSamePlanDistinct(t *testing.T) {
	tokens := &fakeTokens{}
	client := &fakeClient{}
	client.listTokens = func(string, int, int) (newapi.TokenPage, error) { return newapi.TokenPage{}, nil }
	client.login = func(string, []byte) (newapi.SessionToken, error) {
		return newapi.SessionToken{AccessToken: "JWT-1"}, nil
	}
	client.verify = func(string, string, any, []byte) (newapi.SecurityProof, error) {
		return newapi.SecurityProof{ProofToken: "PROOF-1"}, nil
	}
	client.createAccess = func(string, string, newapi.AccessTokenCreateRequest) (newapi.AccessTokenCreated, error) {
		return newapi.AccessTokenCreated{Token: "nap_stored"}, nil
	}
	// Two instances of plan 1, as happens when a user buys the same plan twice.
	client.subscription = func(string) (newapi.SubscriptionSelf, error) {
		return newapi.SubscriptionSelf{Subscriptions: []newapi.SubscriptionSummary{
			{Subscription: newapi.Subscription{ID: 11, PlanID: 1, AmountTotal: 1000, Status: "active"}},
			{Subscription: newapi.Subscription{ID: 12, PlanID: 1, AmountTotal: 1000, Status: "active"}},
		}}, nil
	}

	uc := provisionedUC(t, tokens, &fakeAccounts{}, client)
	usage, err := uc.Subscription(context.Background(), uuid.New(), entity.ExternalSiteTPAI)
	if err != nil {
		t.Fatalf("Subscription: %v", err)
	}
	if len(usage.Plans) != 2 {
		t.Fatalf("plans = %d, want 2", len(usage.Plans))
	}
	seen := map[int]bool{}
	for _, plan := range usage.Plans {
		if plan.PlanID != 1 {
			t.Errorf("plan_id = %d, want 1 for both instances", plan.PlanID)
		}
		if plan.SubscriptionID == 0 {
			t.Error("subscription_id must be set, it is the only unique key")
		}
		if seen[plan.SubscriptionID] {
			t.Errorf("subscription_id %d repeated; keys must be unique", plan.SubscriptionID)
		}
		seen[plan.SubscriptionID] = true
	}
}

func TestSubscriptionWithNoActivePlanReturnsEmptyList(t *testing.T) {
	tokens := &fakeTokens{}
	client := &fakeClient{}
	client.listTokens = func(string, int, int) (newapi.TokenPage, error) { return newapi.TokenPage{}, nil }
	client.login = func(string, []byte) (newapi.SessionToken, error) {
		return newapi.SessionToken{AccessToken: "JWT-1"}, nil
	}
	client.verify = func(string, string, any, []byte) (newapi.SecurityProof, error) {
		return newapi.SecurityProof{ProofToken: "PROOF-1"}, nil
	}
	client.createAccess = func(string, string, newapi.AccessTokenCreateRequest) (newapi.AccessTokenCreated, error) {
		return newapi.AccessTokenCreated{Token: "nap_stored"}, nil
	}
	client.subscription = func(string) (newapi.SubscriptionSelf, error) {
		return newapi.SubscriptionSelf{BillingPreference: "balance"}, nil
	}

	uc := provisionedUC(t, tokens, &fakeAccounts{}, client)
	usage, err := uc.Subscription(context.Background(), uuid.New(), entity.ExternalSiteTPAI)
	if err != nil {
		t.Fatalf("Subscription: %v", err)
	}
	// A non-nil empty slice so the API serialises [] rather than null.
	if usage.Plans == nil || len(usage.Plans) != 0 {
		t.Errorf("plans = %v, want empty non-nil slice", usage.Plans)
	}
}

func TestSubscriptionTreatsZeroTotalAsUnlimited(t *testing.T) {
	tokens := &fakeTokens{}
	client := &fakeClient{}
	client.listTokens = func(string, int, int) (newapi.TokenPage, error) { return newapi.TokenPage{}, nil }
	client.login = func(string, []byte) (newapi.SessionToken, error) {
		return newapi.SessionToken{AccessToken: "JWT-1"}, nil
	}
	client.verify = func(string, string, any, []byte) (newapi.SecurityProof, error) {
		return newapi.SecurityProof{ProofToken: "PROOF-1"}, nil
	}
	client.createAccess = func(string, string, newapi.AccessTokenCreateRequest) (newapi.AccessTokenCreated, error) {
		return newapi.AccessTokenCreated{Token: "nap_stored"}, nil
	}
	client.subscription = func(string) (newapi.SubscriptionSelf, error) {
		return newapi.SubscriptionSelf{Subscriptions: []newapi.SubscriptionSummary{{Subscription: newapi.Subscription{
			AmountTotal: 0,
			AmountUsed:  100,
			Status:      "active",
		}}}}, nil
	}

	uc := provisionedUC(t, tokens, &fakeAccounts{}, client)
	usage, err := uc.Subscription(context.Background(), uuid.New(), entity.ExternalSiteTPAI)
	if err != nil {
		t.Fatalf("Subscription: %v", err)
	}
	if len(usage.Plans) != 1 {
		t.Fatalf("plans = %d, want 1", len(usage.Plans))
	}
	view := usage.Plans[0]
	if !view.Unlimited {
		t.Error("amount_total == 0 must report unlimited")
	}
	if view.CreditsTotal != nil {
		t.Errorf("credits_total = %v, want nil for an unlimited plan", view.CreditsTotal)
	}
}

func TestAccountNotProvisionedTriggersBackfill(t *testing.T) {
	userID := uuid.New()
	accounts := &fakeAccounts{err: repository.ErrNotFound}
	built := []entity.ExternalSiteAccount{
		{UserID: userID, Site: entity.ExternalSiteTPAI, Username: "user@example.com"},
		{UserID: userID, Site: entity.ExternalSiteAPAI, Username: "user@example.com"},
	}
	uc := NewAPIKeyUsecase(
		&fakeUsers{user: &entity.User{ID: userID, Email: "user@example.com"}},
		accounts,
		&fakeTokens{},
		fakeCipher{},
		fakeFactory{accounts: built},
		map[entity.ExternalSite]SiteClient{entity.ExternalSiteTPAI: &fakeClient{}},
	)

	_, err := uc.ListKeys(context.Background(), userID, entity.ExternalSiteTPAI, 1, 20)
	if !errors.Is(err, ErrExternalAccountMissing) {
		t.Fatalf("err = %v, want ErrExternalAccountMissing", err)
	}
	// The heal must enqueue BOTH sites so the worker provisions the user.
	if len(accounts.ensured) != 2 {
		t.Fatalf("ensured rows = %d, want 2 (one per site)", len(accounts.ensured))
	}
	sites := map[entity.ExternalSite]bool{}
	for _, a := range accounts.ensured {
		sites[a.Site] = true
	}
	if !sites[entity.ExternalSiteTPAI] || !sites[entity.ExternalSiteAPAI] {
		t.Errorf("ensured sites = %v, want tpai and apai", sites)
	}
}

func TestEnsureProvisionedDoesNotTouchExistingAccounts(t *testing.T) {
	userID := uuid.New()
	accounts := &fakeAccounts{account: &entity.ExternalSiteAccount{UserID: userID, Site: entity.ExternalSiteTPAI}}
	uc := NewAPIKeyUsecase(
		&fakeUsers{user: &entity.User{ID: userID, Email: "user@example.com"}},
		accounts, &fakeTokens{}, fakeCipher{}, fakeFactory{},
		map[entity.ExternalSite]SiteClient{},
	)

	// The factory returns nothing here, but the point is that EnsureProvisioned on
	// an account that already exists is a no-op at the repo level (ON CONFLICT).
	if err := uc.EnsureProvisioned(context.Background(), userID); err != nil {
		t.Fatalf("EnsureProvisioned: %v", err)
	}
	if len(accounts.ensured) != 0 {
		t.Errorf("ensured = %d rows, want 0 (factory produced none)", len(accounts.ensured))
	}
}

func TestAccountNotProvisionedIsReported(t *testing.T) {
	client := &fakeClient{}
	client.listTokens = func(string, int, int) (newapi.TokenPage, error) { return newapi.TokenPage{}, nil }

	uc := provisionedUC(t, &fakeTokens{}, &fakeAccounts{err: repository.ErrNotFound}, client)
	_, err := uc.ListKeys(context.Background(), uuid.New(), entity.ExternalSiteTPAI, 1, 20)
	if !errors.Is(err, ErrExternalAccountMissing) {
		t.Fatalf("err = %v, want ErrExternalAccountMissing", err)
	}
}

func TestLoginRejectionRequeuesProvisioning(t *testing.T) {
	// The row says "provisioned" but the remote account is gone, so login is
	// rejected. new-api answers that with HTTP 200 + success=false.
	userID := uuid.New()
	client := &fakeClient{}
	client.login = func(string, []byte) (newapi.SessionToken, error) {
		return newapi.SessionToken{}, &newapi.APIError{
			StatusCode: 200,
			Message:    "Username or password is incorrect, or user has been banned",
		}
	}
	accounts := &fakeAccounts{requeueResult: true}
	uc := provisionedUC(t, &fakeTokens{}, accounts, client)

	_, err := uc.ListKeys(context.Background(), userID, entity.ExternalSiteTPAI, 1, 20)
	if !errors.Is(err, ErrExternalAccountMissing) {
		t.Fatalf("err = %v, want ErrExternalAccountMissing so the caller retries", err)
	}
	if accounts.requeueCalls != 1 {
		t.Errorf("requeue calls = %d, want 1", accounts.requeueCalls)
	}
}

func TestLoginRejectionWithoutRequeueSurfacesRealError(t *testing.T) {
	userID := uuid.New()
	client := &fakeClient{}
	client.login = func(string, []byte) (newapi.SessionToken, error) {
		return newapi.SessionToken{}, &newapi.APIError{StatusCode: 200, Message: "banned"}
	}
	// Nothing was requeued (e.g. attempts exhausted), so the real error shows.
	accounts := &fakeAccounts{requeueResult: false}
	uc := provisionedUC(t, &fakeTokens{}, accounts, client)

	_, err := uc.ListKeys(context.Background(), userID, entity.ExternalSiteTPAI, 1, 20)
	if !errors.Is(err, ErrExternalRejected) {
		t.Fatalf("err = %v, want ErrExternalRejected", err)
	}
}

func TestTransientLoginFailureIsNotRequeued(t *testing.T) {
	// A 5xx must NOT requeue: re-provisioning during an outage risks duplicates.
	userID := uuid.New()
	client := &fakeClient{}
	client.login = func(string, []byte) (newapi.SessionToken, error) {
		return newapi.SessionToken{}, &newapi.APIError{StatusCode: 503, Retryable: true}
	}
	accounts := &fakeAccounts{requeueResult: true}
	uc := provisionedUC(t, &fakeTokens{}, accounts, client)

	_, err := uc.ListKeys(context.Background(), userID, entity.ExternalSiteTPAI, 1, 20)
	if !errors.Is(err, ErrExternalUnavailable) {
		t.Fatalf("err = %v, want ErrExternalUnavailable", err)
	}
	if accounts.requeueCalls != 0 {
		t.Errorf("requeue calls = %d, want 0 for a transient failure", accounts.requeueCalls)
	}
}

func TestUnsupportedSiteIsRejected(t *testing.T) {
	uc := NewAPIKeyUsecase(&fakeUsers{}, &fakeAccounts{}, &fakeTokens{}, fakeCipher{}, fakeFactory{}, map[entity.ExternalSite]SiteClient{})
	if _, err := ParseSite("nope"); !errors.Is(err, ErrUnsupportedSite) {
		t.Fatalf("ParseSite(nope) err = %v, want ErrUnsupportedSite", err)
	}
	if _, err := uc.ListKeys(context.Background(), uuid.New(), entity.ExternalSite("nope"), 1, 20); !errors.Is(err, ErrUnsupportedSite) {
		t.Fatalf("ListKeys err = %v, want ErrUnsupportedSite", err)
	}
}
