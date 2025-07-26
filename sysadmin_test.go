// sysadmin_endpoints_test.go
package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

// small helper to drop a Sys‑Admin straight into Bolt
func seedSysAdmin(t *testing.T, db *bolt.DB, email string) {
	t.Helper()
	err := db.Update(func(tx *bolt.Tx) error {
		return createSysAdminTx(tx, email, EncodePassword("pw"), "Root", "Admin")
	})
	require.NoError(t, err)
}

// -----------------------------------------------------------------------------
// /api/schools  (GET)
// -----------------------------------------------------------------------------
func TestMessageAll(t *testing.T) {
	db, tearDown := FullStartTestServer("MessageAll", 8088, "")
	defer tearDown()

	// 1 school, 1 teacher, 1 classes, 4 students
	_, _, _, _, _, err := CreateTestAccounts(db, 1, 1, 1, 4)
	require.NoError(t, err)

	admin := "admin@example.com"

	seedSysAdmin(t, db, admin)
	SetTestLoginUser(admin)

	client := &http.Client{Timeout: 2 * time.Second}

	body := openapi.RequestMessage{
		Message: "Hello World",
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:8088/api/message", bytes.NewBuffer(b))

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestMessageAllSchool(t *testing.T) {
	db, tearDown := FullStartTestServer("MessageAllSchool", 8088, "")
	defer tearDown()

	// 1 school, 1 teacher, 1 classes, 4 students
	_, schools, _, _, _, err := CreateTestAccounts(db, 1, 1, 1, 4)
	require.NoError(t, err)

	admin := "admin@example.com"

	seedSysAdmin(t, db, admin)
	SetTestLoginUser(admin)

	client := &http.Client{Timeout: 2 * time.Second}

	body := openapi.RequestMessage{
		Message:  "Hello World",
		SchoolId: schools[0],
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:8088/api/message/school", bytes.NewBuffer(b))

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestMessageAllStaff(t *testing.T) {
	db, tearDown := FullStartTestServer("MessageAllStaff", 8088, "")
	defer tearDown()

	// 1 school, 1 teacher, 1 classes, 4 students
	_, _, _, _, _, err := CreateTestAccounts(db, 1, 1, 1, 4)
	require.NoError(t, err)

	admin := "admin@example.com"

	seedSysAdmin(t, db, admin)
	SetTestLoginUser(admin)

	client := &http.Client{Timeout: 2 * time.Second}

	body := openapi.RequestMessage{
		Message: "Hello World",
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:8088/api/message/staff", bytes.NewBuffer(b))

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestMessageAllStudents(t *testing.T) {
	db, tearDown := FullStartTestServer("MessageAllStudents", 8088, "")
	defer tearDown()

	// 1 school, 1 teacher, 1 classes, 4 students
	_, _, _, _, _, err := CreateTestAccounts(db, 1, 1, 1, 4)
	require.NoError(t, err)

	admin := "admin@example.com"

	seedSysAdmin(t, db, admin)
	SetTestLoginUser(admin)

	client := &http.Client{Timeout: 2 * time.Second}

	body := openapi.RequestMessage{
		Message: "Hello World",
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:8088/api/message/students", bytes.NewBuffer(b))

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestMessageAllSchoolStaff(t *testing.T) {
	db, tearDown := FullStartTestServer("MessageAllSchoolStaff", 8088, "")
	defer tearDown()

	// 1 school, 1 teacher, 1 classes, 4 students
	_, schools, _, _, _, err := CreateTestAccounts(db, 1, 1, 1, 4)
	require.NoError(t, err)

	admin := "admin@example.com"

	seedSysAdmin(t, db, admin)
	SetTestLoginUser(admin)

	client := &http.Client{Timeout: 2 * time.Second}

	body := openapi.RequestMessage{
		Message:  "Hello World",
		SchoolId: schools[0],
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:8088/api/message/school/staff", bytes.NewBuffer(b))

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestMessageAllSchoolStudents(t *testing.T) {
	db, tearDown := FullStartTestServer("MessageAllSchoolStudents", 8088, "")
	defer tearDown()

	// 1 school, 1 teacher, 1 classes, 4 students
	_, schools, _, _, _, err := CreateTestAccounts(db, 1, 1, 1, 4)
	require.NoError(t, err)

	admin := "admin@example.com"

	seedSysAdmin(t, db, admin)
	SetTestLoginUser(admin)

	client := &http.Client{Timeout: 2 * time.Second}

	body := openapi.RequestMessage{
		Message:  "Hello World",
		SchoolId: schools[0],
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:8088/api/message/school/students", bytes.NewBuffer(b))

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestMessageUser(t *testing.T) {
	db, tearDown := FullStartTestServer("MessageUser", 8088, "")
	defer tearDown()

	// 1 school, 1 teacher, 1 classes, 4 students
	_, _, _, _, students, err := CreateTestAccounts(db, 1, 1, 1, 4)
	require.NoError(t, err)

	admin := "admin@example.com"

	seedSysAdmin(t, db, admin)
	SetTestLoginUser(admin)

	client := &http.Client{Timeout: 2 * time.Second}

	body := openapi.RequestMessage{
		Message: "Hello World",
		UserId:  students[0],
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:8088/api/message/user", bytes.NewBuffer(b))

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestSysAdmin_GetSchools(t *testing.T) {
	db, tearDown := FullStartTestServer("GetSchools", 8088, "")
	defer tearDown()

	// create 2 schools, 2 teachers, 2 classes, 2 students
	_, schools, teachers, _, _, err := CreateTestAccounts(db, 2, 2, 2, 2)
	require.NoError(t, err)
	require.Len(t, schools, 2)

	// log in as Sys‑Admin
	admin := "admin@example.com"
	seedSysAdmin(t, db, admin)
	SetTestLoginUser(admin)

	client := &http.Client{Timeout: 2 * time.Second}
	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8088/api/schools", nil)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	var payload []openapi.School
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&payload))
	assert.Len(t, payload, 2+1) // 2 schools + 1 default school that is created by running FullStartTestServer

	// negative: teacher should get 401
	SetTestLoginUser(teachers[0])
	req2, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8088/api/schools", nil)
	resp2, err := client.Do(req2)
	require.NoError(t, err)
	assert.Equal(t, 401, resp2.StatusCode)
}

// -----------------------------------------------------------------------------
// /api/schools/users  (GET ?schoolId=)
// -----------------------------------------------------------------------------
func TestSysAdmin_GetSchoolUsers(t *testing.T) {
	db, tearDown := FullStartTestServer("GetSchoolUsers", 8088, "")
	defer tearDown()

	// 1 school, 1 teacher, 1 classes, 4 students
	_, schools, teachers, _, students, err := CreateTestAccounts(db, 1, 1, 1, 4)
	require.NoError(t, err)

	admin := "admin@example.com"
	seedSysAdmin(t, db, admin)
	SetTestLoginUser(admin)

	client := &http.Client{Timeout: 2 * time.Second}
	url := "http://127.0.0.1:8088/api/schools/users?_id=" + schools[0]
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	var list []openapi.UserNoHistory
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&list))
	assert.Len(t, list, len(students)+len(teachers)+1) //+1 for the default admin that is created when a school is created

	// teacher trying the same → 401
	SetTestLoginUser(teachers[0])
	reqBad, _ := http.NewRequest(http.MethodGet, url, nil)
	respBad, err := client.Do(reqBad)
	require.NoError(t, err)
	assert.Equal(t, 401, respBad.StatusCode)
}

// -----------------------------------------------------------------------------
// /api/impersonate  (POST)
// -----------------------------------------------------------------------------
func TestSysAdmin_Impersonate(t *testing.T) {
	db, tearDown := FullStartTestServer("Impersonate", 8088, "")
	defer tearDown()

	// 1 teacher, 1 student
	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.NoError(t, err)

	admin := "admin@example.com"
	seedSysAdmin(t, db, admin)
	SetTestLoginUser(admin)

	client := &http.Client{Timeout: 2 * time.Second}

	// happy path
	body := openapi.RequestImpersonate{UserId: students[0]}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:8088/api/impersonate", bytes.NewBuffer(b))

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	var ok openapi.ResponseImpersonate
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&ok))
	assert.NotEmpty(t, ok.Token)
	assert.NotEmpty(t, ok.Xsrf)

	// teacher cannot impersonate – gets 401
	SetTestLoginUser(teachers[0])
	reqBad, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:8088/api/impersonate", bytes.NewBuffer(b))
	respBad, err := client.Do(reqBad)
	require.NoError(t, err)
	assert.Equal(t, 401, respBad.StatusCode)
}

func TestSysAdmin_DeleteSchool(t *testing.T) {
	db, tearDown := FullStartTestServer("DeleteSchool", 8088, "")
	defer tearDown()

	// 1 school, 1 teacher, 1 classes, 4 students
	_, schools, teachers, _, _, err := CreateTestAccounts(db, 1, 1, 1, 4)
	require.NoError(t, err)

	admin := "admin@example.com"
	seedSysAdmin(t, db, admin)

	client := &http.Client{Timeout: 2 * time.Second}

	// teacher cannot delete school – gets 401
	SetTestLoginUser(teachers[0])
	reqBad, _ := http.NewRequest(http.MethodDelete, "http://127.0.0.1:8088/api/schools/school?_id="+schools[0], nil)
	respBad, err := client.Do(reqBad)
	require.NoError(t, err)
	assert.Equal(t, 401, respBad.StatusCode)

	SetTestLoginUser(admin)
	// happy path
	req, _ := http.NewRequest(http.MethodDelete, "http://127.0.0.1:8088/api/schools/school?_id="+schools[0], nil)
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)
}

func TestSysAdmin_MakeSchool(t *testing.T) {
	db, tearDown := FullStartTestServer("MakeSchool", 8088, "")
	defer tearDown()

	admin := "admin@example.com"
	seedSysAdmin(t, db, admin)
	SetTestLoginUser(admin)

	client := &http.Client{Timeout: 2 * time.Second}

	// happy path
	body := openapi.RequestMakeSchool{
		School:    "New School",
		FirstName: "Admin",
		LastName:  "Super",
		Email:     "aa@aa.com",
		City:      "town",
		Zip:       97554,
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:8088/api/schools/school", bytes.NewBuffer(b))

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, 200, resp.StatusCode)

	var ok openapi.ResponseResetPassword
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&ok))
	assert.NotEmpty(t, ok.Password)
}
