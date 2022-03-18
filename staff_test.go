package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

func TestMakeClass(t *testing.T) {
	teacherId := "teacherName"
	addCode := RandomString(6)

	db, teardown := FullStartTestServer("", 8090, teacherId)
	defer teardown()

	schoolId, _ := FindOrCreateSchool(db, "test school", "no city", 0)
	err := db.Update(func(tx *bolt.Tx) error {
		_ = AddUserTx(tx, UserInfo{
			Name:      teacherId,
			Email:     teacherId,
			Confirmed: true,
			SchoolId:  schoolId,
			Role:      UserRoleTeacher,
		})

		school, _ := SchoolByIdTx(tx, schoolId)
		_ = school.Put([]byte("addCode"), []byte(addCode))
		teachers, _ := school.CreateBucketIfNotExists([]byte("teachers"))
		_, _ = teachers.CreateBucket([]byte(teacherId))
		//class, _ := teacher.CreateBucket([]byte(RandomString(12)))
		//_ = class.Put([]byte("name"), []byte("1"))
		//_ = class.Put([]byte("period"), itob32(13))
		//_ = class.Put([]byte("addCode"), []byte(addCode))
		return nil
	})
	assert.Nil(t, err)
	client := &http.Client{}

	body := openapi.RequestMakeClass{
		Name:    "new name",
		OwnerId: "",
		Period:  12,
	}
	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/classes/class",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	defer resp.Body.Close()
	var respData []openapi.Class
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&respData)

	assert.Equal(t, 1, len(respData))
	assert.Equal(t, 6, len(respData[0].AddCode))

}

func TestEditClass(t *testing.T) {
	db, tearDown := FullStartTestServer("addCode", 8090, "")
	defer tearDown()
	members := 2

	_, _, teachers, classes, _, err := CreateTestAccounts(db, 2, 2, 2, members)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}
	body := openapi.RequestEditClass{
		Name:   "Test Name",
		Period: 4,
		Id:     classes[0],
	}

	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8090/api/classes/class",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data openapi.Class
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, members, len(data.Members))
	require.Equal(t, "Test Name", data.Name)
	require.Equal(t, int32(4), data.Period)
	require.Equal(t, classes[0], data.Id)
}

func TestDeleteClass(t *testing.T) {
	db, tearDown := FullStartTestServer("addCode", 8090, "")
	defer tearDown()
	noClasses := 2

	_, _, teachers, classes, _, err := CreateTestAccounts(db, 2, 2, noClasses, 2)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}
	body := openapi.RequestUser{
		Id: classes[0],
	}

	marshal, err := json.Marshal(body)

	req, err := http.NewRequest(http.MethodDelete,
		"http://127.0.0.1:8090/api/classes/class",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data openapi.Class
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)
}
