package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// If you just can spec then ResponseSchool model should not have omitEmpty on students
func TestSearchSchools(t *testing.T) {

	db, teardown := FullStartTestServer("searchSchools", 8090, "")
	defer teardown()

	CreateTestAccounts(db, 3, 2, 2, 2)

	sysAdmin := UserInfo{
		FirstName: "top",
		LastName:  "boss",
		Role:      UserRoleSysAdmin,
		Email:     "boss@boss.com",
	}

	_, err := CreateSysAdmin(db, sysAdmin)
	require.Nil(t, err)

	SA, err := getUserInLocalStore(db, sysAdmin.Email)
	require.Nil(t, err)

	SetTestLoginUser(SA.Email)

	assert.Nil(t, err)
	client := &http.Client{}

	u, err := url.ParseRequestURI("http://127.0.0.1:8090/api/schools")
	require.Nil(t, err)

	q := u.Query()
	q.Set("zip", "")
	u.RawQuery = q.Encode()

	req, _ := http.NewRequest(http.MethodGet, u.String(), nil)

	resp, err := client.Do(req)
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)

	defer resp.Body.Close()
	var respData []openapi.ResponseSchools
	decoder := json.NewDecoder(resp.Body)
	_ = decoder.Decode(&respData)

	assert.Equal(t, 4, len(respData)) // one greater because there is some setup school that has nothing from InitDefaultAccounts
	assert.Equal(t, int32(8), respData[0].Students)

}
