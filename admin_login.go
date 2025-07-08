package main

// admin_login.go consolidates all logic related to provisioning and logging in
// AL Sys‑Admins with TOTP‑based two‑factor authentication.  It keeps the
// existing dev‑friendly in‑memory store but adds a BoltDB‑backed store for
// production, QR‑code provisioning, and a hard requirement that every
// Sys‑Admin has 2‑FA enabled.

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/token"
	"github.com/golang-jwt/jwt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
	bolt "go.etcd.io/bbolt"
)

//---------------------------------------------------------------------
//  Constants & helpers
//---------------------------------------------------------------------

const totpBucket = "totp_secrets" // Bolt bucket keeping persistent secrets

func randHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// generateTOTPSecret returns a fresh base‑32 secret and a PNG (byte slice)
// containing a QR code that encodes the *otpauth://* URL.
func generateTOTPSecret(email, issuer string) (secret string, png []byte, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: email,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", nil, err
	}
	png, err = qrcode.Encode(key.URL(), qrcode.Medium, 256)
	return key.Secret(), png, err
}

//---------------------------------------------------------------------
//  TOTP secret store implementations
//---------------------------------------------------------------------

type TOTPStore interface {
	// Secret returns a base‑32 secret and enabled==true if it exists.
	Secret(email string) (secret string, enabled bool)
	// CacheTemp stores a secret against a tmp‑token for phase‑2 validation.
	CacheTemp(tmpToken, secret string)
	// SecretForTmp retrieves previously cached secret.
	SecretForTmp(tmpToken string) (string, bool)
}

//------------------------------// PROD: BoltDB‑backed

type boltTOTPStore struct {
	db    *bolt.DB
	cache map[string]string // same tmp‑token cache as dev store
}

func NewBoltTOTPStore(db *bolt.DB) *boltTOTPStore {
	// ensure bucket exists
	_ = db.Update(func(tx *bolt.Tx) error {
		_, _ = tx.CreateBucketIfNotExists([]byte(totpBucket))
		return nil
	})
	return &boltTOTPStore{db: db, cache: map[string]string{}}
}

func (b *boltTOTPStore) Secret(email string) (string, bool) {
	var sec []byte
	_ = b.db.View(func(tx *bolt.Tx) error {
		sec = tx.Bucket([]byte(totpBucket)).Get([]byte(email))
		return nil
	})
	if len(sec) == 0 {
		return "", false
	}
	return string(sec), true
}

func (b *boltTOTPStore) CacheTemp(tmpToken, secret string) { b.cache[tmpToken] = secret }
func (b *boltTOTPStore) SecretForTmp(tmpToken string) (string, bool) {
	sec, ok := b.cache[tmpToken]
	return sec, ok
}

//---------------------------------------------------------------------
//  JSON request / response DTOs
//---------------------------------------------------------------------

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

//---------------------------------------------------------------------
//  Provisioning – create a new Sys‑Admin and return QR code
//---------------------------------------------------------------------

func newSysAdminHandler(db *bolt.DB, cfg ServerConfig, ts *boltTOTPStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var req struct {
			Email       string `json:"Email"`
			FirstName   string `json:"FirstName"`
			LastName    string `json:"LastName"`
			PasswordSha string `json:"PasswordSha"` //it comes in clear you need to encode it
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
			return
		}

		if req.Email == "" {
			http.Error(w, "email is mandatory", http.StatusBadRequest)
			return
		}
		if req.FirstName == "" {
			http.Error(w, "first name is mandatory", http.StatusBadRequest)
			return
		}
		if req.LastName == "" {
			http.Error(w, "last name is mandatory", http.StatusBadRequest)
			return
		}
		if req.PasswordSha == "" {
			http.Error(w, "password is mandatory", http.StatusBadRequest)
			return
		}

		req.PasswordSha = EncodePassword(req.PasswordSha)
		var png []byte
		var secret string

		err := db.Update(func(tx *bolt.Tx) error {
			// 1️⃣  create user (implementation found elsewhere in the project)
			err := createSysAdminTx(tx, req.Email, req.PasswordSha, req.FirstName, req.LastName)
			if err != nil {
				return err
			}

			// 2️⃣  generate secret + QR
			secret, png, err = generateTOTPSecret(req.Email, cfg.TOTPIssuer)
			if err != nil {
				return err
			}

			// 3️⃣  persist secret
			err = tx.Bucket([]byte(totpBucket)).Put([]byte(req.Email), []byte(secret))
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// 4️⃣  respond
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(struct {
			QRBase64 string `json:"qr_base64"`
			Secret   string `json:"secret"`
		}{
			QRBase64: "data:image/png;base64," + base64.StdEncoding.EncodeToString(png),
			Secret:   secret,
		})
	})
}

//---------------------------------------------------------------------
//  Login – phase 1 (password) + phase 2 (TOTP)
//---------------------------------------------------------------------

// adminLoginStartHandler verifies email+password and returns a *temporary*
// JWT (no cookie) which the client must exchange with a TOTP code in phase 2.
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

		// 🔐 Ensure they *have* TOTP configured
		secret, enabled := ts.Secret(u.Email)
		if !enabled {
			http.Error(w, "2FA not provisioned for this user", http.StatusForbidden)
			return
		}

		uid := "al_" + token.HashID(sha1.New(), u.Email)
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

// adminLoginVerifyHandler consumes the tmp‑token + TOTP code and returns the
// *real* JWT as a secure cookie.
func adminLoginVerifyHandler(authSvc *auth.Service, ts TOTPStore) http.HandlerFunc {
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
		email := claims.User.Name
		uid := claims.User.ID

		secret, ok := ts.SecretForTmp(req.TmpToken)
		if !ok || !totp.Validate(req.Code, secret) {
			http.Error(w, "invalid code", http.StatusUnauthorized)
			return
		}

		finalClaims := token.Claims{
			User: &token.User{ID: uid, Name: email},
			StandardClaims: jwt.StandardClaims{
				Id:        randHex(16),
				Issuer:    "AL",
				ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
			},
			SessionOnly: false,
		}
		if _, err := authSvc.TokenService().Set(w, finalClaims); err != nil {
			http.Error(w, "token error", http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(userPublic{ID: uid, Name: email})
	}
}
