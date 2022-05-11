package main

import (
	"context"
	"fmt"
	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/lgr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"
)

type TestClock struct {
	Current time.Time
}

func (t *TestClock) Now() time.Time {
	if t.Current.IsZero() {
		t.Current, _ = time.Parse(time.RFC3339, "2020-01-02T15:04:05Z")
	}
	t.Tick()
	lgr.Printf("DEBUG current time - %v", t.Current)
	return t.Current
}
func (t *TestClock) Tick() {
	t.Current = t.Current.Add(time.Millisecond)
}
func (t *TestClock) TickOne(d time.Duration) {
	t.Current = t.Current.Add(d)
}

func OpenTestDB(suffix string) (db *bolt.DB, teardown func()) {
	_ = os.Mkdir("testdata", 0755)
	ldb, err := bolt.Open("testdata/db"+suffix+".db", 0666, nil)
	if err != nil {
		lgr.Printf("FATAL cannot open %v", err)
		return nil, func() {}
	}

	return ldb, func() {
		ldb.Close()
		os.Remove("testdata/db" + suffix + ".db")
	}
}

var testLoginUser string

func SetTestLoginUser(username string) {
	testLoginUser = username
}
func InitTestServer(port int, db *bolt.DB, userName string) (teardown func()) {
	SetTestLoginUser(userName)
	mux := createRouter(db)
	mux.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			user := token.User{
				Name: testLoginUser,
			}
			ctx := request.Context()
			ctx = context.WithValue(ctx, "user", user)
			handler.ServeHTTP(writer, request.WithContext(ctx))
		})
	})
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	l, _ := net.Listen("tcp", addr)
	ts := httptest.NewUnstartedServer(mux)

	ts.Listener = l
	ts.Start()
	return func() {
		ts.Close()
	}
}

func FullStartTestServer(dbSuffix string, port int, userName string) (db *bolt.DB, teardown func()) {
	db, dbTearDown := OpenTestDB(dbSuffix)
	InitDefaultAccounts(db)
	netTearDown := InitTestServer(port, db, userName)

	return db, func() {
		dbTearDown()
		netTearDown()
	}

}

// creates entities
// classes contains noClasses for every teacher in every school + 4 mandatory classes for every school
func CreateTestAccounts(db *bolt.DB, noSchools, noTeachers, noClasses, noStudents int) (admins, schools, teachers, classes, students []string, errE error) {

	schools = make([]string, 0)
	admins = make([]string, 0)
	teachers = make([]string, 0)
	students = make([]string, 0)
	classes = make([]string, 0)
	mandatoryClasses := make([]string, 0)

	for s := 0; s < noSchools; s++ {
		schoolId, err := FindOrCreateSchool(db, fmt.Sprintf("scool %d", s), "sc, ca", s)
		if err != nil {
			errE = err
			return
		}
		schools = append(schools, schoolId)
		admin := UserInfo{
			Name:        fmt.Sprintf("test%d@admin.com", s),
			Email:       fmt.Sprintf("test%d@admin.com", s),
			PasswordSha: EncodePassword("123qwe"),
			FirstName:   "test",
			LastName:    "admin",
			Role:        UserRoleAdmin,
			SchoolId:    schoolId,
		}

		_, errE = CreateSchoolAdmin(db, admin)
		if errE != nil {
			lgr.Printf("ERROR school admin is not created: %v", err)
			return
		}
		admins = append(admins, admin.Name)

		schoolClasses := getSchoolClasses(db, schoolId)
		if schoolClasses != nil {
			for _, class := range schoolClasses {
				mandatoryClasses = append(mandatoryClasses, class.Id)
			}
		}

		for t := 0; t < noTeachers; t++ {
			teacher := UserInfo{
				Name:        fmt.Sprintf("test%d@teacher.com", s*noTeachers+t),
				Email:       fmt.Sprintf("test%d@teacher.com", s*noTeachers+t),
				PasswordSha: EncodePassword("123qwe"),
				FirstName:   "test",
				LastName:    "admin",
				Role:        UserRoleTeacher,
				SchoolId:    schoolId,
			}
			errE = createTeacher(db, teacher)
			if errE != nil {
				return
			}
			teachers = append(teachers, teacher.Name)

			for class := 0; class < noClasses; class++ {

				classId, _, err := CreateClass(db, schoolId, teacher.Name,
					RandomString(10), class)
				if err != nil {
					errE = err
					return
				}

				classes = append(classes, classId)

				for st := 0; st < noStudents; st++ {
					studentId := fmt.Sprintf("%s@student.com", RandomString(15))
					student := UserInfo{
						Name:        studentId,
						Email:       studentId,
						PasswordSha: EncodePassword("123qwe"),
						FirstName:   "test",
						LastName:    "admin",
						Role:        UserRoleStudent,
						SchoolId:    schoolId,
					}
					errE = createStudent(db, student, PathId{
						schoolId:  schoolId,
						teacherId: teacher.Email,
						classId:   classId,
					})
					if errE != nil {
						return
					}
					students = append(students, studentId)
				}
			}
		}
	}

	classes = append(classes, mandatoryClasses...)

	return
}

func TestFirst(t *testing.T) {

	db, teardown := OpenTestDB("")
	defer teardown()

	if db == nil {
		t.Fatalf("db not opened")
	}
}

func TestSchema(t *testing.T) {
	db, teardown := OpenTestDB("")
	defer teardown()

	InitDefaultAccounts(db)

	ok, err := checkUserInLocalStore(db, "test@admin.com", "123qwe")

	assert.True(t, ok, err)
}

func TestIntegrationAuth(t *testing.T) {

	db, teardown := OpenTestDB("-integration")
	defer teardown()

	InitDefaultAccounts(db)
	auth := initAuth(db, ServerConfig{})

	mux := createRouter(db)

	m := auth.Middleware()
	mux.Use(buildAuthMiddleware(m))
	authRoute, _ := auth.Handlers()
	mux.Handle("/auth/al/login", authRoute)

	l, _ := net.Listen("tcp", "127.0.0.1:8089")

	ts := httptest.NewUnstartedServer(mux)
	assert.NoError(t, ts.Listener.Close())
	ts.Listener = l
	ts.Start()
	defer func() {
		ts.Close()
	}()

	formData := url.Values{
		"user":   {"test@admin.com"},
		"passwd": {"123qwe"},
	}
	resp, err := http.PostForm("http://127.0.0.1:8089/auth/al/login", formData)

	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	resp, err = http.Get("http://127.0.0.1:8089/auth/al/login?user=test@admin.com&passwd=123qwe")

	require.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	resp, err = http.Get("http://127.0.0.1:8089/api/classes/teachers")

	assert.NotNil(t, resp)
	assert.Equal(t, 401, resp.StatusCode)

}

func TestIntegrationLoginPage(t *testing.T) {
	db, teardown := OpenTestDB("-integration")
	defer teardown()

	InitDefaultAccounts(db)
	mux := createRouter(db)
	l, _ := net.Listen("tcp", "127.0.0.1:8089")

	ts := httptest.NewUnstartedServer(mux)
	assert.NoError(t, ts.Listener.Close())
	ts.Listener = l
	ts.Start()
	defer func() {
		ts.Close()
	}()
}

func TestInitialDB(t *testing.T) {
	db, teardown := OpenTestDB("inidb")
	defer teardown()

	noSchools := 2
	noTeachers := 3
	noClasses := 4
	noStudents := 5

	a, s, tt, c, st, err := CreateTestAccounts(db, noSchools, noTeachers, noClasses, noStudents)
	require.Nil(t, err)
	require.Equal(t, noSchools, len(s))
	require.Equal(t, noSchools, len(a))
	require.Equal(t, noSchools*noTeachers, len(tt))
	require.Equal(t, noSchools*noTeachers*noClasses+noSchools*4, len(c))
	require.Equal(t, noSchools*noTeachers*noClasses*noStudents, len(st))

	store, err := getUserInLocalStore(db, st[0])

	require.Nil(t, err)
	require.Equal(t, st[0], store.Name)

	teach, err := getUserInLocalStore(db, tt[0])
	require.Nil(t, err)
	require.Equal(t, tt[0], teach.Name)

	err = db.View(func(tx *bolt.Tx) error {
		schools := tx.Bucket([]byte(KeySchools))
		for i, si := range s {
			school := schools.Bucket([]byte(si))

			if school == nil {
				return fmt.Errorf("school %d not found", i)
			}

			teachers := school.Bucket([]byte(KeyTeachers))
			for j, ti := range tt[i*noTeachers : (i+1)*noTeachers] {
				teachInfo, err := getUserInLocalStoreTx(tx, ti)
				if err != nil {
					return fmt.Errorf("teacher %d in school %d is not found", j, i)
				}

				if teachInfo.SchoolId != si {
					return fmt.Errorf("mismatch school id for teacher %d", j)
				}
				teach := teachers.Bucket([]byte(ti))
				if teach == nil {
					return fmt.Errorf("teachers bucket not found for %d in school %d", j, i)
				}
			}

		}
		return nil
	})

	require.Nil(t, err)
}

func TestBackupSecured(t *testing.T) {

	db, teardown := OpenTestDB("-integration")
	defer teardown()

	InitDefaultAccounts(db)
	auth := initAuth(db, ServerConfig{
		AdminPassword: "test1",
	})
	mux := createRouter(db)

	m := auth.Middleware()
	mux.Use(buildAuthMiddleware(m))
	mux.Handle("/admin/backup", backUpHandler(db))

	l, _ := net.Listen("tcp", "127.0.0.1:8089")

	ts := httptest.NewUnstartedServer(mux)
	assert.NoError(t, ts.Listener.Close())
	ts.Listener = l
	ts.Start()
	defer func() {
		ts.Close()
	}()

	client := &http.Client{}

	// access allowed
	req, _ := http.NewRequest(http.MethodGet,
		"http://localhost:8089/admin/backup",
		nil)
	req.Header.Add("Authorization", "Basic YWRtaW46dGVzdDE=")
	get, err := client.Do(req)
	require.Nil(t, err, fmt.Sprintf("backup failed: %v", err))
	assert.NotNil(t, get)
	assert.Equal(t, 200, get.StatusCode)

	// wrong password
	req, _ = http.NewRequest(http.MethodGet,
		"http://localhost:8089/admin/backup",
		nil)
	req.Header.Add("Authorization", "Basic YWRtaW46dGVzdDEx")
	get, err = client.Do(req)
	require.Nil(t, err, fmt.Sprintf("backup failed: %v", err))
	assert.NotNil(t, get)
	assert.Equal(t, 401, get.StatusCode)

	// malformed password
	req, _ = http.NewRequest(http.MethodGet,
		"http://localhost:8089/admin/backup",
		nil)
	req.Header.Add("Authorization", "Basic malformed")
	get, err = client.Do(req)
	require.Nil(t, err, fmt.Sprintf("backup failed: %v", err))
	assert.NotNil(t, get)
	assert.Equal(t, 401, get.StatusCode)

}
