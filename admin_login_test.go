// admin_login_test.go
package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/token"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

// -----------------------------------------------------------------------------
// test helpers
// -----------------------------------------------------------------------------

// spin up an in‑memory Bolt DB + the three admin endpoints wired to httptest
func newLoginTestServer(t *testing.T) (srv *httptest.Server, db *bolt.DB, ts *boltTOTPStore, cleanup func()) {
	t.Helper()

	tmp, err := os.CreateTemp("", "al_bolt_*.db")
	require.NoError(t, err)

	db, err = bolt.Open(tmp.Name(), 0600, &bolt.Options{Timeout: 1 * time.Second})
	require.NoError(t, err)

	// tiny auth.Service good enough for unit tests
	var testSecret token.SecretFunc = func(_ string) (string, error) {
		return "test_secret", nil
	}
	authSvc := auth.NewService(auth.Opts{
		SecretReader:   testSecret,
		TokenDuration:  time.Hour,
		CookieDuration: time.Hour,
		DisableXSRF:    true,
		Issuer:         "AL",
	})
	require.NoError(t, err)

	ts = NewBoltTOTPStore(db)
	cfg := ServerConfig{TOTPIssuer: "AL‑test"}

	mux := http.NewServeMux()
	mux.Handle("/admin/provision", newSysAdminHandler(db, cfg, ts))
	mux.Handle("/admin/login/start", adminLoginStartHandler(authSvc, db, ts))
	mux.Handle("/admin/login/verify", adminLoginVerifyHandler(authSvc, ts))

	srv = httptest.NewServer(mux)

	cleanup = func() {
		srv.Close()
		db.Close()
		_ = os.Remove(tmp.Name())
	}

	return
}

// convenient helpers
func postJSON(t *testing.T, c *http.Client, url string, body any) (*http.Response, []byte) {
	t.Helper()
	b, _ := json.Marshal(body)
	resp, err := c.Post(url, "application/json", bytes.NewReader(b))
	require.NoError(t, err)
	out, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	return resp, out
}

// -----------------------------------------------------------------------------
// tests
// -----------------------------------------------------------------------------

func TestAdminProvisioningAndFullLoginFlow(t *testing.T) {
	srv, db, ts, cleanup := newLoginTestServer(t)
	defer cleanup()

	client := &http.Client{Timeout: 2 * time.Second}

	// 1️⃣  Provision a brand‑new Sys‑Admin
	provReq := map[string]string{
		"Email":       "root@example.com",
		"FirstName":   "Ada",
		"LastName":    "Lovelace",
		"PasswordSha": "s3cretPW!",
	}
	provResp, body := postJSON(t, client, srv.URL+"/admin/provision", provReq)
	assert.Equal(t, http.StatusCreated, provResp.StatusCode)

	var provData struct {
		QRBase64 string `json:"qr_base64"`
		Secret   string `json:"secret"`
	}
	require.NoError(t, json.Unmarshal(body, &provData))
	require.NotEmpty(t, provData.QRBase64)
	require.NotEmpty(t, provData.Secret)

	// make sure secret persisted
	sec, ok := ts.Secret("root@example.com")
	assert.True(t, ok)
	assert.Equal(t, provData.Secret, sec)

	// 2️⃣  Phase‑1 login – correct email/pswd should yield 202 + tmp token
	login1 := map[string]string{"user": "root@example.com", "passwd": "s3cretPW!"}
	l1Resp, l1Body := postJSON(t, client, srv.URL+"/admin/login/start", login1)
	assert.Equal(t, http.StatusAccepted, l1Resp.StatusCode)

	var l1Data struct {
		Status   string `json:"status"`
		TmpToken string `json:"tmp_token"`
	}
	require.NoError(t, json.Unmarshal(l1Body, &l1Data))
	assert.Equal(t, "TOTP_REQUIRED", l1Data.Status)
	require.NotEmpty(t, l1Data.TmpToken)

	// 3️⃣  Phase‑2 login – craft valid TOTP code from the secret we just got
	code, err := totp.GenerateCode(provData.Secret, time.Now())
	require.NoError(t, err)

	login2 := map[string]string{"tmp_token": l1Data.TmpToken, "code": code}
	l2Resp, l2Body := postJSON(t, client, srv.URL+"/admin/login/verify", login2)
	assert.Equal(t, http.StatusOK, l2Resp.StatusCode)

	// should receive a JWT cookie
	cookies := l2Resp.Cookies()
	require.NotEmpty(t, cookies)
	assert.Contains(t, cookies[0].Name, "JWT")

	var pub userPublic
	require.NoError(t, json.Unmarshal(l2Body, &pub))
	assert.Equal(t, "root@example.com", pub.Name)

	// sanity: wrong TOTP → 401
	bad := map[string]string{"tmp_token": l1Data.TmpToken, "code": "000000"}
	respBad, _ := postJSON(t, client, srv.URL+"/admin/login/verify", bad)
	assert.Equal(t, http.StatusUnauthorized, respBad.StatusCode)

	// -----------------------------------------------------------------
	// Extra edge: user w/o TOTP should be blocked at phase‑1
	// -----------------------------------------------------------------
	err = db.Update(func(tx *bolt.Tx) error {
		return createSysAdminTx(tx, "no2fa@example.com", EncodePassword("pw"), "No", "Totp")
	})
	require.NoError(t, err)

	no2fa := map[string]string{"user": "no2fa@example.com", "passwd": "pw"}
	resp, _ := postJSON(t, client, srv.URL+"/admin/login/start", no2fa)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}
