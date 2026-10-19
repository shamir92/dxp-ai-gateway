package newapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// SessionToken is the short-lived bearer returned by POST /api/user/login. It is
// only used to drive the access-token mint sequence; it is never persisted.
type SessionToken struct {
	AccessToken     string `json:"access_token"`
	TokenType       string `json:"token_type"`
	AccessExpiresAt int64  `json:"access_expires_at"`
}

// Login exchanges the provisioned username and password for a browser-session
// bearer. The password is decrypted by the caller and cleared afterwards.
func (c *Client) Login(ctx context.Context, username string, password []byte) (SessionToken, error) {
	payload := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{Username: username, Password: string(password)}

	data, err := c.do(ctx, requestOptions{
		method:  http.MethodPost,
		path:    loginPath,
		referer: "/login",
		body:    payload,
	})
	if err != nil {
		return SessionToken{}, err
	}

	var token SessionToken
	if err := json.Unmarshal(data, &token); err != nil {
		return SessionToken{}, fmt.Errorf("decode login response")
	}
	if strings.TrimSpace(token.AccessToken) == "" {
		// A pending step-up challenge replaces the session payload; the password
		// login itself succeeded but needs additional verification.
		return SessionToken{}, fmt.Errorf("login requires additional verification")
	}
	return token, nil
}

// SecurityProof is a single-use proof authorising exactly one operation. It is
// short-lived and consumed by the upstream on first use.
type SecurityProof struct {
	ProofToken string `json:"proof_token"`
	ExpiresAt  int64  `json:"expires_at"`
	Method     string `json:"method"`
	Scope      string `json:"scope"`
}

// AccessTokenGenerateContext is the verification context required to mint an
// access token. Both fields are mandatory and the keys must be exactly these.
type AccessTokenGenerateContext struct {
	Scopes    []string `json:"scopes"`
	ExpiresAt int64    `json:"expires_at"`
}

// Verify requests a proof for one scope. The bearer must be a browser-session
// token from Login: access-token generation can never be proved by a PAT.
func (c *Client) Verify(ctx context.Context, bearer, scope string, proofContext any, password []byte) (SecurityProof, error) {
	payload := struct {
		Method   string `json:"method"`
		Scope    string `json:"scope"`
		Context  any    `json:"context,omitempty"`
		Password string `json:"password,omitempty"`
	}{
		Method:   verificationMethodPassword,
		Scope:    scope,
		Context:  proofContext,
		Password: string(password),
	}

	data, err := c.do(ctx, requestOptions{
		method:  http.MethodPost,
		path:    verifyPath,
		bearer:  bearer,
		referer: "/keys",
		body:    payload,
	})
	if err != nil {
		return SecurityProof{}, err
	}

	var proof SecurityProof
	if err := json.Unmarshal(data, &proof); err != nil {
		return SecurityProof{}, fmt.Errorf("decode verification response")
	}
	if strings.TrimSpace(proof.ProofToken) == "" {
		return SecurityProof{}, fmt.Errorf("verification did not return a proof")
	}
	return proof, nil
}

// AccessToken is a stored personal access token ("nap_…").
type AccessToken struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	TokenRef  string   `json:"token_ref"`
	TokenHint string   `json:"token_hint"`
	ExpiresAt int64    `json:"expires_at"`
	Scopes    []string `json:"scopes"`
	CreatedAt int64    `json:"created_at"`
}

// AccessTokenCreated carries the plaintext token. It is returned by the upstream
// exactly once, at creation time.
type AccessTokenCreated struct {
	Token string      `json:"token"`
	Item  AccessToken `json:"item"`
}

// AccessTokenCreateRequest is the mint body. ExpiresAt uses 0 for "never".
type AccessTokenCreateRequest struct {
	Name      string   `json:"name"`
	Scopes    []string `json:"scopes"`
	ExpiresAt int64    `json:"expires_at"`
}

// CreateAccessToken mints an access token using the proof from Verify. The
// proof header must carry the raw proof token.
func (c *Client) CreateAccessToken(ctx context.Context, bearer, proofToken string, request AccessTokenCreateRequest) (AccessTokenCreated, error) {
	data, err := c.do(ctx, requestOptions{
		method:  http.MethodPost,
		path:    accessTokensPath,
		bearer:  bearer,
		referer: "/keys",
		body:    request,
		headers: map[string]string{proofHeader: proofToken},
	})
	if err != nil {
		return AccessTokenCreated{}, err
	}

	var created AccessTokenCreated
	if err := json.Unmarshal(data, &created); err != nil {
		return AccessTokenCreated{}, fmt.Errorf("decode access token response")
	}
	if strings.TrimSpace(created.Token) == "" {
		return AccessTokenCreated{}, fmt.Errorf("access token response did not include a token")
	}
	return created, nil
}
