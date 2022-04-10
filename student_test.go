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

	_, _, _, _, students, err := CreateTestAccounts(db, 2, 2, 2, 3)

	require.Nil(t, err)
	require.Equal(t, 24, len(students))

	for _, s := range students {
		err := addUbuck2Student(db, &clock, s, decimal.NewFromFloat(1.01), "daily payment")
		require.Nil(t, err)
	}

	var balance decimal.Decimal

	var studentNetWo decimal.Decimal

	_ = db.View(func(tx *bolt.Tx) error {
		cb := tx.Bucket([]byte(KeyCB))
		accounts := cb.Bucket([]byte(KeyAccounts))
		ub := accounts.Bucket([]byte(CurrencyUBuck))
		v := ub.Get([]byte(KeyBalance))

		_ = balance.UnmarshalText(v)

		studentNetWo = StudentNetWorthTx(tx, students[0])
		return nil
	})

	assert.Equal(t, decimal.NewFromFloat(-24.24), balance)
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

	require.Equal(t, 121.32, StudentNetWorth(db, student.Name).InexactFloat64())

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

func TestStudentAddClass(t *testing.T) {
	db, tearDown := FullStartTestServer("studentAddClass", 8090, "test@admin.com")
	defer tearDown()
	_, _, _, classes, students, err := CreateTestAccounts(db, 2, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/classes/class?_id="+classes[1],
		nil)

	resp, _ := client.Do(req)
	var data openapi.ClassWithMembers
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	body := openapi.RequestAddClass{
		Id:      students[0],
		AddCode: data.AddCode,
	}

	marshal, _ := json.Marshal(body)
	req, _ = http.NewRequest(http.MethodPut, "http://127.0.0.1:8090/api/classes/addClass", bytes.NewBuffer(marshal))
	resp, err = client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v []openapi.ResponseMemberClass
	decoder = json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)

	assert.Equal(t, 2, len(v))
}
