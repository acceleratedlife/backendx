package main

import (
	"bytes"
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"net/url"
	"testing"

	openapi "github.com/acceleratedlife/backend/go"
	bolt "go.etcd.io/bbolt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchAdminTeacherClass(t *testing.T) {
	db, tearDown := FullStartTestServer("searchAdminTeacherClass", 8088, "")
	defer tearDown()

	admins, _, teachers, _, _, _ := CreateTestAccounts(db, 1, 2, 1, 3)

	SetTestLoginUser(admins[0])

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8088/api/classes/teachers",
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
	db, tearDown := FullStartTestServer("searchAdminTeacherClass", 8088, "")
	defer tearDown()

	admins, schools, _, _, _, _ := CreateTestAccounts(db, 1, 1, 1, 25)

	SetTestLoginUser(admins[0])

	client := &http.Client{}

	u, _ := url.ParseRequestURI("http://127.0.0.1:8088/api/schools/school/count")
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

func TestExecuteTax(t *testing.T) {
	db, tearDown := FullStartTestServer("ExecuteTax", 8088, "")
	defer tearDown()

	admins, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 10)

	SetTestLoginUser(students[0])

	client := &http.Client{}
	bodyFlat := openapi.RequestTax{
		TaxRate: 17,
	}

	marshal, _ := json.Marshal(bodyFlat)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8088/api/schools/school/tax",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 401, resp.StatusCode)

	SetTestLoginUser(teachers[0])
	req, _ = http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8088/api/schools/school/tax",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 401, resp.StatusCode)

	SetTestLoginUser(admins[0])

	marshal, _ = json.Marshal(bodyFlat)
	req, _ = http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8088/api/schools/school/tax",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	bodyProgressive := openapi.RequestTax{
		TaxRate: 0,
	}

	err = db.Update(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte(KeyUsers))
		for _, d := range students {
			studentData := users.Get([]byte(d))
			var student UserInfo
			err = json.Unmarshal(studentData, &student)
			if err != nil {
				return err
			}

			student.TaxableIncome = int32(rand.IntN(600-260) + 260)
			marshal, err := json.Marshal(student)
			if err != nil {
				return err
			}

			err = users.Put([]byte(d), marshal)
			if err != nil {
				return err
			}

		}
		return err
	})

	require.NotNil(t, resp)

	marshal, _ = json.Marshal(bodyProgressive)
	req, _ = http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8088/api/schools/school/tax",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestProgressiveBrackets(t *testing.T) {
	db, tearDown := FullStartTestServer("ProgressiveBrackets", 8088, "")
	defer tearDown()

	admins, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 5)

	SetTestLoginUser(admins[0])

	err := db.Update(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte(KeyUsers))
		for _, d := range students {
			studentData := users.Get([]byte(d))
			var student UserInfo
			err := json.Unmarshal(studentData, &student)
			if err != nil {
				return err
			}

			student.TaxableIncome = int32(rand.IntN(600-260) + 260)
			marshal, err := json.Marshal(student)
			if err != nil {
				return err
			}

			err = users.Put([]byte(d), marshal)
			if err != nil {
				return err
			}

		}
		return nil
	})
	require.Nil(t, err)

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8088/api/schools/school/tax",
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}
