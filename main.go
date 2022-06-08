/*
 * OpenAPI Petstore
 *
 * This is a sample server Petstore server. For this sample, you can use the api key `special-key` to test the authorization filters.
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/middleware"
	"github.com/go-pkgz/auth/provider"
	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/lgr"
	"github.com/gorilla/mux"
	bolt "go.etcd.io/bbolt"
)

const (
	OperationDebit      = 1
	OperationCredit     = 2
	CurrencyUBuck       = "ubuck"
	KeyAuctions         = "auctions"
	KeyCB               = "cb"
	KeyUsers            = "users"
	KeyAccounts         = "accounts"
	KeybAccounts        = "bAccounts"
	KeycAccounts        = "cAccounts"
	KeyBalance          = "balance"
	KeyTransactions     = "transactions"
	KeyTeachers         = "teachers"
	KeySchools          = "schools"
	KeyStudents         = "students"
	KeyClasses          = "classes"
	KeyAddCode          = "addCode"
	KeyName             = "name"
	KeyCity             = "city"
	KeyZip              = "zip"
	KeyDayPayment       = "dayPayment"
	KeyPeriod           = "period"
	KeyAdmins           = "admins"
	KeyEmail            = "Email"
	KeyFirstName        = "FirstName"
	KeyLastName         = "LastName"
	KeyCollege          = "College"
	KeyCareerTransition = "CareerTransition"
	KeyCareerEnd        = "CareerEnd"
	KeyCollegeEnd       = "CollegeEnd"
	KeyHistory          = "History"
	KeyEntireSchool     = "Entire School"
	KeyTeacherClasses   = "teacherClasses"
	KeyFreshman         = "Freshman"
	KeySophomores       = "Sophomores"
	KeyJuniors          = "Juniors"
	KeySeniors          = "Seniors"
	KeyBid              = "bid"
	KeyMaxBid           = "maxBid"
	KeyDescription      = "description"
	KeyEndDate          = "endDate"
	KeyStartDate        = "startDate"
	KeyOwnerId          = "owner_id"
	KeyVisibility       = "visibility"
	KeyWinnerId         = "winner_id"
	KeyTime             = "2006-01-02 15:04:05.999999999 -0700 MST"
	KeyValue            = "value"
	KeyMMA              = "MMA"
)

type ServerConfig struct {
	AdminPassword string
	SecureCookies bool
	EnableXSRF    bool
	SecretKey     string
}

type Clock interface {
	Now() time.Time
}

type AppClock struct {
}

func (*AppClock) Now() time.Time {
	return time.Now()
}

func main() {
	lgr.Printf("server started")

	config := loadConfig()

	db, err := bolt.Open("al.db", 0666, nil)
	if err != nil {
		panic(fmt.Errorf("cannot open db %v", err))
	}
	defer db.Close()

	InitDefaultAccounts(db)

	// ***

	authService := initAuth(db, config)
	authRoute, _ := authService.Handlers()
	m := authService.Middleware()

	// *** auth
	router := createRouter(db)
	router.Handle("/auth/al/login", authRoute)
	router.Handle("/auth/al/logout", authRoute)

	// backup
	router.Handle("/admin/backup", backUpHandler(db))

	router.Use(buildAuthMiddleware(m))

	log.Fatal(http.ListenAndServe(":5000", router))
}

func createRouter(db *bolt.DB) *mux.Router {
	clock := &AppClock{}

	teacherApiServiceImpl := NewTeacherApiServiceImpl(db)
	teacherApiController := openapi.NewTeacherApiController(teacherApiServiceImpl)

	StudentApiServiceImpl := NewStudentApiServiceImpl(db)
	StudentApiController := openapi.NewStudentApiController(StudentApiServiceImpl)

	SchoolAdminApiService := NewSchoolAdminServiceImpl(db)
	SchoolAdminApiController := openapi.NewSchoolAdminApiController(SchoolAdminApiService)

	allService := NewAllApiServiceImpl(db, clock)
	allController := openapi.NewAllApiController(allService)

	sysAdminApiServiceImpl := NewSysAdminApiServiceImpl(db)
	sysAdminCtrl := openapi.NewSysAdminApiController(sysAdminApiServiceImpl)

	unregisteredServiceImpl := NewUnregisteredApiServiceImpl(db)
	unregisteredApiController := openapi.NewUnregisteredApiController(unregisteredServiceImpl)
	openapi.WithUnregisteredApiErrorHandler(ErrorHandler)

	allSchoolApiServiceImpl := NewAllSchoolApiServiceImpl(db)
	schoolApiController := openapi.NewAllSchoolApiController(allSchoolApiServiceImpl)

	staffApiServiceImpl := NewStaffApiServiceImpl(db, clock)
	staffApiController := openapi.NewStaffApiController(staffApiServiceImpl)

	return openapi.NewRouter(SchoolAdminApiController,
		allController,
		sysAdminCtrl,
		schoolApiController,
		staffApiController,
		unregisteredApiController,
		StudentApiController,
		teacherApiController)

}

func InitDefaultAccounts(db *bolt.DB) {

	schoolId, err := FindOrCreateSchool(db, "test school", "no city", 0)
	if err != nil {
		lgr.Printf("ERROR school does not exist: %v", err)
	}

	lgr.Printf("INFO test school id - %s", schoolId)

	admin := UserInfo{
		Name:        "test@admin.com",
		Email:       "test@admin.com",
		PasswordSha: EncodePassword("123qwe"),
		FirstName:   "test",
		LastName:    "admin",
		Role:        UserRoleAdmin,
		SchoolId:    schoolId,
	}

	_, err = CreateSchoolAdmin(db, admin)
	if err != nil {
		lgr.Printf("ERROR school admin is not created: %v", err)
		return
	}

}

func initAuth(db *bolt.DB, config ServerConfig) *auth.Service {
	options := auth.Opts{
		SecretReader: token.SecretFunc(func(id string) (string, error) { // secret key for JWT
			return config.SecretKey, nil
		}),
		DisableXSRF:    !config.EnableXSRF,
		TokenDuration:  time.Minute * 5, // token expires in 5 minutes
		CookieDuration: time.Hour * 24,  // cookie expires in 1 day and will enforce re-login
		SecureCookies:  config.SecureCookies,
		Issuer:         "AL",
		//URL:            "http://127.0.0.1:8080",
		AvatarStore: avatar.NewNoOp(),
		Validator: token.ValidatorFunc(func(_ string, claims token.Claims) bool {
			// allow only dev_* names
			//return claims.User != nil && strings.HasPrefix(claims.User.Name, "dev_")
			return true
		}),
		AdminPasswd: config.AdminPassword,
	}

	// create auth service with providers
	service := auth.NewService(options)
	service.AddDirectProvider("al", provider.CredCheckerFunc(func(user, password string) (ok bool, err error) {
		ok, err = checkUserInLocalStore(db, user, password)
		return
	}))
	return service
}

func buildAuthMiddleware(m middleware.Authenticator) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// if not authentication related pass through auth
			if strings.HasPrefix(r.URL.Path, "/auth") {
				handler.ServeHTTP(w, r)
			} else if r.URL.Path == "/api/users/register" {
				// or auth-free
				handler.ServeHTTP(w, r)

			} else {
				h := m.Auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					userInfo, err := token.GetUserInfo(r)
					if err != nil {
						lgr.Printf("failed to get user info, %s", err)
						w.WriteHeader(http.StatusForbidden)
						return
					}
					ctx := context.WithValue(r.Context(), "user", userInfo)
					handler.ServeHTTP(w, r.WithContext(ctx))
				}))
				h.ServeHTTP(w, r)
			}
			return
		})
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error, result *openapi.ImplResponse) {
	lgr.Printf("ERROR %s", err)
}

func loadConfig() ServerConfig {
	config := ServerConfig{
		SecretKey:     "secret",
		AdminPassword: "admin",
	}

	yamlFile, err := ioutil.ReadFile("./alcfg.yml")
	if err != nil {
		lgr.Printf("ERROR cannot load config %v", err)
		return config
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		lgr.Printf("ERROR cannot decode config %v", err)
	}

	return config
}
