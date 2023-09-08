package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllSchoolApiServiceImpl_AddCodeClass(t *testing.T) {
	db, tearDown := FullStartTestServer("addCode", 8090, "test@admin.com")
	defer tearDown()

	_, _, teachers, classes, _, _ := CreateTestAccounts(db, 2, 2, 2, 2)

	SetTestLoginUser(teachers[0])

	// initialize http client
	client := &http.Client{}

	body := openapi.RequestEditClass{
		Id: classes[0],
	}

	marshal, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8090/api/classes/class/addCode", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	var v openapi.ClassWithMembers
	_ = decoder.Decode(&v)
	assert.GreaterOrEqual(t, 9, len(v.AddCode))
}

func TestRemoveClass(t *testing.T) {
	db, tearDown := FullStartTestServer("RemoveClass", 8090, "test@admin.com")
	defer tearDown()
	members := 2
	_, _, _, classes, students, err := CreateTestAccounts(db, 1, 2, 2, members)

	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	body := openapi.RequestKickClass{
		Id:     classes[0],
		KickId: students[0],
	}

	marshal, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8090/api/classes/removeAdmin", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v []openapi.ResponseMemberClass
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)

	require.Nil(t, err)

	assert.Equal(t, 0, len(v))
}

func TestSearchMyClasses(t *testing.T) {
	db, tearDown := FullStartTestServer("searchMyClasses", 8090, "test@admin.com")
	defer tearDown()
	_, _, _, _, students, err := CreateTestAccounts(db, 1, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	body := openapi.RequestUser{
		Id: students[0],
	}

	marshal, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8090/api/classes/member", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v []openapi.ResponseMemberClass
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)

	assert.Equal(t, 1, len(v))
}

// changes the add code for frosh, soph, jr, sr. Students need it to register
func TestAddCodeClass_adminClass(t *testing.T) {
	db, tearDown := FullStartTestServer("addCodeClass", 8090, "test@admin.com")
	defer tearDown()
	admins, _, _, classes, _, err := CreateTestAccounts(db, 1, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(admins[0])

	// initialize http client
	client := &http.Client{}

	body := openapi.RequestUser{
		Id: classes[0],
	}

	marshal, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8090/api/classes/class/addCode", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v openapi.Class
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
	assert.NotNil(t, v.AddCode)
	assert.NotEqual(t, v.AddCode, "")
}

// changes the school code that teachers use to register
func TestAddCodeClass_SchoolClass(t *testing.T) {
	db, tearDown := FullStartTestServer("addCodeClass", 8090, "test@admin.com")
	defer tearDown()
	admins, schools, _, _, _, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.Nil(t, err)

	SetTestLoginUser(admins[0])

	// initialize http client
	client := &http.Client{}

	body := openapi.RequestUser{
		Id: schools[0],
	}

	marshal, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8090/api/classes/class/addCode", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v openapi.Class
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
	assert.NotNil(t, v.AddCode)
	assert.NotEqual(t, v.AddCode, "")
}

// changes the add code for a teachers class. The students will use this to register
func TestAddCodeClass_TeacherClass(t *testing.T) {
	db, tearDown := FullStartTestServer("addCodeClass", 8090, "test@admin.com")
	defer tearDown()
	_, _, teachers, classes, _, err := CreateTestAccounts(db, 1, 2, 2, 2)
	require.Nil(t, err)

	// SetTestLoginUser(teachers[0])
	SetTestLoginUser(teachers[0])

	// initialize http client
	client := &http.Client{}

	body := openapi.RequestUser{
		Id: classes[0],
	}

	marshal, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8090/api/classes/class/addCode", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v openapi.Class
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
	assert.NotNil(t, v.AddCode)
	assert.NotEqual(t, v.AddCode, "")
}
