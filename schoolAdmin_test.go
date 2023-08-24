package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	openapi "github.com/acceleratedlife/backend/go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchAdminTeacherClass(t *testing.T) {
	db, tearDown := FullStartTestServer("searchAdminTeacherClass", 8090, "")
	defer tearDown()

	admins, _, teachers, _, _, _ := CreateTestAccounts(db, 1, 2, 1, 3)

	SetTestLoginUser(admins[0])

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/classes/teachers",
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data openapi.ClassWithMembers
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, admins[0], data.OwnerId)
	require.Equal(t, 2, len(data.Members))
	require.Equal(t, teachers[0], data.Members[0].Id)
	require.Equal(t, teachers[1], data.Members[1].Id)
}

func TestGetStudentCountEndpoint(t *testing.T) {
	db, tearDown := FullStartTestServer("searchAdminTeacherClass", 8090, "")
	defer tearDown()

	admins, schools, _, _, _, _ := CreateTestAccounts(db, 1, 1, 1, 25)

	SetTestLoginUser(admins[0])

	client := &http.Client{}

	u, _ := url.ParseRequestURI("http://127.0.0.1:8090/api/schools/school/count")
	q := u.Query()
	q.Set("schoolId", schools[0])
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet,
		u.String(),
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var data openapi.ResponseStudentCount
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, int32(25), data.Count)
}
