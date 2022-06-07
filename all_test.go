package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/shopspring/decimal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	_, tearDown := FullStartTestServer("auth", 8090, "test@admin.com")
	defer tearDown()

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/users/auth",
		nil)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data openapi.ResponseAuth2
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, UserRoleAdmin, data.Role)
	require.Equal(t, "test@admin.com", data.Email)
	require.Equal(t, "test@admin.com", data.Id)
	require.True(t, data.IsAuth)
	require.True(t, data.IsAdmin)
	require.False(t, len(data.SchoolId) == 0)
}

func TestSearchStudents(t *testing.T) {
	db, tearDown := FullStartTestServer("searchStudents", 8090, "")
	defer tearDown()

	_, _, teachers, _, _, err := CreateTestAccounts(db, 1, 1, 1, 1)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/users",
		nil)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data []openapi.UserNoHistory
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, 1, len(data))
}

func TestSearchStudent(t *testing.T) {
	db, tearDown := FullStartTestServer("searchStudent", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 1, 2, 6)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/users/user?_id="+students[0],
		nil)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data openapi.User
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, students[0], data.Id)
	require.NotNil(t, data.FirstName)
	require.NotNil(t, data.LastName)
	require.NotNil(t, data.Income)
}

func TestSearchClass(t *testing.T) {
	db, tearDown := FullStartTestServer("searchClass", 8090, "")
	defer tearDown()
	members := 10

	_, _, _, classes, students, err := CreateTestAccounts(db, 3, 3, 3, members)

	SetTestLoginUser(students[0])

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/classes/class?_id="+classes[0],
		nil)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data openapi.ClassWithMembers
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	flag := false
	for i := range data.Members {
		if data.Members[i].Id == students[0] {
			flag = true
			break
		}
	}

	require.Equal(t, flag, true)
	require.Equal(t, classes[0], data.Id)
	require.Equal(t, len(data.Members), members)
}

func TestUserEdit(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("userEdit", 8090, "test@admin.com")
	defer tearDown()
	_, _, _, _, students, err := CreateTestAccounts(db, 1, 2, 2, 2)
	require.Nil(t, err)

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(15000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	body := openapi.UsersUserBody{
		FirstName:        "test",
		LastName:         "user",
		Password:         "123qwe",
		College:          true,
		CareerTransition: true,
	}

	marshal, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8090/api/users/user", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v openapi.User
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)

	assert.Equal(t, body.FirstName, v.FirstName)
	assert.Equal(t, body.LastName, v.LastName)
	assert.Equal(t, body.College, v.College)
	assert.Equal(t, body.CareerTransition, v.CareerTransition)
}

// The test should pass but it won't as negatives are not currently allowed.
func TestUserEditNegative(t *testing.T) {
	db, tearDown := FullStartTestServer("userEditNegative", 8090, "test@admin.com")
	defer tearDown()
	_, _, _, _, students, err := CreateTestAccounts(db, 1, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}
	require.Nil(t, err)

	body := openapi.UsersUserBody{
		FirstName:        "test",
		LastName:         "user",
		Password:         "123qwe",
		College:          true,
		CareerTransition: true,
	}

	marshal, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8090/api/users/user", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v openapi.User
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)

	assert.Equal(t, body.FirstName, v.FirstName)
	assert.Equal(t, body.LastName, v.LastName)
	assert.Equal(t, body.College, v.College)
	assert.Equal(t, body.CareerTransition, v.CareerTransition)
}

func TestSearchStudentBuck(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("searchStudentBuck", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, err := CreateTestAccounts(db, 3, 3, 3, 3)

	SetTestLoginUser(students[0])

	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[0], "pre load")
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/accounts/all",
		nil)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data []openapi.ResponseAccounts
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	assert.Equal(t, data[0].Balance, float32(1000))
}
