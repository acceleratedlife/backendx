package main

// This test will fail in github as it does not know the correct config variables
// This test will fail locally unless you supply the correct email and password in local config file
// commented out the test as github will fail it
// func TestResetStaffPassword(t *testing.T) {
// db, tearDown := FullStartTestServer("resetStaffPassword", 8088, "")
// defer tearDown()

// _, _, teachers, _, _, _ := CreateTestAccounts(db, 1, 1, 0, 0)

// client := &http.Client{}
// body := openapi.RequestUser{
// 	Id: teachers[0],
// }

// marshal, _ := json.Marshal(body)

// req, _ := http.NewRequest(http.MethodPost,
// 	"http://127.0.0.1:8088/api/users/resetStaffPassword",
// 	bytes.NewBuffer(marshal))

// resp, err := client.Do(req)
// defer resp.Body.Close()
// require.Nil(t, err)
// require.NotNil(t, resp)
// assert.Equal(t, 200, resp.StatusCode)

// }

//works great just a new problem with coin gecko changes

// func TestGetCryptos(t *testing.T) {
// 	db, tearDown := FullStartTestServer("getCryptos", 8088, "")
// 	defer tearDown()
// 	err := coinGecko(db)
// 	require.Nil(t, err)

// 	client := &http.Client{}

// 	req, _ := http.NewRequest(http.MethodGet,
// 		"http://127.0.0.1:8088/api/allCrypto", nil)

// 	resp, err := client.Do(req)
// 	require.Nil(t, err)
// 	defer resp.Body.Close()
// 	require.Nil(t, err)
// 	require.NotNil(t, resp)
// 	assert.Equal(t, 200, resp.StatusCode)

// 	var data []openapi.ResponseCryptoPrice
// 	decoder := json.NewDecoder(resp.Body)
// 	_ = decoder.Decode(&data)

// 	require.Greater(t, len(data), 20)
// }
