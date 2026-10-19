package newapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// newTestClient points a real Client at an httptest TLS server, preserving the
// production redirect and timeout policy while trusting the test certificate.
func newTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewTLSServer(handler)
	t.Cleanup(srv.Close)

	client, err := NewClient(srv.URL, "admin-token", 5*time.Second)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	client.http.Transport = srv.Client().Transport
	return client, srv
}

func writeEnvelope(t *testing.T, w http.ResponseWriter, payload string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if _, err := io.WriteString(w, payload); err != nil {
		t.Fatalf("write response: %v", err)
	}
}

func decodeBody(t *testing.T, r *http.Request, into any) {
	t.Helper()
	if err := json.NewDecoder(r.Body).Decode(into); err != nil {
		t.Fatalf("decode request body: %v", err)
	}
}

func TestLoginSendsCredentialsAndParsesSession(t *testing.T) {
	var (
		gotPath    string
		gotBody    map[string]string
		gotAuthHdr string
	)
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuthHdr = r.Header.Get("Authorization")
		decodeBody(t, r, &gotBody)
		writeEnvelope(t, w, `{"success":true,"message":"","data":{"access_token":"JWT-1","token_type":"Bearer","access_expires_at":123}}`)
	})

	token, err := client.Login(context.Background(), "user@example.com", []byte("s3cret"))
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if gotPath != "/api/user/login" {
		t.Errorf("path = %q, want /api/user/login", gotPath)
	}
	if gotAuthHdr != "" {
		t.Errorf("login must not send an Authorization header, got %q", gotAuthHdr)
	}
	if gotBody["username"] != "user@example.com" || gotBody["password"] != "s3cret" {
		t.Errorf("body = %v, want the supplied credentials", gotBody)
	}
	if token.AccessToken != "JWT-1" || token.AccessExpiresAt != 123 {
		t.Errorf("token = %+v, want JWT-1/123", token)
	}
}

func TestLoginWithoutAccessTokenIsAnError(t *testing.T) {
	// A pending step-up challenge replaces the session payload.
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeEnvelope(t, w, `{"success":true,"message":"","data":{"method":"2fa","available":true}}`)
	})

	if _, err := client.Login(context.Background(), "user@example.com", []byte("pw")); err == nil {
		t.Fatal("Login must fail when no access_token is returned")
	}
}

func TestVerifySendsProofContextWithBearer(t *testing.T) {
	var (
		gotAuth  string
		gotProof string
		gotPath  string
		gotBody  struct {
			Method   string         `json:"method"`
			Scope    string         `json:"scope"`
			Context  map[string]any `json:"context"`
			Password string         `json:"password"`
			Code     string         `json:"code"`
		}
	)
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		gotProof = r.Header.Get(proofHeader)
		decodeBody(t, r, &gotBody)
		writeEnvelope(t, w, `{"success":true,"data":{"proof_token":"PROOF-1","method":"password","scope":"access_token.generate"}}`)
	})

	proof, err := client.Verify(
		context.Background(),
		"JWT-1",
		ScopeAccessTokenGenerate,
		AccessTokenGenerateContext{Scopes: []string{ScopeAPIKeyRead}, ExpiresAt: 999},
		[]byte("pw"),
	)
	if err != nil {
		t.Fatalf("Verify: %v", err)
	}
	if gotPath != "/api/verify" {
		t.Errorf("path = %q, want /api/verify", gotPath)
	}
	if gotAuth != "Bearer JWT-1" {
		t.Errorf("authorization = %q, want Bearer JWT-1", gotAuth)
	}
	if gotProof != "" {
		t.Errorf("verify must not send %s, got %q", proofHeader, gotProof)
	}
	if gotBody.Method != "password" || gotBody.Scope != "access_token.generate" {
		t.Errorf("body method/scope = %q/%q", gotBody.Method, gotBody.Scope)
	}
	if gotBody.Context["expires_at"] != float64(999) {
		t.Errorf("context.expires_at = %v, want 999", gotBody.Context["expires_at"])
	}
	if proof.ProofToken != "PROOF-1" {
		t.Errorf("proof = %+v, want PROOF-1", proof)
	}
}

func TestCreateAccessTokenSendsProofHeader(t *testing.T) {
	var (
		gotProof string
		gotBody  AccessTokenCreateRequest
		gotPath  string
	)
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotProof = r.Header.Get(proofHeader)
		decodeBody(t, r, &gotBody)
		writeEnvelope(t, w, `{"success":true,"data":{"token":"nap_abc","item":{"id":7,"name":"dxp-gateway"}}}`)
	})

	created, err := client.CreateAccessToken(context.Background(), "JWT-1", "PROOF-1", AccessTokenCreateRequest{
		Name:      "dxp-gateway",
		Scopes:    []string{ScopeAPIKeyRead, ScopeAPIKeyWrite, ScopeWalletRead},
		ExpiresAt: 1000,
	})
	if err != nil {
		t.Fatalf("CreateAccessToken: %v", err)
	}
	if gotPath != "/api/user/access_tokens" {
		t.Errorf("path = %q, want /api/user/access_tokens", gotPath)
	}
	if gotProof != "PROOF-1" {
		t.Errorf("%s = %q, want PROOF-1", proofHeader, gotProof)
	}
	if len(gotBody.Scopes) != 3 || gotBody.Scopes[0] != ScopeAPIKeyRead {
		t.Errorf("scopes = %v", gotBody.Scopes)
	}
	if created.Token != "nap_abc" || created.Item.ID != 7 {
		t.Errorf("created = %+v", created)
	}
}

func TestListTokensUsesPageQueryAndToleratesNullAutoGroups(t *testing.T) {
	var gotQuery string
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.RawQuery
		writeEnvelope(t, w, `{"success":true,"data":{"page":1,"page_size":20,"total":1,"items":[
			{"id":4,"name":"Mimo Code","key":"ECiB**********djAC","status":1,"unlimited_quota":true,
			 "remain_quota":-276313,"expired_time":-1,"allow_ips":"","auto_groups":null,"DeletedAt":null}]}}`)
	})

	page, err := client.ListTokens(context.Background(), "nap_abc", 1, 20)
	if err != nil {
		t.Fatalf("ListTokens: %v", err)
	}
	if gotQuery != "p=1&page_size=20" {
		t.Errorf("query = %q, want p=1&page_size=20", gotQuery)
	}
	if page.Total != 1 || len(page.Items) != 1 {
		t.Fatalf("page = %+v", page)
	}
	if page.Items[0].AutoGroups == nil {
		t.Error("nil auto_groups must normalise to an empty slice")
	}
	if page.Items[0].RemainQuota != -276313 {
		t.Errorf("remain_quota = %d", page.Items[0].RemainQuota)
	}
}

func TestTokenMutationsUseExpectedMethodsAndPaths(t *testing.T) {
	type call struct {
		method, path string
		body         TokenInput
	}
	var calls []call

	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		c := call{method: r.Method, path: r.URL.Path}
		// Only the create/update calls carry a JSON body.
		if r.URL.Path == "/api/token/" && (r.Method == http.MethodPost || r.Method == http.MethodPut) {
			decodeBody(t, r, &c.body)
		}
		calls = append(calls, c)

		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/token/":
			writeEnvelope(t, w, `{"success":true,"message":""}`)
		case r.Method == http.MethodPost && r.URL.Path == "/api/token/5/key":
			writeEnvelope(t, w, `{"success":true,"data":{"key":"sk-revealed"}}`)
		default:
			writeEnvelope(t, w, `{"success":true,"data":{"id":5,"name":"renamed"}}`)
		}
	})

	ctx := context.Background()
	if err := client.CreateToken(ctx, "nap_abc", TokenInput{Name: "new key", UnlimitedQuota: true, ExpiredTime: -1}); err != nil {
		t.Fatalf("CreateToken: %v", err)
	}
	if _, err := client.UpdateToken(ctx, "nap_abc", TokenInput{ID: 5, Name: "renamed"}); err != nil {
		t.Fatalf("UpdateToken: %v", err)
	}
	if err := client.DeleteToken(ctx, "nap_abc", 5); err != nil {
		t.Fatalf("DeleteToken: %v", err)
	}
	key, err := client.RevealTokenKey(ctx, "nap_abc", 5)
	if err != nil {
		t.Fatalf("RevealTokenKey: %v", err)
	}
	if key != "sk-revealed" {
		t.Errorf("key = %q", key)
	}

	want := []call{
		{method: http.MethodPost, path: "/api/token/"},
		{method: http.MethodPut, path: "/api/token/"},
		{method: http.MethodDelete, path: "/api/token/5"},
		{method: http.MethodPost, path: "/api/token/5/key"},
	}
	if len(calls) != len(want) {
		t.Fatalf("calls = %+v, want %+v", calls, want)
	}
	for i := range want {
		if calls[i].method != want[i].method || calls[i].path != want[i].path {
			t.Errorf("call %d = %s %s, want %s %s", i, calls[i].method, calls[i].path, want[i].method, want[i].path)
		}
	}
	if calls[0].body.ID != 0 {
		t.Errorf("create must omit id, got %d", calls[0].body.ID)
	}
	if calls[1].body.ID != 5 {
		t.Errorf("update must carry id 5, got %d", calls[1].body.ID)
	}
}

func TestSubscriptionSelfUnwrapsNestedSubscription(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/subscription/self" {
			t.Errorf("path = %q", r.URL.Path)
		}
		writeEnvelope(t, w, `{"success":true,"data":{"billing_preference":"subscription_only",
			"subscriptions":[{"subscription":{"id":2,"plan_id":1,"amount_total":2941177,"amount_used":368207,"status":"active","end_time":1794926115}}],
			"all_subscriptions":[{"subscription":{"id":2,"amount_total":2941177,"amount_used":368207}}]}}`)
	})

	result, err := client.SubscriptionSelf(context.Background(), "nap_abc")
	if err != nil {
		t.Fatalf("SubscriptionSelf: %v", err)
	}
	if len(result.Subscriptions) != 1 {
		t.Fatalf("subscriptions = %d, want 1", len(result.Subscriptions))
	}
	sub := result.Subscriptions[0].Subscription
	if sub.AmountTotal != 2941177 || sub.AmountUsed != 368207 || sub.Status != "active" {
		t.Errorf("subscription = %+v", sub)
	}
}

func TestBusinessErrorIsNotRetryableAndCarriesMessage(t *testing.T) {
	// Business errors arrive with HTTP 200 and success=false.
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeEnvelope(t, w, `{"success":false,"message":"已达到最大令牌数量限制 (20)"}`)
	})

	_, err := client.ListTokens(context.Background(), "nap_abc", 1, 20)
	if err == nil {
		t.Fatal("expected an error")
	}
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if apiErr.Retryable {
		t.Error("a business error must not be retryable")
	}
	if apiErr.Message == "" {
		t.Error("expected the upstream message to be preserved")
	}
}

func TestUnauthorizedIsFlagged(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		writeEnvelope(t, w, `{"success":false,"message":"invalid access token"}`)
	})

	_, err := client.ListTokens(context.Background(), "nap_stale", 1, 20)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if !apiErr.Unauthorized() {
		t.Errorf("Unauthorized() = false for status %d, want true", apiErr.StatusCode)
	}
	if apiErr.Retryable {
		t.Error("401 must not be retryable")
	}
}

func TestServerErrorIsRetryable(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.ListTokens(context.Background(), "nap_abc", 1, 20)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if !apiErr.Retryable {
		t.Error("5xx must be retryable")
	}
}

func TestStatusReturnsAPIInfoWithoutBearer(t *testing.T) {
	var gotAuth string
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		writeEnvelope(t, w, `{"success":true,"data":{"api_info":[{"id":1,"color":"blue","description":"for sea user","route":"for sea user","url":"https://tp-sg.fromdxp.com"}]}}`)
	})

	status, err := client.Status(context.Background())
	if err != nil {
		t.Fatalf("Status: %v", err)
	}
	if gotAuth != "" {
		t.Errorf("status must not send an Authorization header, got %q", gotAuth)
	}
	if len(status.APIInfo) != 1 || status.APIInfo[0].URL != "https://tp-sg.fromdxp.com" {
		t.Errorf("api_info = %+v", status.APIInfo)
	}
}
