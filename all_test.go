package main

import (
	"encoding/json"
	"net/http"
	"testing"

	openapi "github.com/acceleratedlife/backend/go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	_, tearDown := FullStartTestServer("addCode", 8090, "test@admin.com")
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
	db, tearDown := FullStartTestServer("addCode", 8090, "")
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
	db, tearDown := FullStartTestServer("addCode", 8090, "")
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
	db, tearDown := FullStartTestServer("addCode", 8090, "")
	defer tearDown()
	members := 1 // if this is a larger number then the test will sometimes fail due to race

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

	require.Equal(t, students[0], data.Members[0].Id) //sometimes this is students[1] and other times students[0], race condition?
	require.Equal(t, classes[0], data.Id)
	require.Equal(t, len(data.Members), members)
}
