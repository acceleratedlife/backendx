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

func TestEditClass(t *testing.T) {
	db, tearDown := FullStartTestServer("editClass", 8090, "")
	defer tearDown()
	members := 2

	_, _, teachers, classes, _, _ := CreateTestAccounts(db, 2, 2, 2, members)

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
	require.Nil(t, err)
	defer resp.Body.Close()
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

	_, _, teachers, classes, students, _ := CreateTestAccounts(db, 1, 1, 1, members)

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
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDeleteClass(t *testing.T) {
	db, tearDown := FullStartTestServer("DeleteClass", 8090, "")
	defer tearDown()
	noClasses := 2

	_, _, teachers, classes, _, _ := CreateTestAccounts(db, 2, 2, noClasses, 2)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodDelete,
		"http://127.0.0.1:8090/api/classes/class?_id="+classes[0],
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestSearchAuctionsTeacher(t *testing.T) {
	clock := TestClock{}
	db, teardown := FullStartTestServer("searchAuctionsTeacher", 8090, "")
	defer teardown()

	_, schools, teachers, classes, _, _ := CreateTestAccounts(db, 1, 1, 2, 2)

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

	err := MakeAuctionImpl(db, UserInfo{
		Name:     teachers[0],
		SchoolId: schools[0],
		Role:     UserRoleTeacher,
	}, body, true)

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

func TestSearchAuctionsTeacherStudent(t *testing.T) {
	clock := TestClock{}
	db, teardown := FullStartTestServer("searchAuctionsTeacherStudent", 8090, "")
	defer teardown()

	_, schools, _, classes, students, _ := CreateTestAccounts(db, 1, 1, 2, 2)

	SetTestLoginUser(students[0])

	body := openapi.RequestMakeAuction{
		Bid:         4,
		MaxBid:      4,
		Description: "Test Auction",
		EndDate:     clock.Now().Add(500),
		StartDate:   clock.Now(),
		OwnerId:     students[0],
		Visibility:  classes,
	}

	err := MakeAuctionImpl(db, UserInfo{
		Name:     students[0],
		SchoolId: schools[0],
		Role:     UserRoleStudent,
	}, body, true)

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

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 1, 1, numStudents)

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
	require.Nil(t, err)
	defer resp.Body.Close()
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

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 1, 1, numStudents)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodDelete,
		"http://127.0.0.1:8090/api/users/user?_id="+students[0],
		bytes.NewBuffer(nil))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

}

func TestResetPassword(t *testing.T) {
	db, tearDown := FullStartTestServer("resetPassword", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 2, 2, 2, 2)

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
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data openapi.ResponseResetPassword
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.GreaterOrEqual(t, len(data.Password), 6)
}

func TestSearchEvents(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("searchEvents", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, _ := CreateTestAccounts(db, 1, 1, 1, 2)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}

	event := EventRequest{
		Positive:    false,
		Description: "Pay Taxes",
		Title:       "Taxes",
	}

	marshal, _ := json.Marshal(event)

	err := createJobOrEvent(db, marshal, KeyNEvents, "Teacher")
	require.Nil(t, err)

	event = EventRequest{
		Positive:    true,
		Description: "Pay Taxes",
		Title:       "Taxes",
	}

	marshal, _ = json.Marshal(event)

	err = createJobOrEvent(db, marshal, KeyPEvents, "Teacher")
	require.Nil(t, err)

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
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var respData []openapi.ResponseEvents
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&respData)

	assert.Greater(t, len(respData), 0)
	assert.NotZero(t, respData[0].Value)

}

func TestAuctionsAllGet(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("auctionsAllGet", 8090, "")
	defer tearDown()

	_, _, teachers, classes, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	SetTestLoginUser(teachers[0])

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	client := &http.Client{}

	body := openapi.RequestMakeAuction{
		Bid:         4,
		MaxBid:      4,
		Description: "Test Auction",
		EndDate:     clock.Now().Add(time.Minute * 100),
		StartDate:   clock.Now().Add(time.Minute * -10),
		OwnerId:     teachers[0],
		Visibility:  classes,
	}

	err = MakeAuctionImpl(db, teacher, body, true)
	require.Nil(t, err)

	body = openapi.RequestMakeAuction{
		Bid:         4,
		MaxBid:      4,
		Description: "Test Auction",
		EndDate:     clock.Now().Add(time.Minute * 10),
		StartDate:   clock.Now().Add(time.Minute * -10),
		OwnerId:     students[0],
		Visibility:  classes,
	}

	err = MakeAuctionImpl(db, teacher, body, true)
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/auctions/all",
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var respData []openapi.ResponseAuctionStudent
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&respData)

	assert.Equal(t, 2, len(respData))

}

func TestAuctionApprove(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("auctionApprove", 8090, "")
	defer tearDown()

	_, _, teachers, classes, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	SetTestLoginUser(teachers[0])

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	client := &http.Client{}

	body := openapi.RequestMakeAuction{
		Bid:         4,
		MaxBid:      4,
		Description: "Test Auction",
		EndDate:     clock.Now().Add(time.Minute * 100),
		StartDate:   clock.Now().Add(time.Minute * -10),
		OwnerId:     students[0],
		Visibility:  classes,
	}

	err = MakeAuctionImpl(db, student, body, false)
	require.Nil(t, err)

	auctions, err := getAllAuctions(db, &clock, teacher)
	require.Nil(t, err)

	action := openapi.RequestAuctionAction{
		AuctionId: auctions[0].Id.Format(time.RFC3339Nano),
	}

	marshal, err := json.Marshal(action)
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8090/api/auctions/approve",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

}

func TestAuctionReject(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("auctionReject", 8090, "")
	defer tearDown()

	_, _, teachers, classes, students, _ := CreateTestAccounts(db, 1, 1, 1, 1)

	SetTestLoginUser(teachers[0])

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	client := &http.Client{}

	body := openapi.RequestMakeAuction{
		Bid:         4,
		MaxBid:      4,
		Description: "Test Auction",
		EndDate:     clock.Now().Add(time.Minute * 100),
		StartDate:   clock.Now().Add(time.Minute * -10),
		OwnerId:     students[0],
		Visibility:  classes,
	}

	err = MakeAuctionImpl(db, student, body, false)
	require.Nil(t, err)

	body = openapi.RequestMakeAuction{
		Bid:         4,
		MaxBid:      4,
		Description: "Test Auction",
		EndDate:     clock.Now().Add(time.Minute * 10),
		StartDate:   clock.Now().Add(time.Minute * -10),
		OwnerId:     students[0],
		Visibility:  classes,
	}

	err = MakeAuctionImpl(db, teacher, body, true)
	require.Nil(t, err)

	auctions, err := getAllAuctions(db, &clock, teacher)
	require.Nil(t, err)
	assert.Equal(t, 2, len(auctions))

	timeId := auctions[0].Id.Format(time.RFC3339Nano)

	u, _ := url.ParseRequestURI("http://127.0.0.1:8090/api/auctions/reject")
	q := u.Query()
	q.Set("_id", timeId)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodDelete,
		u.String(),
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	auctions, err = getAllAuctions(db, &clock, teacher)
	require.Nil(t, err)
	assert.Equal(t, 1, len(auctions))

}

// omit empty should not be on Student2students or Lottery
func TestStudent2StudentOmitEmpty(t *testing.T) {

	settings := openapi.Settings{
		Student2student: true,
		CurrencyLock:    false,
		Lottery:         true,
		Odds:            5000,
	}

	marshal, err := json.Marshal(settings)
	require.Nil(t, err)

	var postSettings openapi.Settings //testing to see that omit empty is in the correct places

	err = json.Unmarshal(marshal, &postSettings)
	require.Nil(t, err)

	require.True(t, postSettings.Student2student)
	require.True(t, postSettings.Lottery)
	require.Equal(t, int32(5000), postSettings.Odds)

	settings = openapi.Settings{
		Student2student: false,
		CurrencyLock:    false,
		Lottery:         false,
	}

	marshal, err = json.Marshal(settings)
	require.Nil(t, err)

	err = json.Unmarshal(marshal, &postSettings)
	require.Nil(t, err)

	require.False(t, postSettings.Student2student)
	require.False(t, postSettings.Lottery)

}

// if the following test is failing and you just ran the spec then it is probably because student2student and/or Lottery has been changed
// to omitempty on openapi.Settings. This causes nothing to be sent back because false is seen as empty and is omitted.
func TestGetSettingsAdmin(t *testing.T) {
	db, tearDown := FullStartTestServer("getSettings", 8090, "")
	defer tearDown()

	admins, _, _, _, _, _ := CreateTestAccounts(db, 1, 2, 1, 3)

	SetTestLoginUser(admins[0])

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/settings",
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data openapi.Settings
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.False(t, data.Student2student)
	require.False(t, data.Lottery)

	admin, err := getUserInLocalStore(db, admins[0])
	require.Nil(t, err)

	setSettings(db, admin, openapi.Settings{
		Student2student: true,
		Lottery:         true,
		Odds:            2010,
	})

	settings, err := getSettings(db, admin)
	require.Nil(t, err)
	require.True(t, settings.Student2student)
	require.True(t, settings.Lottery)
	require.Equal(t, int32(2010), settings.Odds)

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	decoder = json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.True(t, data.Student2student)
	require.True(t, settings.Lottery)
	require.Equal(t, int32(2010), settings.Odds)
}

func TestGetSettingsTeacher(t *testing.T) {
	db, tearDown := FullStartTestServer("getSettings", 8090, "")
	defer tearDown()

	_, _, teachers, _, _, _ := CreateTestAccounts(db, 1, 2, 1, 3)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/settings",
		nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	var data openapi.Settings
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&data)

	require.False(t, data.CurrencyLock)
}

// if the following test is failing and you just ran the spec then it is probably because you set student2student and/or lottery as required in the spec
// This causes openapi.Settings to require student2student to be true so when you make it false marshal does not know what to do with it.
// The only solution at this time is to omit the requirement in spec. This will make getsettings fail until you delete omitempty
func TestSetSettingsAdmin(t *testing.T) {
	db, tearDown := FullStartTestServer("setSettings", 8090, "")
	defer tearDown()

	admins, _, _, _, _, _ := CreateTestAccounts(db, 1, 2, 1, 3)

	SetTestLoginUser(admins[0])

	client := &http.Client{}

	admin, err := getUserInLocalStore(db, admins[0])
	require.Nil(t, err)

	settings, err := getSettings(db, admin)
	require.Nil(t, err)
	require.False(t, settings.Student2student)
	require.False(t, settings.Lottery)

	settings = openapi.Settings{
		Student2student: true,
		CurrencyLock:    false,
		Lottery:         true,
		Odds:            2012,
	}

	marshal, err := json.Marshal(settings)
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8090/api/settings",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	settings, err = getSettings(db, admin)
	require.Nil(t, err)
	require.True(t, settings.Student2student)
	require.True(t, settings.Lottery)
	require.Equal(t, int32(2012), settings.Odds)

	settings = openapi.Settings{
		Student2student: false,
		CurrencyLock:    false,
		Lottery:         false,
	}

	marshal, _ = json.Marshal(settings)

	req, _ = http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8090/api/settings",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	settings, err = getSettings(db, admin)
	require.Nil(t, err)
	require.False(t, settings.Student2student)
	require.False(t, settings.Lottery)
}

func TestSetSettingsTeacher(t *testing.T) {
	db, tearDown := FullStartTestServer("setSettings", 8090, "")
	defer tearDown()

	_, _, teachers, _, _, _ := CreateTestAccounts(db, 1, 2, 1, 3)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	require.False(t, teacher.Settings.CurrencyLock)

	settings := openapi.Settings{
		CurrencyLock: true,
	}

	marshal, err := json.Marshal(settings)
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8090/api/settings",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	teacher, err = getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)
	require.True(t, teacher.Settings.CurrencyLock)

	settings = openapi.Settings{
		CurrencyLock: false,
	}

	marshal, _ = json.Marshal(settings)

	req, _ = http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8090/api/settings",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	teacher, err = getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)
	require.False(t, teacher.Settings.CurrencyLock)
}

func TestStudent2Student(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("payTransactions_student", 8090, "")
	defer tearDown()

	admins, _, _, _, students, err := CreateTestAccounts(db, 1, 2, 2, 2)
	require.Nil(t, err)

	admin, err := getUserInLocalStore(db, admins[0])
	require.Nil(t, err)

	setSettings(db, admin, openapi.Settings{Student2student: true})

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

	settings, err := getSettings(db, admin)
	require.Nil(t, err)
	require.True(t, settings.Student2student)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/transactions/payTransaction",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
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
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 400, resp.StatusCode)

}

// ********** this is currently testing the same as student2student. Need to rewrite once the logic has been implemented.
func TestLottery(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("lottery", 8090, "")
	defer tearDown()

	admins, _, _, _, students, err := CreateTestAccounts(db, 1, 2, 2, 2)
	require.Nil(t, err)

	admin, err := getUserInLocalStore(db, admins[0])
	require.Nil(t, err)

	setSettings(db, admin, openapi.Settings{Student2student: true})

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

	settings, err := getSettings(db, admin)
	require.Nil(t, err)
	require.True(t, settings.Student2student)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/transactions/payTransaction",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
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
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 400, resp.StatusCode)

}

func TestCurrencyLock(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("payTransactions_student", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 2, 2, 2)
	require.Nil(t, err)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	teacher.Settings.CurrencyLock = true
	err = userEdit(db, &clock, teacher, openapi.RequestUserEdit{})
	require.Nil(t, err)

	SetTestLoginUser(students[0])

	client := &http.Client{}
	body := openapi.RequestBuckConvert{
		AccountFrom: CurrencyUBuck,
		AccountTo:   teachers[0],
		Amount:      100,
	}

	for _, student := range students {
		userDetails, err := getUserInLocalStore(db, student)
		require.Nil(t, err)
		err = addUbuck2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), "pre load")
		require.Nil(t, err)
		err = addBuck2Student(db, &clock, userDetails, decimal.NewFromFloat(1000), teachers[0], "pre load")
		require.Nil(t, err)
	}

	marshal, _ := json.Marshal(body)

	teacher, err = getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)
	require.True(t, teacher.Settings.CurrencyLock)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/transactions/conversionTransaction",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 400, resp.StatusCode)

	teacher.Settings.CurrencyLock = false
	err = userEdit(db, &clock, teacher, openapi.RequestUserEdit{})
	require.Nil(t, err)

	req, _ = http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/transactions/conversionTransaction",
		bytes.NewBuffer(marshal))

	resp, err = client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestMakeMarketItem(t *testing.T) {
	db, tearDown := FullStartTestServer("makeMarketItem", 8090, "")
	defer tearDown()

	_, _, teachers, _, _, _ := CreateTestAccounts(db, 2, 2, 2, 2)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}
	body := openapi.RequestMakeMarketItem{
		Title: "Candy",
		Count: 4,
		Cost:  99,
	}

	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodPost,
		"http://127.0.0.1:8090/api/marketItems",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

}

func TestMarketItemResolve(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("marketItemResolve", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.Nil(t, err)
	SetTestLoginUser(teachers[0])

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(1000), teachers[0], "pre load")
	require.Nil(t, err)

	client := &http.Client{}

	itemId, err := makeMarketItem(db, &clock, teacher, openapi.RequestMakeMarketItem{
		Title: "Candy",
		Count: 4,
		Cost:  3,
	})
	require.Nil(t, err)

	purchaseId, err := buyMarketItem(db, &clock, student, teacher, itemId)
	require.Nil(t, err)

	request := openapi.RequestMarketRefund{
		Id:     itemId,
		UserId: purchaseId,
	}

	marshal, err := json.Marshal(request)
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8090/api/marketItems/resolve",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

}

func TestMarketItemRefund(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("marketItemRefund", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.Nil(t, err)
	SetTestLoginUser(teachers[0])

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(1000), teachers[0], "pre load")
	require.Nil(t, err)

	client := &http.Client{}

	itemId, err := makeMarketItem(db, &clock, teacher, openapi.RequestMakeMarketItem{
		Title: "Candy",
		Count: 4,
		Cost:  3,
	})
	require.Nil(t, err)

	purchaseId, err := buyMarketItem(db, &clock, student, teacher, itemId)
	require.Nil(t, err)

	request := openapi.RequestMarketRefund{
		Id:     itemId,
		UserId: purchaseId,
	}

	marshal, err := json.Marshal(request)
	require.Nil(t, err)

	req, _ := http.NewRequest(http.MethodPut,
		"http://127.0.0.1:8090/api/marketItems/refund",
		bytes.NewBuffer(marshal))

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

}

func TestMarketItemDelete(t *testing.T) {
	clock := TestClock{}
	db, tearDown := FullStartTestServer("marketItemDelete", 8090, "")
	defer tearDown()

	_, _, teachers, _, students, err := CreateTestAccounts(db, 1, 1, 1, 1)
	require.Nil(t, err)
	SetTestLoginUser(teachers[0])

	student, err := getUserInLocalStore(db, students[0])
	require.Nil(t, err)

	teacher, err := getUserInLocalStore(db, teachers[0])
	require.Nil(t, err)

	err = pay2Student(db, &clock, student, decimal.NewFromFloat(1000), teachers[0], "pre load")
	require.Nil(t, err)

	client := &http.Client{}

	itemId, err := makeMarketItem(db, &clock, teacher, openapi.RequestMakeMarketItem{
		Title: "Candy",
		Count: 4,
		Cost:  3,
	})
	require.Nil(t, err)

	_, err = buyMarketItem(db, &clock, student, teacher, itemId)
	require.Nil(t, err)

	u, err := url.ParseRequestURI("http://127.0.0.1:8090/api/marketItems")
	require.Nil(t, err)

	q := u.Query()
	q.Set("_id", itemId)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodDelete, u.String(), nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	defer resp.Body.Close()
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

}

func TestGetStudentCount(t *testing.T) {
	db, tearDown := FullStartTestServer("getStudentCount", 8090, "")
	defer tearDown()

	_, schools, _, _, _, _ := CreateTestAccounts(db, 1, 2, 2, 10)

	count, err := getStudentCount(db, schools[0])
	require.Nil(t, err)

	require.Equal(t, int32(40), count)

}
