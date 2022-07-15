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

func TestAddTeacher(t *testing.T) {
	db, tearDown := FullStartTestServer("addTeacher", 8090, "test@admin.com")
	defer tearDown()

	schoolId, _ := FindOrCreateSchool(db, "test school", "no city", 0)

	err := db.Update(func(tx *bolt.Tx) error {
		school, _ := SchoolByIdTx(tx, schoolId)
		return school.Put([]byte("addCode"), []byte("123123"))
	})

	require.Nil(t, err)

	client := &http.Client{}

	body := openapi.RequestRegister{
		Email:     "new@teacher.com",
		Password:  "weqr",
		AddCode:   "123123",
		FirstName: "fssdf",
		LastName:  "sdfsdf11",
	}

	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/users/register",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var teach UserInfo
	err = db.View(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte("users"))
		techdata := users.Get([]byte("new@teacher.com"))

		_ = json.Unmarshal(techdata, &teach)
		return nil

	})
	require.Equal(t, UserRoleTeacher, teach.Role)
}

func TestAddStudent(t *testing.T) {
	db, tearDown := FullStartTestServer("addStudent", 8090, "test@admin.com")
	defer tearDown()

	schoolId, _ := FindOrCreateSchool(db, "test school", "no city", 0)
	_ = db.Update(func(tx *bolt.Tx) error {
		school, _ := SchoolByIdTx(tx, schoolId)
		return school.Put([]byte("addCode"), []byte("123123"))
	})
	s := UnregisteredApiServiceImpl{
		db: db,
	}

	job := Job{
		Title:       "Teacher",
		Pay:         53000,
		Description: "Teach Stuff",
		College:     true,
	}

	marshal, err := json.Marshal(job)

	err = createJobOrEvent(db, marshal, KeyCollegeJobs)
	require.Nil(t, err)

	job2 := Job{
		Title:       "Teacher",
		Pay:         53000,
		Description: "Teach Stuff",
		College:     false,
	}

	marshal, err = json.Marshal(job2)

	err = createJobOrEvent(db, marshal, KeyJobs)
	require.Nil(t, err)

	_, err = s.Register(nil, openapi.RequestRegister{
		Email:     "testaddstu@teacher.com",
		Password:  "123",
		AddCode:   "123123",
		FirstName: "qw",
		LastName:  "qw",
	})

	require.Nil(t, err)

	s1 := StaffApiServiceImpl{
		db: db,
	}
	classes, _ := s1.MakeClassImpl(UserInfo{
		Name:     "testaddstu@teacher.com",
		SchoolId: schoolId,
		Role:     UserRoleTeacher,
	}, openapi.RequestMakeClass{
		Name:    "qwe",
		OwnerId: "",
		Period:  0,
	})

	assert.Equal(t, 6, len(classes[0].AddCode))

	regResponse, err := s.Register(nil, openapi.RequestRegister{
		Email:     "testaddstu@student.com",
		Password:  "123321",
		AddCode:   classes[0].AddCode,
		FirstName: "qw1",
		LastName:  "qw2",
	})

	require.Nil(t, err)
	assert.Equal(t, 200, regResponse.Code)
}
