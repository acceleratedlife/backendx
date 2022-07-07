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
)

func TestMakeClass(t *testing.T) {

	db, teardown := FullStartTestServer("makeClass", 8090, "")
	defer teardown()

	_, _, teachers, _, _, err := CreateTestAccounts(db, 2, 2, 2, 2)

	SetTestLoginUser(teachers[0])
	// openapi.NewAllApiService().SearchClass(classes[0])

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

	assert.Equal(t, 3, len(respData))
	assert.Equal(t, 6, len(respData[0].AddCode))

}

func TestSearchClasses(t *testing.T) {
	db, tearDown := FullStartTestServer("searchClasses", 8090, "")
	defer tearDown()
	classCount := 2

	_, _, teachers, _, _, err := CreateTestAccounts(db, 2, 2, classCount, 2)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/classes",
		nil)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data []openapi.Class
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, classCount, len(data))
	require.Equal(t, teachers[0], data[0].OwnerId)
}

func TestEditClass(t *testing.T) {
	db, tearDown := FullStartTestServer("editClass", 8090, "")
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

func TestKickClass(t *testing.T) {
	db, tearDown := FullStartTestServer("kickClass", 8090, "")
	defer tearDown()
	members := 2

	_, _, teachers, classes, students, err := CreateTestAccounts(db, 1, 1, 1, members)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}
	body := openapi.RequestKickClass{
		KickId: students[0],
		Id:     classes[0],
	}

	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8090/api/classes/class/kick",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDeleteClass(t *testing.T) {
	db, tearDown := FullStartTestServer("DeleteClass", 8090, "")
	defer tearDown()
	noClasses := 2

	_, _, teachers, classes, _, err := CreateTestAccounts(db, 2, 2, noClasses, 2)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodDelete,
		"http://127.0.0.1:8090/api/classes/class?_id="+classes[0],
		nil)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDeleteAuction(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServerClock("deleteAuction", 8090, "", &clock)
	defer tearDown()

	_, _, teachers, classes, students, err := CreateTestAccounts(db, 2, 2, 2, 2)

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

	teacherDetails, err := getUserInLocalStore(db, teachers[0])

	//to be deleted
	auctions, err := MakeAuctionImpl(db, teacherDetails, body)

	require.Nil(t, err)

	client := &http.Client{}

	u, err := url.ParseRequestURI("http://127.0.0.1:8090/api/auctions/auction")
	q := u.Query()
	q.Set("_id", auctions[0].Id.String())
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodDelete,
		u.String(),
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var respData []openapi.Auction
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&respData)

	assert.Equal(t, 0, len(respData))

	//to be deleted
	auctions, err = MakeAuctionImpl(db, teacherDetails, body)

	userDetails, _ := getUserInLocalStore(db, students[0])
	addUbuck2Student(db, &clock, userDetails, decimal.NewFromInt32(100), "loading")

	_, err = placeBid(db, &clock, userDetails, auctions[0].Id.String(), 20)
	require.Nil(t, err)

	q.Set("_id", auctions[0].Id.String())
	u.RawQuery = q.Encode()

	req, err = http.NewRequest(http.MethodDelete,
		u.String(),
		nil)

	resp, err = client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	decoder = json.NewDecoder(resp.Body)
	_ = decoder.Decode(&respData)

	assert.Equal(t, 0, len(respData))

	//to be de-activated
	auctions, err = MakeAuctionImpl(db, teacherDetails, body)
	require.Nil(t, err)

	_, err = placeBid(db, &clock, userDetails, auctions[0].Id.String(), 20)
	require.Nil(t, err)
	clock.TickOne(time.Hour * 1)

	q.Set("_id", auctions[0].Id.String())
	u.RawQuery = q.Encode()

	req, err = http.NewRequest(http.MethodDelete,
		u.String(),
		nil)

	resp, err = client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	decoder = json.NewDecoder(resp.Body)
	_ = decoder.Decode(&respData)

	assert.Equal(t, 1, len(respData))

}

func TestMakeAuction(t *testing.T) {
	clock := TestClock{}
	db, teardown := FullStartTestServer("makeClass", 8090, "")
	defer teardown()

	_, _, teachers, classes, _, err := CreateTestAccounts(db, 1, 1, 2, 2)

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

	defer resp.Body.Close()
	var respData []openapi.Auction
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&respData)

	assert.Equal(t, len(classes), len(respData[0].Visibility))
	assert.Equal(t, 1, len(respData))

}

func TestSearchAuctionsTeacher(t *testing.T) {
	clock := TestClock{}
	db, teardown := FullStartTestServer("searchAuctionsTeacher", 8090, "")
	defer teardown()

	_, schools, teachers, classes, _, err := CreateTestAccounts(db, 1, 1, 2, 2)

	SetTestLoginUser(teachers[0])

	body := openapi.RequestMakeAuction{
		Bid:         4,
		MaxBid:      4,
		Description: "Test Auction",
		EndDate:     clock.Now().Add(500),
		StartDate:   clock.Now(),
		OwnerId:     teachers[0],
		Visibility:  classes,
	}

	_, err = MakeAuctionImpl(db, UserInfo{
		Name:     teachers[0],
		SchoolId: schools[0],
		Role:     UserRoleTeacher,
	}, body)

	require.Nil(t, err)

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/auctions",
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	defer resp.Body.Close()
	var respData []openapi.Auction
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&respData)

	assert.Equal(t, 1, len(respData))
	assert.Equal(t, len(classes), len(respData[0].Visibility))

}

func TestSearchTransactions(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("searchTransactions", 8090, "")
	defer tearDown()
	numStudents := 15

	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 1, 1, numStudents)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}

	for _, student := range students {
		userDetails, err := getUserInLocalStore(db, student)
		require.Nil(t, err)
		err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[0], "pre load")
		require.Nil(t, err)
		err = chargeStudent(db, &clock, userDetails, decimal.NewFromFloat(100), teachers[0], "charge", false)
		require.Nil(t, err)
	}

	marshal, _ := json.Marshal(teachers[0])

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/transactions?_id="+teachers[0],
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var respData []openapi.ResponseTransactions
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&respData)

	assert.Equal(t, numStudents*2, len(respData))

}

func TestDeleteStudent(t *testing.T) {
	db, tearDown := FullStartTestServer("deleteStudent", 8090, "")
	defer tearDown()
	numStudents := 50

	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 1, 1, numStudents)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodDelete,
		"http://127.0.0.1:8090/api/users/user?_id="+students[0],
		bytes.NewBuffer(nil))

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

}

func TestResetPassword(t *testing.T) {
	db, tearDown := FullStartTestServer("resetPassword", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, err := CreateTestAccounts(db, 2, 2, 2, 2)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}
	body := openapi.RequestUser{
		Id: students[0],
	}

	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/users/resetPassword",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data openapi.ResponseResetPassword
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.Equal(t, 6, len(data.Password))
}

func TestSearchEvents(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("searchEvents", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 1, 1, 2)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}

	for _, student := range students {
		userDetails, err := getUserInLocalStore(db, student)
		require.Nil(t, err)
		err = pay2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[0], "pre load")
		require.Nil(t, err)
		EventIfNeeded(db, &clock, userDetails)
	}

	clock.TickOne(time.Hour * 240)

	for _, student := range students {
		userDetails, err := getUserInLocalStore(db, student)
		require.Nil(t, err)
		EventIfNeeded(db, &clock, userDetails)
	}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/events", nil)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var respData []openapi.ResponseEvents
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&respData)

	assert.NotZero(t, respData[0].Value)

}
