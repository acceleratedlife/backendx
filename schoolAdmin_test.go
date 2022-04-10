package main

import (
	"encoding/json"
	"net/http"
	"testing"

	openapi "github.com/acceleratedlife/backend/go"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchAdminTeacherClass(t *testing.T) {
	db, tearDown := FullStartTestServer("searchAdminTeacherClass", 8090, "")
	defer tearDown()

	admins, _, teachers, _, _, err := CreateTestAccounts(db, 1, 2, 1, 3)

	SetTestLoginUser(admins[0])

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/classes/teachers",
		nil)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
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
