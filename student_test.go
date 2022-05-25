package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/lgr"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

func Test_addUbuck2Student(t *testing.T) {
	clock := AppClock{}
	db, dbTearDown := OpenTestDB("")
	defer dbTearDown()

	_, schools, _, _, students, err := CreateTestAccounts(db, 2, 2, 2, 3)

	require.Nil(t, err)
	require.Equal(t, 24, len(students))

	for _, s := range students {
		userInfo, _ := getUserInLocalStore(db, s)
		err := addUbuck2Student(db, &clock, userInfo, decimal.NewFromFloat(1.01), "daily payment")
		require.Nil(t, err)
	}

	var balance decimal.Decimal

	var studentNetWo decimal.Decimal

	_ = db.View(func(tx *bolt.Tx) error {
		cb := tx.Bucket([]byte(KeyCB))
		accounts := cb.Bucket([]byte(KeybAccounts))
		ub := accounts.Bucket([]byte(CurrencyUBuck))
		v := ub.Get([]byte(KeyBalance))

		_ = balance.UnmarshalText(v)

		studentNetWo = StudentNetWorthTx(tx, students[0])
		return nil
	})

	assert.Equal(t, -12.12, balance.InexactFloat64())
	assert.Equal(t, 1.01, studentNetWo.InexactFloat64())
}

//func TestStudentNetWorthTx(t *testing.T) {
//	type args struct {
//		tx       *bolt.Tx
//		userName string
//	}
//	tests := []struct {
//		name    string
//		args    args
//		wantRes decimal.Decimal
//	}{
//		{name: , args: , wantRes: },
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			assert.Equalf(t, tt.wantRes, StudentNetWorthTx(tt.args.tx, tt.args.userName), "StudentNetWorthTx(%v, %v)", tt.args.tx, tt.args.userName)
//		})
//	}
//}

func TestDailyPayment(t *testing.T) {

	lgr.Printf("INFO TestDailyPayment")
	t.Log("INFO TestDailyPayment")
	clock := TestClock{}
	db, dbTearDown := OpenTestDB("-pay")
	defer dbTearDown()
	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	student, _ := getUserInLocalStore(db, students[0])

	r := DailyPayIfNeeded(db, &clock, student)

	require.True(t, r)

	require.Equal(t, float64(student.Income), StudentNetWorth(db, student.Name).InexactFloat64())

	r = DailyPayIfNeeded(db, &clock, student)
	require.False(t, r)

	clock.Tick()
	r = DailyPayIfNeeded(db, &clock, student)
	require.False(t, r)

	clock.TickOne(24 * time.Hour)
	r = DailyPayIfNeeded(db, &clock, student)
	require.True(t, r)

	netWorth := decimal.Zero
	_ = db.View(func(tx *bolt.Tx) error {
		netWorth = StudentNetWorthTx(tx, students[0])
		return nil
	})

	require.True(t, netWorth.GreaterThan(decimal.NewFromInt(200)))
}

func TestStudentAddClass_Teachers(t *testing.T) {
	db, tearDown := FullStartTestServer("studentAddClass_Teachers", 8090, "test@admin.com")
	defer tearDown()
	_, schools, teachers, _, students, err := CreateTestAccounts(db, 2, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	var classAddCode string

	_ = db.View(func(tx *bolt.Tx) error {
		school, _ := SchoolByIdTx(tx, schools[0])
		teachersBucket := school.Bucket([]byte(KeyTeachers))
		teacher := teachersBucket.Bucket([]byte(teachers[0]))
		classesBucket := teacher.Bucket([]byte(KeyClasses))
		c := classesBucket.Cursor()
		k, _ := c.First()
		class := classesBucket.Bucket(k)
		classAddCode = string(class.Get([]byte(KeyAddCode)))
		return nil
	})

	body := openapi.RequestAddClass{
		Id:      students[0],
		AddCode: classAddCode,
	}

	marshal, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8090/api/classes/addClass", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v []openapi.ResponseMemberClass
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
}

func TestStudentAddClass_Schools(t *testing.T) {
	db, tearDown := FullStartTestServer("studentAddClass_Schools", 8090, "test@admin.com")
	defer tearDown()
	_, schools, _, _, students, err := CreateTestAccounts(db, 2, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])
	var freshmanAddCode string

	_ = db.View(func(tx *bolt.Tx) error {
		school, _ := SchoolByIdTx(tx, schools[0])
		classes := school.Bucket([]byte(KeyClasses))
		c := classes.Cursor()
		k, _ := c.First()
		class := classes.Bucket(k)
		freshmanAddCode = string(class.Get([]byte(KeyAddCode)))
		return nil
	})

	// initialize http client
	client := &http.Client{}

	body := openapi.RequestAddClass{
		Id:      students[0],
		AddCode: freshmanAddCode,
	}

	marshal, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8090/api/classes/addClass", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v []openapi.ResponseMemberClass
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)

	assert.Equal(t, 2, len(v))
}

func TestStudentAddClass_InvalidCode(t *testing.T) {
	db, tearDown := FullStartTestServer("studentAddClass_InvalidCode", 8090, "test@admin.com")
	defer tearDown()
	_, _, _, _, students, err := CreateTestAccounts(db, 2, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	body := openapi.RequestAddClass{
		Id:      students[0],
		AddCode: "invalid1",
	}

	marshal, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8090/api/classes/addClass", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 404, resp.StatusCode, resp)

	var v openapi.ResponseRegister4
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
	require.Equal(t, "Invalid Add Code", v.Message)
}
