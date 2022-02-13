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

	_, _, teachers, _, _, err := CreateTestAccounts(db, 1, 1, 2, 6)

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

	require.Equal(t, 12, len(data))
}
