package main

import (
	"encoding/json"
	"net/http"
	"time"

	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"

	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/token"
	"github.com/golang-jwt/jwt"
	"github.com/pquerna/otp/totp"
	bolt "go.etcd.io/bbolt"
)

// TOTPStore provides secrets and a temporary-token cache for 2FA
// You need a concrete implementation, e.g. the in-memory store below.

type memTOTPStore struct {
	data  map[string]string // email   -> base32 secret
	cache map[string]string // tmpTok  -> secret
}

// NewMemTOTPStore returns a dev-only in-memory store for TOTP secrets.
func NewMemTOTPStore() *memTOTPStore {
	return &memTOTPStore{
		data:  map[string]string{},
		cache: map[string]string{},
	}
}

// Secret returns the saved secret for an email, enabled==true if present.
func (s *memTOTPStore) Secret(email string) (string, bool) {
	sec, ok := s.data[email]
	return sec, ok
}

// CacheTemp caches a secret under the tmpToken until verification.
func (s *memTOTPStore) CacheTemp(tmpToken, secret string) {
	s.cache[tmpToken] = secret
}

// SecretForTmp retrieves the secret by tmpToken for phase-2 validation.
func (s *memTOTPStore) SecretForTmp(tmpToken string) (string, bool) {
	sec, ok := s.cache[tmpToken]
	return sec, ok
}

// TOTPStore interface
// Secret(email) returns base32 secret and enabled flag
// CacheTemp stores secret against tmpToken
// SecretForTmp retrieves secret by tmpToken

type TOTPStore interface {
	Secret(email string) (secret string, enabled bool)
	CacheTemp(tmpToken, secret string)
	SecretForTmp(tmpToken string) (string, bool)
}

type adminLoginReq struct {
	User   string `json:"user"`
	Passwd string `json:"passwd"`
}

type adminTOTPReq struct {
	TmpToken string `json:"tmp_token"`
	Code     string `json:"code"`
}

type userPublic struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture,omitempty"`
}

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Phase 1: password check, issue tmp-token only via Token(), no cookie
func adminLoginStartHandler(authSvc *auth.Service, db *bolt.DB, ts TOTPStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req adminLoginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		ok, _ := checkUserInLocalStore(db, req.User, req.Passwd)
		if !ok {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		u, _ := getUserInLocalStore(db, req.User)
		if u.Role != UserRoleSysAdmin {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		uid := "al_" + token.HashID(sha1.New(), u.Email)
		secret, enabled := ts.Secret(u.Email)
		if !enabled {
			// no TOTP: issue final JWT here
			claims := token.Claims{
				User: &token.User{ID: uid, Name: u.Email},
				StandardClaims: jwt.StandardClaims{
					Id:        randHex(16),
					Issuer:    "AL",
					ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
				},
				SessionOnly: false,
			}
			_, err := authSvc.TokenService().Set(w, claims)
			if err != nil {
				http.Error(w, "token error", http.StatusInternalServerError)
				return
			}
			_ = json.NewEncoder(w).Encode(userPublic{ID: uid, Name: u.Email})
			return
		}
		// TOTP enabled: issue tmp token without Set()
		claims := token.Claims{
			User: &token.User{ID: uid, Name: u.Email},
			StandardClaims: jwt.StandardClaims{
				Id:        randHex(16),
				Issuer:    "AL",
				ExpiresAt: time.Now().Add(2 * time.Minute).Unix(),
			},
			SessionOnly: false,
		}
		tmpToken, err := authSvc.TokenService().Token(claims)
		if err != nil {
			http.Error(w, "token error", http.StatusInternalServerError)
			return
		}
		ts.CacheTemp(tmpToken, secret)
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(struct {
			Status   string `json:"status"`
			TmpToken string `json:"tmp_token"`
		}{"TOTP_REQUIRED", tmpToken})
	}
}

// Phase 2: verify TOTP then issue real JWT via Set()
func adminLoginVerifyHandler(authSvc *auth.Service, db *bolt.DB, ts TOTPStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req adminTOTPReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		claims, err := authSvc.TokenService().Parse(req.TmpToken)
		if err != nil || claims.ExpiresAt < time.Now().Unix() {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		uid := claims.User.ID
		email := claims.User.Name
		secret, ok := ts.SecretForTmp(req.TmpToken)
		if !ok || !totp.Validate(req.Code, secret) {
			http.Error(w, "invalid code", http.StatusUnauthorized)
			return
		}
		// issue final JWT
		finalClaims := token.Claims{
			User: &token.User{ID: uid, Name: email},
			StandardClaims: jwt.StandardClaims{
				Id:        randHex(16),
				Issuer:    "AL",
				ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
			},
			SessionOnly: false,
		}
		_, err = authSvc.TokenService().Set(w, finalClaims)
		if err != nil {
			http.Error(w, "token error", http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(userPublic{ID: uid, Name: email})
	}
}
