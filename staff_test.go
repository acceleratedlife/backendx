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
	body := openapi.RequestUser{
		Id: teachers[0],
	}

	marshal, _ := json.Marshal(body)

	req, _ := http.NewRequest(http.MethodGet,
		"http://127.0.0.1:8090/api/classes",
		bytes.NewBuffer(marshal))

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

func TestMakeAuction(t *testing.T) {
	clock := AppClock{}
	db, teardown := FullStartTestServer("makeClass", 8090, "")
	defer teardown()

	_, _, teachers, classes, _, err := CreateTestAccounts(db, 2, 2, 2, 2)

	SetTestLoginUser(teachers[0])

	client := &http.Client{}

	body := openapi.RequestMakeAuction{
		Bid:         4,
		MaxBid:      4,
		Description: "Test Auction",
		EndDate:     clock.Now().Add(500),
		StartDate:   clock.Now(),
		Owner_id:    teachers[0],
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

	assert.Equal(t, len(classes)-4, len(respData[0].Visibility)) // -4 because both schools have freshman sophomores... the key is the same so only added once
	assert.Equal(t, 1, len(respData))
	// assert.Equal(t, 6, len(respData[0].AddCode))

}

// func TestPayTransaction_credit(t *testing.T) {
// 	db, tearDown := FullStartTestServer("payTransaction_credit", 8090, "")
// 	defer tearDown()

// 	_, _, teachers, _, students, err := CreateTestAccounts(db, 2, 2, 2, 2)

// 	SetTestLoginUser(teachers[0])

// 	client := &http.Client{}
// 	body := openapi.RequestPayTransaction{
// 		Owner:       "",
// 		Description: "credit",
// 		Amount:      100,
// 		Student:     students[0],
// 	}

// 	marshal, _ := json.Marshal(body)

// 	req, _ := http.NewRequest(http.MethodPost,
// 		"http://127.0.0.1:8090/api/transactions/payTransaction",
// 		bytes.NewBuffer(marshal))

// 	resp, err := client.Do(req)
// 	defer resp.Body.Close()
// 	require.Nil(t, err)
// 	require.NotNil(t, resp)
// 	assert.Equal(t, 200, resp.StatusCode)

// require.Equal(t, members, len(data.Members))
// require.Equal(t, "Test Name", data.Name)
// require.Equal(t, int32(4), data.Period)
// require.Equal(t, classes[0], data.Id)
// }
