package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"testing"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/shopspring/decimal"
	bolt "go.etcd.io/bbolt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuth(t *testing.T) {
	_, tearDown := FullStartTestServer("auth", 8090, "test@admin.com")
	defer tearDown()
	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/users/auth",
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
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
	clock := TestClock{}
	db, tearDown := FullStartTestServer("searchStudents", 8090, "")
	defer tearDown()
	coinGecko(db)

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 3)

	SetTestLoginUser(teachers[0])

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	_ = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(15000), CurrencyUBuck, "pre load")
	body := openapi.RequestCryptoConvert{
		Name: "cardano",
		Buy:  10,
		Sell: 0,
	}
	err = cryptoTransaction(db, &clock, userDetails, body)
	require.Nil(t, err)

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/users",
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data []openapi.UserNoHistory
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, 1, len(data))
	require.NotZero(t, data[0].NetWorth)
}

func TestSearchStudent(t *testing.T) {
	db, tearDown := FullStartTestServer("searchStudent", 8090, "")
	defer tearDown()
	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 1, 2, 6)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/users/user?_id="+students[0],
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data openapi.User
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, students[0], data.Id)
	require.NotNil(t, data.FirstName)
	require.NotNil(t, data.LastName)
	require.NotNil(t, data.Income)
}

func TestSearchClass(t *testing.T) {
	db, tearDown := FullStartTestServer("searchClass", 8090, "")
	defer tearDown()
	members := 10

	_, _, _, classes, students, _ := CreateTestAccounts(db, 3, 3, 3, members)

	SetTestLoginUser(students[0])

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/classes/class?_id="+classes[0],
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data openapi.ClassWithMembers
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	flag := false
	for i := range data.Members {
		if data.Members[i].Id == students[0] {
			flag = true
			break
		}
	}

	require.Equal(t, flag, true)
	require.Equal(t, classes[0], data.Id)
	require.Equal(t, len(data.Members), members)
}

func TestUserEdit(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServerClock("userEdit", 8090, "test@admin.com", &clock)
	defer tearDown()
	_, _, _, _, students, err := CreateTestAccounts(db, 1, 2, 2, 2)
	require.Nil(t, err)

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(15000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	body := openapi.RequestUserEdit{
		FirstName:        "test",
		LastName:         "user",
		Password:         "123qwe",
		College:          true,
		CareerTransition: true,
	}

	marshal, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8090/api/users/user", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v openapi.User
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)

	regex, err := regexp.Compile(string(body.LastName[0]) + "[0-9]{4}")
	require.Nil(t, err)

	assert.Equal(t, body.FirstName, v.FirstName)
	assert.True(t, regex.MatchString(v.LastName))
	assert.Equal(t, body.College, v.College)
	assert.Equal(t, body.CareerTransition, v.CareerTransition)
}

func TestUserEditStaff(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServerClock("userEditStaff", 8090, "test@admin.com", &clock)
	defer tearDown()
	_, _, teachers, _, _, err := CreateTestAccounts(db, 1, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(teachers[0])

	// initialize http client
	client := &http.Client{}

	body := openapi.RequestUserEdit{
		FirstName:        "test",
		LastName:         "user",
		Password:         "123qwe",
		College:          false,
		CareerTransition: false,
	}

	marshal, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8090/api/users/user", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v openapi.User
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)

	assert.Equal(t, body.FirstName, v.FirstName)
	assert.Equal(t, body.LastName, v.LastName)
	assert.Equal(t, body.College, v.College)
	assert.Equal(t, body.CareerTransition, v.CareerTransition)
}

func TestUserEditNegative(t *testing.T) {
	db, tearDown := FullStartTestServer("userEditNegative", 8090, "test@admin.com")
	defer tearDown()
	_, _, _, _, students, err := CreateTestAccounts(db, 1, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}
	require.Nil(t, err)

	body := openapi.RequestUserEdit{
		FirstName:        "test",
		LastName:         "user",
		Password:         "123qwe",
		College:          true,
		CareerTransition: false,
	}

	marshal, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8090/api/users/user", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v openapi.User
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)

	regex, err := regexp.Compile(string(body.LastName[0]) + "[0-9]{4}")
	require.Nil(t, err)

	assert.Equal(t, body.FirstName, v.FirstName)
	assert.True(t, regex.MatchString(v.LastName))
	assert.Equal(t, body.College, v.College)
	assert.Equal(t, body.CareerTransition, v.CareerTransition)

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	require.NotEqual(t, student.Name, "")

	var account openapi.ResponseCurrencyExchange
	err = db.View(func(tx *bolt.Tx) error {
		studentBucket, err := getStudentBucketRx(tx, student.Name)
		if err != nil {
			return err
		}
		accounts := studentBucket.Bucket([]byte(KeyAccounts))
		debt := accounts.Bucket([]byte(KeyDebt))
		account, err = getStudentAccountRx(tx, debt, student.Name)
		if err != nil {
			return err
		}
		return nil
	})

	require.Nil(t, err)

	require.Greater(t, account.Balance, float32(5000))
	require.Less(t, account.Balance, float32(8000))
}

func TestSearchStudentBucks(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("searchStudentBucks", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 3, 3, 3, 3)

	SetTestLoginUser(students[0])

	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[0], "pre load")
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/accounts/all",
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data []openapi.ResponseAccount
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	assert.Equal(t, data[0].Balance, float32(1000))
}

func TestSearchStudentBucksNegative(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("searchStudentBucksNegative", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 3, 3, 3, 3)

	SetTestLoginUser(students[0])

	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	userDetails1, err := getUserInLocalStore(db, students[1])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[1], "pre load")
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails1, decimal.NewFromFloat(1000), teachers[0], "pre load")
	require.Nil(t, err)
	err = chargeStudent(db, &clock, userDetails, decimal.NewFromFloat(1001), teachers[0], "charge", false)
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/accounts/all",
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data []openapi.ResponseAccount
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	if data[0].Id == KeyDebt {
		assert.Equal(t, float32(1001), data[0].Balance)
	} else {
		assert.Equal(t, float32(1000), data[0].Balance)
	}

	assert.Equal(t, len(data), 2)
}

func TestSearchStudentBucksUbuck(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("searchStudentBucksUbuck", 8090, "")
	defer tearDown()

	_, _, _, _, students, _ := CreateTestAccounts(db, 3, 3, 3, 3)

	SetTestLoginUser(students[0])

	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), CurrencyUBuck, "daily pay")
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/accounts/all",
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data []openapi.ResponseAccount
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	assert.Equal(t, float32(1000), data[0].Balance)
	assert.Equal(t, float32(1), data[0].Conversion)
}

func TestSearchAllBucks(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("searchAllBucks", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 3, 1, 1)

	SetTestLoginUser(students[0])

	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), CurrencyUBuck, "daily pay")
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[0], "daily pay")
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[1], "daily pay")
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[2], "daily pay")
	require.Nil(t, err)
	err = chargeStudent(db, &clock, userDetails, decimal.NewFromFloat(10000), teachers[0], "charge", false)
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/bucks",
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data []openapi.Buck
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	assert.Equal(t, "Debt", data[3].Name)
	assert.Equal(t, "UBuck", data[4].Name)
}

func TestExchangeRate(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("exchangeRate", 8090, "")
	defer tearDown()
	members := 10

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 3, 3, members)

	SetTestLoginUser(students[0])

	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), CurrencyUBuck, "daily pay")
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[0], "daily pay")
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(2000), teachers[1], "daily pay")
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[2], "daily pay")
	require.Nil(t, err)
	err = chargeStudent(db, &clock, userDetails, decimal.NewFromFloat(10000), teachers[0], "charge", false)
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/accounts/exchangeRate?sellCurrency="+teachers[0]+"&buyCurrency="+teachers[1],
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data []openapi.ResponseAccount
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, float32(1000), data[0].Balance)
	require.Equal(t, float32(2), data[0].Conversion)
	require.Equal(t, float32(2000), data[1].Balance)

	req, _ = http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/accounts/exchangeRate?sellCurrency="+CurrencyUBuck+"&buyCurrency="+teachers[1],
		nil)

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	decoder = json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, float32(1000), data[0].Balance)
	require.Equal(t, float32(1.5), data[0].Conversion)
	require.Equal(t, float32(2000), data[1].Balance)

}

func TestExchangeRate_ubuck(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("exchangeRate_ubuck", 8090, "")
	defer tearDown()
	members := 10

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 3, 3, members)

	SetTestLoginUser(students[0])

	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), CurrencyUBuck, "daily pay")
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[0], "daily pay")
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(2000), teachers[1], "daily pay")
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[2], "daily pay")
	require.Nil(t, err)
	err = chargeStudent(db, &clock, userDetails, decimal.NewFromFloat(10000), teachers[0], "charge", false)
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/accounts/exchangeRate?sellCurrency="+CurrencyUBuck+"&buyCurrency="+teachers[1],
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data []openapi.ResponseAccount
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, float32(1000), data[0].Balance)
	require.Equal(t, float32(1.5), data[0].Conversion)
	require.Equal(t, float32(2000), data[1].Balance)

}

func TestExchangeRate_debt(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("exchangeRate_debt", 8090, "")
	defer tearDown()
	members := 10

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 3, 3, members)

	SetTestLoginUser(students[0])

	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), CurrencyUBuck, "daily pay")
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[0], "daily pay")
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(2000), teachers[1], "daily pay")
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[2], "daily pay")
	require.Nil(t, err)
	err = chargeStudent(db, &clock, userDetails, decimal.NewFromFloat(10000), teachers[0], "charge", false)
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/accounts/exchangeRate?sellCurrency="+KeyDebt+"&buyCurrency="+CurrencyUBuck,
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data []openapi.ResponseAccount
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, float32(-1), data[0].Conversion)

	req, _ = http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/accounts/exchangeRate?sellCurrency="+CurrencyUBuck+"&buyCurrency="+KeyDebt,
		nil)

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	decoder = json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, float32(-1), data[0].Conversion)

	req, _ = http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/accounts/exchangeRate?sellCurrency="+KeyDebt+"&buyCurrency="+teachers[1],
		nil)

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	decoder = json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, float32(-1.5), data[0].Conversion)

	req, _ = http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/accounts/exchangeRate?sellCurrency="+teachers[1]+"&buyCurrency="+KeyDebt,
		nil)

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	decoder = json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, float32(-.666667), data[0].Conversion)

}

func TestPayTransaction_credit(t *testing.T) {
	db, tearDown := FullStartTestServer("payTransaction_credit", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 2, 2, 2, 2)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}
	body := openapi.RequestPayTransaction{
		OwnerId:     teachers[0],
		Description: "credit",
		Amount:      100,
		Student:     students[0],
	}

	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/transactions/payTransaction",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

}

func TestPayTransaction_debit(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("payTransaction_debit", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 2, 2, 2, 2)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}
	body := openapi.RequestPayTransaction{
		OwnerId:     teachers[0],
		Description: "debit",
		Amount:      -100,
		Student:     students[0],
	}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[0], "pre load")
	require.Nil(t, err)

	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/transactions/payTransaction",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

}

func TestPayTransactions_credit(t *testing.T) {
	db, tearDown := FullStartTestServer("payTransactions_credit", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 2, 2, 2)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}
	body := openapi.RequestPayTransactions{
		Owner:       teachers[0],
		Description: "credit",
		Amount:      100,
		Students:    students,
	}

	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/transactions/payTransactions",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

}

func TestPayTransactions_debit(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("payTransactions_debit", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 2, 2, 2)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}
	body := openapi.RequestPayTransactions{
		Owner:       teachers[0],
		Description: "debit",
		Amount:      -100,
		Students:    students,
	}

	for _, student := range students {
		userDetails, err := getUserInLocalStore(db, student)
		require.Nil(t, err)
		err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[0], "pre load")
		require.Nil(t, err)
	}

	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/transactions/payTransactions",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

}

func TestPayTransaction_student(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("payTransactions_student", 8090, "")
	defer tearDown()

	admins, _, _, _, students, _ := CreateTestAccounts(db, 1, 2, 2, 2)

	SetTestLoginUser(students[0])

	admin, err := getUserInLocalStore(db, admins[0])
	require.Nil(t, err)

	setSettings(db, admin, openapi.Settings{Student2student: true})

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

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/transactions/payTransaction",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

}

func TestPayTransaction_studentDebt(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("payTransactions_studentDebt", 8090, "")
	defer tearDown()

	_, _, _, _, students, _ := CreateTestAccounts(db, 1, 2, 2, 2)

	SetTestLoginUser(students[0])

	client := &http.Client{}
	body := openapi.RequestPayTransaction{
		OwnerId:     students[0],
		Description: "student2student",
		Amount:      -100,
		Student:     students[1],
	}

	for _, student := range students {
		userDetails, err := getUserInLocalStore(db, student)
		require.Nil(t, err)
		err = addUbuck2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), "pre load")
		require.Nil(t, err)
	}

	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/transactions/payTransaction",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 400, resp.StatusCode)

}

func TestMakeAuction(t *testing.T) {
	clock := TestClock{}
	db, teardown := FullStartTestServer("makeClass", 8090, "")
	defer teardown()

	_, _, teachers, classes, _, _ := CreateTestAccounts(db, 1, 1, 2, 2)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}

	body := openapi.RequestMakeAuction{
		Bid:         4,
		MaxBid:      4,
		Description: "Test Auction",
		EndDate:     clock.Now().Add(500),
		StartDate:   clock.Now(),
		OwnerId:     teachers[0],
		Visibility:  classes,
	}
	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/auctions",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

}

func TestSearchClasses(t *testing.T) {
	db, tearDown := FullStartTestServer("searchClasses", 8090, "")
	defer tearDown()
	classCount := 2

	_, _, teachers, _, _, _ := CreateTestAccounts(db, 2, 2, classCount, 2)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/classes",
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data []openapi.Class
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, classCount, len(data))
	require.Equal(t, teachers[0], data[0].OwnerId)
}

func TestDeleteAuction(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServerClock("deleteAuction", 8090, "", &clock)
	defer tearDown()

	_, _, teachers, classes, students, _ := CreateTestAccounts(db, 2, 2, 2, 2)

	SetTestLoginUser(teachers[0])

	body := openapi.RequestMakeAuction{
		Bid:         4,
		MaxBid:      4,
		Description: "Test Auction",
		EndDate:     clock.Now().Add(time.Second * 10),
		StartDate:   clock.Now(),
		OwnerId:     teachers[0],
		Visibility:  classes,
	}

	teacherDetails, _ := getUserInLocalStore(db, teachers[0])

	//to be deleted
	err := MakeAuctionImpl(db, teacherDetails, body, true)
	require.Nil(t, err)
	auctions, err := getTeacherAuctions(db, teacherDetails)
	require.Nil(t, err)

	client := &http.Client{}
	timeId := auctions[0].Id.Format(time.RFC3339Nano)

	u, _ := url.ParseRequestURI("http://127.0.0.1:8090/api/auctions/auction")
	q := u.Query()
	q.Set("_id", timeId)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodDelete,
		u.String(),
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	//to be deleted
	_ = MakeAuctionImpl(db, teacherDetails, body, true)

	userDetails, _ := getUserInLocalStore(db, students[0])
	addUbuck2Student(db, &clock, userDetails, decimal.NewFromInt32(100), "loading")

	_, err = placeBid(db, &clock, userDetails, timeId, 20)
	require.Nil(t, err)

	q.Set("_id", timeId)
	u.RawQuery = q.Encode()

	req, _ = http.NewRequest(http.MethodDelete,
		u.String(),
		nil)

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	//to be de-activated
	err = MakeAuctionImpl(db, teacherDetails, body, true)
	require.Nil(t, err)

	_, err = placeBid(db, &clock, userDetails, timeId, 20)
	require.Nil(t, err)
	clock.TickOne(time.Hour * 1)

	q.Set("_id", timeId)
	u.RawQuery = q.Encode()

	req, _ = http.NewRequest(http.MethodDelete,
		u.String(),
		nil)

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

}

func TestSearchMarketItems(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("searchMarketItems", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	SetTestLoginUser(students[0])

	userDetails, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)
	makeMarketItem(db, &clock, userDetails, openapi.RequestMakeMarketItem{
		Title: "Candy",
		Count: 4,
		Cost:  56,
	})

	client := &http.Client{}
	u, _ := url.ParseRequestURI("http://127.0.0.1:8090/api/marketItems")
	q := u.Query()
	q.Set("_id", teachers[0])
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet, u.String(), nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data []openapi.ResponseMarketItem
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, 1, len(data))
	require.Equal(t, "Candy", data[0].Title)
}
