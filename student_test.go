package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"
)

func Test_addUbuck2Student(t *testing.T) {
	clock := TestClock{}
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
		cb, err := getCbRx(tx, schools[0])
		require.Nil(t, err)

		accounts := cb.Bucket([]byte(KeyAccounts))
		ub := accounts.Bucket([]byte(CurrencyUBuck))
		v := ub.Get([]byte(KeyBalance))

		_ = balance.UnmarshalText(v)

		studentNetWo = StudentNetWorthTx(tx, students[0])
		return nil
	})

	assert.Equal(t, -12.12, balance.InexactFloat64())
	assert.Equal(t, 1.01, studentNetWo.InexactFloat64())
}

func TestStudentAddClass_Teachers(t *testing.T) {
	db, tearDown := FullStartTestServer("studentAddClass_Teachers", 8088, "test@admin.com")
	defer tearDown()
	_, schools, teachers, _, students, err := CreateTestAccounts(db, 2, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	var classAddCode string

	_ = db.View(func(tx *bolt.Tx) error {
		school, _ := schoolByIdTx(tx, schools[0])
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
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8088/api/classes/addClass", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v []openapi.ResponseMemberClass
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
}

func TestStudentAddClass_Schools(t *testing.T) {
	db, tearDown := FullStartTestServer("studentAddClass_Schools", 8088, "test@admin.com")
	defer tearDown()
	_, schools, _, _, students, err := CreateTestAccounts(db, 2, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])
	var freshmanAddCode string

	_ = db.View(func(tx *bolt.Tx) error {
		school, _ := schoolByIdTx(tx, schools[0])
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
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8088/api/classes/addClass", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v []openapi.ResponseMemberClass
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)

	assert.Equal(t, 2, len(v))
}

func TestStudentAddClass_InvalidCode(t *testing.T) {
	db, tearDown := FullStartTestServer("studentAddClass_InvalidCode", 8088, "test@admin.com")
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
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8088/api/classes/addClass", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 404, resp.StatusCode, resp)

	var v openapi.ResponseRegister4
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
	require.Equal(t, "Invalid Add Code", v.Message)
}

func TestSearchStudentUbuck(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("searchStudentUbuck", 8088, "test@admin.com")
	defer tearDown()
	_, _, _, _, students, err := CreateTestAccounts(db, 2, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8088/api/accounts/account/student", nil)
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v openapi.ResponseSearchStudentUbuck
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
	require.Equal(t, float32(10000), v.Value)
}

func TestLatestLotto(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("LatestLotto", 8088, "test@admin.com")
	defer tearDown()
	_, _, _, _, students, err := CreateTestAccounts(db, 2, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	settings := openapi.Settings{
		Lottery: true,
		Odds:    750,
	}

	err = initializeLottery(db, userDetails, settings, &clock)
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8088/api/lottery/latest", nil)
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v openapi.ResponseLottoLatest
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
	require.Equal(t, settings.Odds, v.Odds)
	require.Equal(t, "", v.Winner)
}

func TestPreviousLotto(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("PreviousLotto", 8088, "test@admin.com")
	defer tearDown()
	_, _, _, _, students, err := CreateTestAccounts(db, 2, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	settings := openapi.Settings{
		Lottery: true,
		Odds:    250,
	}

	err = setSettings(db, &clock, userDetails, settings)
	require.Nil(t, err)

	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(100000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8088/api/lottery/previous", nil)
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v openapi.ResponseLottoLatest
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
	require.Equal(t, "No Previous Raffle", v.Winner)

	winner, err := purchaseLotto(db, &clock, userDetails, 20000)
	require.Nil(t, err)
	require.True(t, winner)

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	decoder = json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
	require.Equal(t, userDetails.Email, v.Winner)
}

func TestSearchAuctionsStudent(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServerClock("searchAuctionsStudent", 8088, "test@admin.com", &clock)
	defer tearDown()
	_, schools, teachers, _, students, err := CreateTestAccounts(db, 1, 1, 2, 1)
	require.Nil(t, err)

	teacherClasses := getTeacherClasses(db, schools[0], teachers[0])
	auctionClasses := make([]string, 0)
	auctionClasses = append(auctionClasses, teacherClasses[0].Id)

	body := openapi.RequestMakeAuction{
		Bid:         4,
		MaxBid:      4,
		Description: "Test Auction",
		EndDate:     clock.Now().Add(time.Minute * 10),
		StartDate:   clock.Now().Add(time.Minute * -10),
		OwnerId:     teachers[0],
		Visibility:  auctionClasses,
	}

	_ = MakeAuctionImpl(db, UserInfo{
		Name:     teachers[0],
		SchoolId: schools[0],
		Role:     UserRoleTeacher,
	}, body, true)

	auctionClasses = make([]string, 0)
	auctionClasses = append(auctionClasses, teacherClasses[1].Id)

	body.Visibility = auctionClasses

	_ = MakeAuctionImpl(db, UserInfo{
		Name:     teachers[0],
		SchoolId: schools[0],
		Role:     UserRoleTeacher,
	}, body, true)

	body2 := openapi.RequestAddClass{
		AddCode: teacherClasses[0].AddCode,
		Id:      "dd",
	}

	body3 := openapi.RequestAddClass{
		AddCode: teacherClasses[1].AddCode,
		Id:      "dd",
	}

	body4 := openapi.RequestKickClass{
		KickId: students[0],
		Id:     teacherClasses[0].Id,
	}

	// initialize http client
	client := &http.Client{}
	SetTestLoginUser(students[0])
	marshal, _ := json.Marshal(body2)

	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8088/api/classes/addClass", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	marshal, _ = json.Marshal(body3)

	req, _ = http.NewRequest(http.MethodPut, "http://127.0.0.1:8088/api/classes/addClass", bytes.NewBuffer(marshal))
	resp, err = client.Do(req)
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	SetTestLoginUser(teachers[0])
	marshal, _ = json.Marshal(body4)

	req, _ = http.NewRequest(http.MethodPut, "http://127.0.0.1:8088/api/classes/class/kick", bytes.NewBuffer(marshal))
	resp, err = client.Do(req)
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	SetTestLoginUser(students[0])

	req, _ = http.NewRequest(http.MethodGet, "http://127.0.0.1:8088/api/auctions/student", nil)
	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v []openapi.ResponseAuctionStudent
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
	require.Equal(t, 1, len(v))
}

func TestSearchBuckTransaction(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("searchBuckTransaction", 8088, "test@admin.com")
	defer tearDown()
	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 4, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	for _, teach := range teachers {
		err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(10000), teach, "pre load")
		require.Nil(t, err)
		err = chargeStudent(db, &clock, userDetails, decimal.NewFromFloat(5000), teach, "charge", false)
		require.Nil(t, err)
	}

	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8088/api/transactions/buckTransactions", nil)
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v []openapi.ResponseBuckTransaction
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
	require.Equal(t, float32(-5000), v[0].Amount)

}

func TestSearchBuckTransactionNegative(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("searchBuckTransactionNegative", 8088, "test@admin.com")
	defer tearDown()
	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 4, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	for _, teach := range teachers {
		err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(10000), teach, "pre load")
		require.Nil(t, err)
		err = chargeStudent(db, &clock, userDetails, decimal.NewFromFloat(50000), teach, "charge", false)
		require.Nil(t, err)
	}

	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8088/api/transactions/buckTransactions", nil)
	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v []openapi.ResponseBuckTransaction
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
	require.Equal(t, float32(50000), v[0].Amount)

}

func TestBuckConvert(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("buckConvert", 8088, "test@admin.com")
	defer tearDown()
	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 2, 3, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(10000), teachers[0], "pre load")
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1), teachers[1], "pre load")
	require.Nil(t, err)

	body := openapi.RequestBuckConvert{
		AccountFrom: teachers[0],
		AccountTo:   teachers[1],
		Amount:      1000,
	}
	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8088/api/transactions/conversionTransaction",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	body = openapi.RequestBuckConvert{
		AccountFrom: teachers[0],
		AccountTo:   teachers[1],
		Amount:      100000,
	}
	marshal, _ = json.Marshal(body)

	req, _ = http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8088/api/transactions/conversionTransaction",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 400, resp.StatusCode, resp)

}

func TestBuckConvertNewCurrency(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("buckConvertNewCurrency", 8088, "test@admin.com")
	defer tearDown()
	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 2, 3, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	otherStudent, err := getUserInLocalStore(db, students[1])
	require.Nil(t, err)

	err = pay2Student(db, &clock, otherStudent, decimal.NewFromFloat(10000), teachers[0], "pre load")
	require.Nil(t, err)
	err = addUbuck2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), "daily payment")
	require.Nil(t, err)

	body := openapi.RequestBuckConvert{
		AccountFrom: CurrencyUBuck,
		AccountTo:   teachers[0],
		Amount:      900,
	}
	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8088/api/transactions/conversionTransaction",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	_ = db.View(func(tx *bolt.Tx) error {
		student, err := getStudentBucketRx(tx, students[0])
		require.Nil(t, err)
		accounts := student.Bucket([]byte(KeyAccounts))
		account := accounts.Bucket([]byte(teachers[0]))
		balanceData := account.Get([]byte(KeyBalance))
		var balance float32
		_ = json.Unmarshal(balanceData, &balance)
		require.Equal(t, float32(900*.99), balance)
		return nil
	})

}

func TestBuckConvert_ubuck(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("buckConvert_ubuck", 8088, "test@admin.com")
	defer tearDown()
	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 4, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(10000), teachers[0], "pre load")
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1), teachers[1], "pre load")
	require.Nil(t, err)
	err = addUbuck2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), "daily payment")
	require.Nil(t, err)

	body := openapi.RequestBuckConvert{
		AccountFrom: teachers[0],
		AccountTo:   CurrencyUBuck,
		Amount:      1000,
	}
	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8088/api/transactions/conversionTransaction",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	body = openapi.RequestBuckConvert{
		AccountFrom: CurrencyUBuck,
		AccountTo:   teachers[1],
		Amount:      100,
	}
	marshal, _ = json.Marshal(body)

	req, _ = http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8088/api/transactions/conversionTransaction",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

}

func TestBuckConvert_debt(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("buckConvert_debt", 8088, "test@admin.com")
	defer tearDown()
	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 4, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(100), teachers[0], "pre load")
	require.Nil(t, err)
	err = chargeStudent(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[0], "charge", false)
	require.Nil(t, err)
	err = addUbuck2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), "daily payment")
	require.Nil(t, err)

	body := openapi.RequestBuckConvert{
		AccountFrom: teachers[0],
		AccountTo:   KeyDebt,
		Amount:      1000,
	}
	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8088/api/transactions/conversionTransaction",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 400, resp.StatusCode, resp)

	body = openapi.RequestBuckConvert{
		AccountFrom: teachers[0],
		AccountTo:   KeyDebt,
		Amount:      50,
	}
	marshal, _ = json.Marshal(body)

	req, _ = http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8088/api/transactions/conversionTransaction",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

}

func TestBuckConvert_debt_ubuck(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("buckConvert_debt_ubuck", 8088, "test@admin.com")
	defer tearDown()
	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 4, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(100), teachers[0], "pre load")
	require.Nil(t, err)
	err = chargeStudent(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[0], "charge", false)
	require.Nil(t, err)
	err = addUbuck2Student(db, &clock, userDetails, decimal.NewFromFloat(10000), "daily payment")
	require.Nil(t, err)

	body := openapi.RequestBuckConvert{
		AccountFrom: CurrencyUBuck,
		AccountTo:   KeyDebt,
		Amount:      500,
	}
	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8088/api/transactions/conversionTransaction",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	body = openapi.RequestBuckConvert{
		AccountFrom: CurrencyUBuck,
		AccountTo:   KeyDebt,
		Amount:      500,
	}
	marshal, _ = json.Marshal(body)

	req, _ = http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8088/api/transactions/conversionTransaction",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

}

func TestAuctionBid(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("auctionBid", 8088, "test@admin.com")
	defer tearDown()
	_, _, teachers, classes, students, err := CreateTestAccounts(db, 1, 1, 1, 2)
	require.Nil(t, err)

	SetTestLoginUser(teachers[0])

	// initialize http client
	client := &http.Client{}

	student0, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	student1, err := getUserInLocalStore(db, students[1])
	require.Nil(t, err)
	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	err = addUbuck2Student(db, &clock, student0, decimal.NewFromFloat(100), "starter")
	require.Nil(t, err)
	err = addUbuck2Student(db, &clock, student1, decimal.NewFromFloat(100), "starter")
	require.Nil(t, err)

	err = MakeAuctionImpl(db, teacher, openapi.RequestMakeAuction{
		Bid:         0,
		MaxBid:      0,
		Description: "test auc",
		EndDate:     clock.Now().Add(time.Minute),
		StartDate:   clock.Now(),
		OwnerId:     teacher.Name,
		Visibility:  classes,
	}, true)
	require.Nil(t, err)

	auctions, err := getTeacherAuctions(db, teacher)
	require.Nil(t, err)

	//overdrawn student 0 max 0 bid 0
	timeId := auctions[0].Id.Format(time.RFC3339Nano)

	body := openapi.RequestAuctionBid{
		Item: timeId,
		Bid:  500,
	}
	marshal, _ := json.Marshal(body)

	SetTestLoginUser(students[0])

	req, _ := http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8088/api/auctions/placeBid",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 400, resp.StatusCode, resp)

	//first bid student 0 max 50 bid 1

	body = openapi.RequestAuctionBid{
		Item: timeId,
		Bid:  50,
	}
	marshal, _ = json.Marshal(body)

	req, _ = http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8088/api/auctions/placeBid",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	//self outbid student0 max 50 bid 1
	body = openapi.RequestAuctionBid{
		Item: timeId,
		Bid:  51,
	}
	marshal, _ = json.Marshal(body)

	req, _ = http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8088/api/auctions/placeBid",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 400, resp.StatusCode, resp)

	//true outbid student1 max 91 bid 51
	body = openapi.RequestAuctionBid{
		Item: timeId,
		Bid:  91,
	}
	marshal, _ = json.Marshal(body)

	SetTestLoginUser(students[1])

	req, _ = http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8088/api/auctions/placeBid",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	//good bid but under max student0 max 91 bid 62
	body = openapi.RequestAuctionBid{
		Item: timeId,
		Bid:  61,
	}
	marshal, _ = json.Marshal(body)
	SetTestLoginUser(students[0])

	req, _ = http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8088/api/auctions/placeBid",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	//good bid but under max student0 max 91 bid 90
	body = openapi.RequestAuctionBid{
		Item: timeId,
		Bid:  89,
	}
	marshal, _ = json.Marshal(body)

	req, _ = http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8088/api/auctions/placeBid",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	//true outbid student0 max 97 bid 92
	body = openapi.RequestAuctionBid{
		Item: timeId,
		Bid:  97,
	}
	marshal, _ = json.Marshal(body)

	SetTestLoginUser(students[0])

	req, _ = http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8088/api/auctions/placeBid",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	//self outbid student0 max 97 bid 92
	body = openapi.RequestAuctionBid{
		Item: timeId,
		Bid:  100,
	}
	marshal, _ = json.Marshal(body)

	SetTestLoginUser(students[0])

	req, _ = http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8088/api/auctions/placeBid",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 400, resp.StatusCode, resp)
}

func TestSearchStudentCrypto(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("searchStudentCrypto", 8088, "test@admin.com")
	defer tearDown()
	coinGecko(db)
	_, _, _, _, students, err := CreateTestAccounts(db, 2, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	body := openapi.RequestCryptoConvert{
		Name: "cardano",
		Buy:  10,
		Sell: 0,
	}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	_ = cryptoTransaction(db, &clock, userDetails, body)

	req, _ := http.NewRequest(http.MethodGet, "http://127.0.0.1:8088/api/accounts/allCrypto", nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v []CryptoDecimal
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
	require.Equal(t, float64(10), v[0].Quantity.InexactFloat64())
	require.NotEqual(t, float64(0), v[0].Basis.InexactFloat64())
	require.NotEqual(t, float64(0), v[0].CurrentPrice.InexactFloat64())
}

func TestSearchCrypto(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("searchCrypto", 8088, "test@admin.com")
	defer tearDown()
	coinGecko(db)
	_, _, _, _, students, err := CreateTestAccounts(db, 2, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	// initialize http client
	client := &http.Client{}

	u, _ := url.ParseRequestURI("http://127.0.0.1:8088/api/accounts/crypto")
	q := u.Query()
	q.Set("name", "cardano")
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet,
		u.String(),
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v openapi.ResponseCrypto
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
	require.Equal(t, float32(0), v.Owned)
	require.Equal(t, "cardano", v.Searched)
	require.Greater(t, v.Usd, float32(0))
}

// if you run tests too many times this might fail due to coin gecko blocking you, wait a few minutes
func TestCryptoConvert(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("cryptoConvert", 8088, "test@admin.com")
	defer tearDown()
	coinGecko(db)
	_, _, _, _, students, err := CreateTestAccounts(db, 2, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	// initialize http client
	client := &http.Client{}

	body := openapi.RequestCryptoConvert{
		Name: "CarDano",
		Buy:  3,
		Sell: 0,
	}
	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8088/api/transactions/cryptoTransaction",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)
}

func TestSearchCryptoTransactions(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("searchCryptoTransactions", 8088, "test@admin.com")
	defer tearDown()
	coinGecko(db)
	_, _, _, _, students, err := CreateTestAccounts(db, 2, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(10000), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	// initialize http client
	client := &http.Client{}

	body := openapi.RequestCryptoConvert{
		Name: "CarDano",
		Buy:  3,
		Sell: 0,
	}

	student, _ := getUserInLocalStore(db, students[0])
	_ = cryptoTransaction(db, &clock, student, body)

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8088/api/transactions/cryptoTransactions", nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v []openapi.ResponseCryptoTransaction
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
	require.Equal(t, float32(3), v[0].Amount)
	require.Equal(t, float32(3), v[0].Balance)
}

func TestBuyMarketItem(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("BuyMarketItem", 8088, "test@admin.com")
	defer tearDown()
	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(10000), teachers[0], "pre load")
	require.Nil(t, err)

	request := openapi.RequestMakeMarketItem{
		Title: "Candy",
		Count: 3,
		Cost:  5,
	}

	id, err := makeMarketItem(db, &clock, teacher, request)
	require.Nil(t, err)

	// initialize http client
	client := &http.Client{}

	body := openapi.RequestMarketRefund{
		Id:        id,
		UserId:    students[0],
		TeacherId: teachers[0],
	}

	marshal, err := json.Marshal(body)
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8088/api/marketItems/buy", bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

}

func TestSearchBuck(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("SearchBuck", 8088, "test@admin.com")
	defer tearDown()
	_, _, teachers, _, students, err := CreateTestAccounts(db, 2, 2, 2, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(10000), teachers[0], "pre load")
	require.Nil(t, err)

	u, _ := url.ParseRequestURI("http://127.0.0.1:8088/api/bucks/buck")
	q := u.Query()
	q.Set("_id", teachers[0])
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	require.Nil(t, err)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var v openapi.ResponseSearchStudentUbuck
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&v)
	require.Nil(t, err)
	require.Equal(t, float32(10000), v.Value)
}

func TestLottoPurchase(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("LottoPurchase", 8088, "test@admin.com")
	defer tearDown()
	_, _, _, _, students, err := CreateTestAccounts(db, 1, 2, 3, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	var settings = openapi.Settings{
		Student2student: true,
		CurrencyLock:    false,
		Lottery:         true,
		Odds:            250,
	}

	err = setSettings(db, &clock, userDetails, settings)
	require.Nil(t, err)

	err = addUbuck2Student(db, &clock, userDetails, decimal.NewFromFloat(100000), "pre load")
	require.Nil(t, err)

	u, _ := url.ParseRequestURI("http://127.0.0.1:8088/api/lottery/purchase")
	q := u.Query()
	q.Set("quantity", "2000")
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodPut,
		u.String(),
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	lottery, err := getLottoPrevious(db, userDetails)
	require.Nil(t, err)

	assert.Equal(t, userDetails.Email, lottery.Winner)

}

func TestBuyCDEndpoint(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("BuyCDEndpoint", 8088, "test@admin.com")
	defer tearDown()
	_, _, _, _, students, err := CreateTestAccounts(db, 1, 2, 3, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(199), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	body := openapi.RequestBuyCd{
		PrinInv: 100,
		Time:    14,
	}

	marshal, err := json.Marshal(body)
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8088/api/transactions/CDTransaction", bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	req, _ = http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8088/api/transactions/CDTransaction", bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 404, resp.StatusCode, resp)
}

func TestSearchCDEndpoint(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("SearchCDEndpoint", 8088, "test@admin.com")
	defer tearDown()
	_, _, _, _, students, err := CreateTestAccounts(db, 1, 2, 3, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(500), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	body := openapi.RequestBuyCd{
		PrinInv: 100,
		Time:    14,
	}

	for i := 0; i < 5; i++ {
		err = buyCD(db, &clock, userDetails, body)
		require.Nil(t, err)
	}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8088/api/accounts/CDS", nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var data []openapi.ResponseCd
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)
	require.Equal(t, 5, len(data))
}

func TestSearchCDTransaction(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("SearchCDTransaction", 8088, "test@admin.com")
	defer tearDown()
	_, _, _, _, students, err := CreateTestAccounts(db, 1, 2, 3, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(500), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	body := openapi.RequestBuyCd{
		PrinInv: 100,
		Time:    14,
	}

	for i := 0; i < 5; i++ {
		err = buyCD(db, &clock, userDetails, body)
		require.Nil(t, err)
	}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8088/api/transactions/CDTransactions", nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	var data []openapi.ResponseTransactions
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)
	require.Equal(t, 5, len(data))
}

func TestRefundCD(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("refundCD", 8088, "test@admin.com")
	defer tearDown()
	_, _, _, _, students, err := CreateTestAccounts(db, 1, 2, 3, 2)
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	// initialize http client
	client := &http.Client{}

	userDetails, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(500), CurrencyUBuck, "pre load")
	require.Nil(t, err)

	body := openapi.RequestBuyCd{
		PrinInv: 100,
		Time:    14,
	}

	for i := 0; i < 5; i++ {
		err = buyCD(db, &clock, userDetails, body)
		require.Nil(t, err)
	}

	items, err := getCDS(db, userDetails)
	require.Nil(t, err)

	payload := openapi.RequestUser{
		Id: items[0].Ts.Format(time.RFC3339Nano),
	}

	marshal, err := json.Marshal(payload)
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8088/api/transactions/CDRefund", bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

}

type mockSSEService struct {
	broadcastCalled bool
	auctionID       string
	eventType       string
	data            interface{}
}

func (m *mockSSEService) BroadcastAuctionEvent(auctionID string, eventType string, data interface{}) {
	m.broadcastCalled = true
	m.auctionID = auctionID
	m.eventType = eventType
	m.data = data
}

func (m *mockSSEService) HandleAuctionEventsSSE(w http.ResponseWriter, r *http.Request) {
	// No-op for testing
}

func (m *mockSSEService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// No-op for testing
}

func TestBroadcastAuctionEvent(t *testing.T) {
	clock := TestClock{}
	db, teardown := FullStartTestServer("broadcastTest", 8088, "")
	defer teardown()

	mockSSE := &mockSSEService{}

	svc := &StudentApiServiceImpl{
		db:         db,
		clock:      &clock,
		sseService: mockSSE,
	}

	// Create test auction
	auction := openapi.UnifiedAuction{
		Bid:      100,
		MaxBid:   200,
		WinnerId: openapi.UnifiedAuctionOwnerId{Id: "testUser"},
	}

	// Call the broadcast directly
	svc.broadcastAuctionEvent("test-auction-id", auction)

	// Verify the broadcast was called with correct parameters
	require.True(t, mockSSE.broadcastCalled)
	require.Equal(t, "test-auction-id", mockSSE.auctionID)
	require.Equal(t, "update", mockSSE.eventType)

	// Verify the auction data was passed correctly
	auctionData, ok := mockSSE.data.(openapi.UnifiedAuction)
	require.True(t, ok)
	require.Equal(t, float32(100), auctionData.Bid)
	require.Equal(t, int32(200), auctionData.MaxBid)
	require.Equal(t, "testUser", auctionData.WinnerId.Id)
}

// Helper function to create a bid with mock SSE
func createBid(t *testing.T, db *bolt.DB, clock Clock, userDetails UserInfo, item string, bid int32) (message string, err error) {
	mockSSE := &mockSSEService{}
	return placeBid(db, clock, userDetails, item, bid, mockSSE)
}

// If you want to verify broadcasts in tests, you can use this helper
func createBidAndVerifyBroadcast(db *bolt.DB, clock Clock, userDetails UserInfo, item string, bid int32) (message string, mockSSE *mockSSEService, err error) {
	mockSSE = &mockSSEService{}
	message, err = placeBid(db, clock, userDetails, item, bid, mockSSE)
	return
}

// Example test that verifies broadcast
func TestBidBroadcast(t *testing.T) {
	clock := TestClock{}
	db, teardown := FullStartTestServer("bidBroadcast", 8088, "")
	defer teardown()

	_, _, teachers, classes, students, err := CreateTestAccounts(db, 1, 1, 1, 2)
	require.Nil(t, err)

	SetTestLoginUser(teachers[0])

	student0, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)
	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	err = addUbuck2Student(db, &clock, student0, decimal.NewFromFloat(100), "starter")
	require.Nil(t, err)

	err = MakeAuctionImpl(db, teacher, openapi.RequestMakeAuction{
		Bid:         0,
		MaxBid:      0,
		Description: "test auction",
		EndDate:     clock.Now().Add(time.Minute),
		StartDate:   clock.Now(),
		OwnerId:     teacher.Name,
		Visibility:  classes,
		TrueAuction: false,
	}, true)
	require.Nil(t, err)

	auctions, err := getTeacherAuctions(db, teacher)
	require.Nil(t, err)

	//overdrawn student 0 max 0 bid 0
	timeId := auctions[0].Id.Format(time.RFC3339Nano)

	_, mockSSE, err := createBidAndVerifyBroadcast(db, &clock, student0, timeId, 5)
	require.Nil(t, err)
	require.True(t, mockSSE.broadcastCalled)
	require.Equal(t, auctions[0].Id.UTC().Format("2006-01-02 15:04:05.000 -0700 MST"), mockSSE.auctionID)
	require.Equal(t, "update", mockSSE.eventType)
}
