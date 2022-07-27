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

//if the following test is failing and you just ran the spec then it is probably because student2student has been change
//to omitempty on openapi.Settings. This causes nothing to be sent back because false is seen as empty and is omitted.
func TestGetSettings(t *testing.T) {
	db, tearDown := FullStartTestServer("getSettings", 8090, "")
	defer tearDown()

	admins, _, _, _, _, _ := CreateTestAccounts(db, 1, 2, 1, 3)

	SetTestLoginUser(admins[0])

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/settings",
		nil)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data openapi.Settings
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.True(t, data.Student2student)

	admin, err := getUserInLocalStore(db, admins[0])
	require.Nil(t, err)

	setSettings(db, admin, openapi.Settings{Student2student: false})
	settings, err := getSettings(db, admin)
	require.Nil(t, err)
	require.False(t, settings.Student2student)

	resp, err = client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	decoder = json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.False(t, data.Student2student)
}

func TestSetSettings(t *testing.T) {
	db, tearDown := FullStartTestServer("setSettings", 8090, "")
	defer tearDown()

	admins, _, _, _, _, _ := CreateTestAccounts(db, 1, 2, 1, 3)

	SetTestLoginUser(admins[0])

	client := &http.Client{}

	admin, err := getUserInLocalStore(db, admins[0])
	require.Nil(t, err)

	settings, err := getSettings(db, admin)
	require.Nil(t, err)
	require.True(t, settings.Student2student)

	settings = openapi.Settings{
		Student2student: false,
	}

	marshal, err := json.Marshal(settings)

	req, _ := http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8090/api/settings",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	settings, err = getSettings(db, admin)
	require.Nil(t, err)
	require.False(t, settings.Student2student)

	settings = openapi.Settings{
		Student2student: true,
	}

	marshal, err = json.Marshal(settings)

	req, _ = http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8090/api/settings",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	settings, err = getSettings(db, admin)
	require.Nil(t, err)
	require.True(t, settings.Student2student)
}

func TestStudent2Student(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("payTransactions_student", 8090, "")
	defer tearDown()

	admins, _, _, _, students, err := CreateTestAccounts(db, 1, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	client := &http.Client{}
	body := openapi.RequestPayTransaction{
		OwnerId:     students[0],
		Description: "student2student",
		Amount:      100,
		Student:     students[1],
	}

	for _, student := range students {
		userDetails, err := getUserInLocalStore(db, student)
		require.Nil(t, err)
		err = addUbuck2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), "pre load")
		require.Nil(t, err)
	}

	marshal, _ := json.Marshal(body)

	admin, err := getUserInLocalStore(db, admins[0])
	require.Nil(t, err)

	settings, err := getSettings(db, admin)
	require.Nil(t, err)
	require.True(t, settings.Student2student)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/transactions/payTransaction",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	setSettings(db, admin, openapi.Settings{Student2student: false})
	settings, err = getSettings(db, admin)
	require.Nil(t, err)
	require.False(t, settings.Student2student)

	req, _ = http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/transactions/payTransaction",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 400, resp.StatusCode)

}
