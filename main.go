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
	"encoding/json"
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
	keyCharge           = 1.01
	KeyDebt             = "debt"
	CurrencyUBuck       = "ubuck"
	KeyCollegeJobs      = "collegeJobs"
	KeyJobs             = "jobs"
	KeyPEvents          = "positiveEvents"
	KeyNEvents          = "negativeEvents"
	KeyAuctions         = "auctions"
	KeyCB               = "cb"
	KeyUsers            = "users"
	KeyAccounts         = "accounts"
	KeyCryptos          = "cryptos"
	KeyConversion       = "conversion"
	KeyBalance          = "balance"
	KeyBasis            = "basis"
	KeyTransactions     = "transactions"
	KeyTeachers         = "teachers"
	KeySchools          = "schools"
	KeyStudents         = "students"
	KeyClasses          = "classes"
	KeyAddCode          = "addCode"
	KeySettings         = "settings"
	KeyName             = "name"
	KeyCity             = "city"
	KeyZip              = "zip"
	KeyDayPayment       = "dayPayment"
	KeyDayEvent         = "dayEvent"
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
	//new school
	router.Handle("/admin/new-school", newSchoolHandler(db))
	//reset staff password
	router.Handle("/admin/resetPassword", resetPasswordHandler(db))
	//add job details
	router.Handle("/admin/addJobs", addJobsHandler(db))
	//add life event
	router.Handle("/admin/addEvents", addEventsHandler(db))
	//add admin
	router.Handle("/admin/addAdmin", addAdminHandler(db))
	//seed db
	router.Handle("/admin/seedDb", seedDbHandler(db))

	router.Use(buildAuthMiddleware(m))

	log.Fatal(http.ListenAndServe(":5000", router))

}

// creates routes for prod
func createRouter(db *bolt.DB) *mux.Router {
	clock := &AppClock{}
	return createRouterClock(db, clock)
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
		StudentApiController)

}

func InitDefaultAccounts(db *bolt.DB) {
	newSchoolRequest := NewSchoolRequest{
		School:    "test school",
		FirstName: "test",
		LastName:  "admin",
		Email:     "test@admin.com",
		City:      "Stockton",
		Zip:       95336,
	}
	_ = createNewSchool(db, newSchoolRequest, "123qwe")
}

func createNewSchool(db *bolt.DB, newSchoolRequest NewSchoolRequest, adminPassword string) error {

	schoolId, err := FindOrCreateSchool(db, newSchoolRequest.School, newSchoolRequest.City, newSchoolRequest.Zip)
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

func seedDb(db *bolt.DB) (err error) {

	school := NewSchoolRequest{
		School:    "JHS",
		FirstName: "Tom",
		LastName:  "Jones",
		Email:     "aa@aa.com",
		City:      "Tacoma",
		Zip:       94558,
	}

	err = createNewSchool(db, school, "123qwe")
	if err != nil {
		lgr.Printf("ERROR school is not created: %v", err)
		return err
	}

	admin, err := getUserInLocalStore(db, "aa@aa.com")
	if err != nil {
		return err
	}

	schoolId := admin.SchoolId

	job := Job{
		Pay:         52000,
		Description: "Teach Stuff",
	}

	marshal, err := json.Marshal(job)

	err = createJobOrEvent(db, marshal, KeyCollegeJobs, "Teacher")
	if err != nil {
		return err
	}

	job = Job{
		Pay:         102000,
		Description: "Rehab Injuries",
	}

	marshal, err = json.Marshal(job)

	err = createJobOrEvent(db, marshal, KeyCollegeJobs, "Physical Therapist")
	if err != nil {
		return err
	}

	job = Job{
		Pay:         160128,
		Description: "Oversee progress of large project",
	}

	marshal, err = json.Marshal(job)

	err = createJobOrEvent(db, marshal, KeyCollegeJobs, "Division Director")
	if err != nil {
		return err
	}

	job = Job{
		Pay:         26000,
		Description: "Wait on tables",
		College:     false,
	}

	marshal, err = json.Marshal(job)

	err = createJobOrEvent(db, marshal, KeyJobs, "Waiter")
	if err != nil {
		return err
	}

	job = Job{
		Pay:         45000,
		Description: "run office",
		College:     false,
	}

	marshal, err = json.Marshal(job)

	err = createJobOrEvent(db, marshal, KeyJobs, "Secretary")
	if err != nil {
		return err
	}

	job = Job{
		Pay:         26000,
		Description: "organize warehouse",
		College:     false,
	}

	marshal, err = json.Marshal(job)

	err = createJobOrEvent(db, marshal, KeyJobs, "Logistics")
	if err != nil {
		return err
	}

	event := eventRequest{
		Description: "Pay Taxes",
	}

	marshal, _ = json.Marshal(event)

	err = createJobOrEvent(db, marshal, KeyNEvents, "Taxes")
	if err != nil {
		return err
	}

	event = eventRequest{
		Description: "Broken Arm",
	}

	marshal, _ = json.Marshal(event)

	err = createJobOrEvent(db, marshal, KeyNEvents, "Injury")
	if err != nil {
		return err
	}

	event = eventRequest{
		Description: "crashed car",
	}

	marshal, _ = json.Marshal(event)

	err = createJobOrEvent(db, marshal, KeyNEvents, "car accident")
	if err != nil {
		return err
	}

	event = eventRequest{
		Description: "won beauty contest",
	}

	marshal, _ = json.Marshal(event)

	err = createJobOrEvent(db, marshal, KeyPEvents, "beauty")
	if err != nil {
		return err
	}

	event = eventRequest{
		Description: "created a viral video",
	}

	marshal, _ = json.Marshal(event)

	err = createJobOrEvent(db, marshal, KeyPEvents, "video")
	if err != nil {
		return err
	}

	event = eventRequest{
		Description: "won fortnite tournament",
	}

	marshal, _ = json.Marshal(event)

	err = createJobOrEvent(db, marshal, KeyPEvents, "fortnite")
	if err != nil {
		return err
	}

	newUser := UserInfo{
		Name:        "tt@tt.com",
		FirstName:   "tt",
		LastName:    "tt",
		Email:       "tt@tt.com",
		Confirmed:   true,
		PasswordSha: EncodePassword("123qwe"),
		SchoolId:    schoolId,
		Role:        UserRoleTeacher,
	}

	err = createTeacher(db, newUser)
	if err != nil {
		return err
	}

	newUser = UserInfo{
		Name:        "tt1@tt.com",
		FirstName:   "tt1",
		LastName:    "tt1",
		Email:       "tt1@tt.com",
		Confirmed:   true,
		PasswordSha: EncodePassword("123qwe"),
		SchoolId:    schoolId,
		Role:        UserRoleTeacher,
	}

	err = createTeacher(db, newUser)
	if err != nil {
		return err
	}

	newUser = UserInfo{
		Name:        "tt2@tt.com",
		FirstName:   "tt2",
		LastName:    "tt2",
		Email:       "tt2@tt.com",
		Confirmed:   true,
		PasswordSha: EncodePassword("123qwe"),
		SchoolId:    schoolId,
		Role:        UserRoleTeacher,
	}

	err = createTeacher(db, newUser)
	if err != nil {
		return err
	}

	classId, _, err := CreateClass(db, schoolId, "tt@tt.com", "math", 1)
	if err != nil {
		return
	}

	path := PathId{
		schoolId:  schoolId,
		classId:   classId,
		teacherId: "tt@tt.com",
	}

	newUser = UserInfo{
		Name:        "ss@ss.com",
		FirstName:   "ss",
		LastName:    "ss",
		Email:       "ss@ss.com",
		Confirmed:   true,
		PasswordSha: EncodePassword("123qwe"),
		SchoolId:    schoolId,
		Role:        UserRoleStudent,
		Job:         getJobId(db, KeyJobs),
	}

	err = createStudent(db, newUser, path)
	if err != nil {
		return err
	}

	newUser = UserInfo{
		Name:        "ss1@ss.com",
		FirstName:   "ss1",
		LastName:    "ss1",
		Email:       "ss1@ss.com",
		Confirmed:   true,
		PasswordSha: EncodePassword("123qwe"),
		SchoolId:    schoolId,
		Role:        UserRoleStudent,
		Job:         getJobId(db, KeyJobs),
	}

	err = createStudent(db, newUser, path)
	if err != nil {
		return err
	}

	newUser = UserInfo{
		Name:        "ss2@ss.com",
		FirstName:   "ss2",
		LastName:    "ss2",
		Email:       "ss2@ss.com",
		Confirmed:   true,
		PasswordSha: EncodePassword("123qwe"),
		SchoolId:    schoolId,
		Role:        UserRoleStudent,
		Job:         getJobId(db, KeyCollegeJobs),
		College:     true,
	}

	err = createStudent(db, newUser, path)
	if err != nil {
		return err
	}

	newUser = UserInfo{
		Name:        "ss3@ss.com",
		FirstName:   "ss3",
		LastName:    "ss3",
		Email:       "ss3@ss.com",
		Confirmed:   true,
		PasswordSha: EncodePassword("123qwe"),
		SchoolId:    schoolId,
		Role:        UserRoleStudent,
		Job:         getJobId(db, KeyCollegeJobs),
		College:     true,
	}

	err = createStudent(db, newUser, path)
	if err != nil {
		return err
	}

	return

}
