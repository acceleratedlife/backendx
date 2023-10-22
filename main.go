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
	"log"
	"net/http"
	"os"
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
	OperationDebit          = 1
	OperationCredit         = 2
	keyCharge               = 1.01
	KeyDebt                 = "debt"
	CurrencyUBuck           = "ubuck"
	KeyCollegeJobs          = "collegeJobs"
	KeyJobs                 = "jobs"
	KeyPEvents              = "positiveEvents"
	KeyNEvents              = "negativeEvents"
	KeyAuctions             = "auctions"
	KeyCB                   = "cb"
	KeyUsers                = "users"
	KeyAccounts             = "accounts"
	KeyCryptos              = "cryptos"
	KeyConversion           = "conversion"
	KeyBalance              = "balance"
	KeyBasis                = "basis"
	KeyTransactions         = "transactions"
	KeyTeachers             = "teachers"
	KeySchools              = "schools"
	KeyStudents             = "students"
	KeyClasses              = "classes"
	KeyAddCode              = "addCode"
	KeySettings             = "settings"
	KeyName                 = "name"
	KeyCity                 = "city"
	KeyZip                  = "zip"
	KeyDayPayment           = "dayPayment"
	KeyDayEvent             = "dayEvent"
	KeyPeriod               = "period"
	KeyAdmins               = "admins"
	KeyEmail                = "Email"
	KeyFirstName            = "FirstName"
	KeyLastName             = "LastName"
	KeyCollege              = "College"
	KeyCareerTransition     = "CareerTransition"
	KeyCareerEnd            = "CareerEnd"
	KeyCollegeEnd           = "CollegeEnd"
	KeyHistory              = "History"
	KeyEntireSchool         = "Entire School"
	KeyTeacherClasses       = "teacherClasses"
	KeyFreshman             = "Freshman"
	KeySophomores           = "Sophomores"
	KeyJuniors              = "Juniors"
	KeySeniors              = "Seniors"
	KeyBid                  = "bid"
	KeyMaxBid               = "maxBid"
	KeyDescription          = "description"
	KeyEndDate              = "endDate"
	KeyStartDate            = "startDate"
	KeyOwnerId              = "owner_id"
	KeyVisibility           = "visibility"
	KeyWinnerId             = "winner_id"
	KeyTime                 = "2006-01-02 15:04:05.999999999 -0700 MST"
	KeyValue                = "value"
	KeyMMA                  = "MMA"
	KeyModMMA               = "modMMA"
	KeyPayFrequency         = "payFrequency"
	KeyRegEnd               = "regEnd"
	KeyCoins                = "ethereum,cardano,bitcoin,chainlink,bnb,xrp,solana,dogecoin,polkadot,shiba-inu,dai,polygon,tron,avalanche,okb,litecoin,ftx,cronos,monery,uniswap,stellar,algorand,chain,flow,vechain,filecoin,frax,apecoin,hedera,eos,decentraland,tezos,quant,elrond,chillz,aave,kucoin,zcash,helium,fantom"
	LoanRate                = 1.015
	KeyMarket               = "market"
	KeyMarketData           = "marketData"
	KeyBuyers               = "buyers"
	KeyLotteries            = "lotteries"
	KeyPricePerTicket       = 5
	KeyCertificateOfDeposit = "certificateOfDeposit"
)

var build_date string

type ServerConfig struct {
	AdminPassword string
	SecureCookies bool
	EnableXSRF    bool
	SecretKey     string
	ServerPort    int
	SeedPassword  string
	EmailSMTP     string
	PasswordSMTP  string
	Production    bool
}

type Clock interface {
	Now() time.Time
}

type AppClock struct {
}

type DemoClock struct {
	Future time.Duration
}

func (t *DemoClock) Now() time.Time {
	lgr.Printf("DEBUG current time - %v", time.Now().Add(t.Future))
	return time.Now().Add(t.Future)
}
func (t *DemoClock) TickOne(d time.Duration) {
	t.Future = t.Future + d
}

func (t *DemoClock) ResetNow() {
	t.Future = time.Duration(0)
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

	runEveryMinute(db)

	// ***

	authService := initAuth(db, config)
	authRoute, _ := authService.Handlers()
	m := authService.Middleware()

	// *** auth
	router, clock := createRouter(db)
	router.Handle("/auth/al/login", authRoute)
	router.Handle("/auth/al/logout", authRoute)

	// backup
	router.Handle("/admin/backup", backUpHandler(db))
	router.HandleFunc("/admin/version", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/octet-stream")
		writer.Write([]byte(build_date))
	})
	//new school
	router.Handle("/admin/new-school", newSchoolHandler(db, &AppClock{}))
	//reset staff password
	router.Handle("/admin/resetPassword", resetPasswordHandler(db))
	//add job details
	router.Handle("/admin/addJobs", addJobsHandler(db))
	//add life event
	router.Handle("/admin/addEvents", addEventsHandler(db))
	//add admin
	router.Handle("/admin/addAdmin", addAdminHandler(db))
	//seed db, dev only
	router.Handle("/admin/seedDb", seedDbHandler(db, clock))
	//advance clock 15 days, dev only
	router.Handle("/admin/nextCollege", nextCollegeHandler(clock))
	//advance clock 5 days, dev only
	router.Handle("/admin/nextCareer", nextCareerHandler(clock))
	//advance clock 24 hours, dev only
	router.Handle("/admin/nextDay", nextDayHandler(clock))
	//advance clock 1 hour, dev only
	router.Handle("/admin/nextHour", nextHourHandler(clock))
	//advance clock 10 minutes, dev only
	router.Handle("/admin/nextMinutes", nextMinutesHandler(clock))
	//reset clock to current time
	router.Handle("/admin/resetClock", resetClockHandler(clock))

	router.Use(buildAuthMiddleware(m))

	addr := fmt.Sprintf(":%d", config.ServerPort)
	log.Fatal(http.ListenAndServe(addr, router))

}

// creates routes for prod
func createRouter(db *bolt.DB) (*mux.Router, *DemoClock) {
	serverConfig := loadConfig()
	if serverConfig.Production {
		clock := &AppClock{}
		return createRouterClock(db, clock), nil
	}

	clock := &DemoClock{}
	return createRouterClock(db, clock), clock
}

func createRouterClock(db *bolt.DB, clock Clock) *mux.Router {

	StudentApiServiceImpl := NewStudentApiServiceImpl(db, clock)
	StudentApiController := openapi.NewStudentApiController(StudentApiServiceImpl)

	SchoolAdminApiService := NewSchoolAdminServiceImpl(db)
	SchoolAdminApiController := openapi.NewSchoolAdminApiController(SchoolAdminApiService)

	allService := NewAllApiServiceImpl(db, clock)
	allController := openapi.NewAllApiController(allService)

	sysAdminApiServiceImpl := NewSysAdminApiServiceImpl(db)
	sysAdminCtrl := openapi.NewSysAdminApiController(sysAdminApiServiceImpl)

	unregisteredServiceImpl := NewUnregisteredApiServiceImpl(db, clock)
	unregisteredApiController := openapi.NewUnregisteredApiController(unregisteredServiceImpl)
	openapi.WithUnregisteredApiErrorHandler(ErrorHandler)

	allSchoolApiServiceImpl := NewAllSchoolApiServiceImpl(db, clock)
	schoolApiController := openapi.NewAllSchoolApiController(allSchoolApiServiceImpl)

	staffApiServiceImpl := NewStaffApiServiceImpl(db, clock)
	staffApiController := openapi.NewStaffApiController(staffApiServiceImpl)

	return openapi.NewRouter(SchoolAdminApiController,
		allController,
		sysAdminCtrl,
		schoolApiController,
		staffApiController,
		unregisteredApiController,
		StudentApiController)

}

func InitDefaultAccounts(db *bolt.DB, clock Clock) {
	newSchoolRequest := NewSchoolRequest{
		School:    "test school",
		FirstName: "test",
		LastName:  "admin",
		Email:     "test@admin.com",
		City:      "Stockton",
		Zip:       95336,
	}
	_ = createNewSchool(db, clock, newSchoolRequest, "123qwe")
}

func createNewSchool(db *bolt.DB, clock Clock, newSchoolRequest NewSchoolRequest, adminPassword string) error {

	schoolId, err := FindOrCreateSchool(db, clock, newSchoolRequest.School, newSchoolRequest.City, newSchoolRequest.Zip)
	if err != nil {
		lgr.Printf("ERROR school does not exist: %v", err)
	}

	lgr.Printf("INFO test school id - %s", schoolId)

	admin := UserInfo{
		Name:        newSchoolRequest.Email,
		Email:       newSchoolRequest.Email,
		PasswordSha: EncodePassword(adminPassword),
		FirstName:   newSchoolRequest.FirstName,
		LastName:    newSchoolRequest.LastName,
		Role:        UserRoleAdmin,
		SchoolId:    schoolId,
	}

	_, err = CreateSchoolAdmin(db, admin)
	if err != nil {
		lgr.Printf("ERROR school admin is not created: %v", err)
		return err
	}

	tEmail := newSchoolRequest.Email[:1] + "." + newSchoolRequest.Email[1:]
	lgr.Printf("INFO teacher's email: %s", tEmail)

	teacher := UserInfo{
		Name:        tEmail,
		Email:       tEmail,
		PasswordSha: EncodePassword(adminPassword),
		FirstName:   newSchoolRequest.FirstName,
		LastName:    newSchoolRequest.LastName,
		Role:        UserRoleTeacher,
		SchoolId:    schoolId,
		Settings: TeacherSettings{
			CurrencyLock: false,
		},
	}

	err = createTeacher(db, teacher)
	if err != nil {
		lgr.Printf("ERROR teacher user is not created")
	}

	return nil
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
			} else if r.URL.Path == "/api/users/register" || r.URL.Path == "/api/users/resetStaffPassword" {
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
		ServerPort:    5000,
		SeedPassword:  "123qwe",
		EmailSMTP:     "qq@qq.com",
		PasswordSMTP:  "123qwe",
		Production:    false,
	}

	yamlFile, err := os.ReadFile("./alcfg.yml")
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

func seedDb(db *bolt.DB, clock Clock, eventRequests []EventRequest, jobRequests []Job) (err error) {

	config := loadConfig()

	school := NewSchoolRequest{
		School:    "JHS",
		FirstName: "Tom",
		LastName:  "Jones",
		Email:     "aa@aa.com",
		City:      "Tacoma",
		Zip:       94558,
	}

	err = createNewSchool(db, clock, school, config.SeedPassword)
	if err != nil {
		lgr.Printf("ERROR school is not created: %v", err)
		return err
	}

	admin, err := getUserInLocalStore(db, "aa@aa.com")
	if err != nil {
		return err
	}

	schoolId := admin.SchoolId

	err = createJobs(db, jobRequests)
	if err != nil {
		return err
	}

	err = createEvents(db, eventRequests)
	if err != nil {
		return err
	}

	err = createTeachers(db, clock, schoolId, config.SeedPassword)
	if err != nil {
		return err
	}

	return

}
