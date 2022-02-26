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

func TestAllSchoolApiServiceImpl_AddCodeClass(t *testing.T) {
	db, tearDown := FullStartTestServer("addCode", 8090, "test@admin.com")
	defer tearDown()

	_, _, _, classes, _, err := CreateTestAccounts(db, 2, 2, 2, 2)

	// initialize http client
	client := &http.Client{}

	body := openapi.RequestEditClass{
		Id: classes[0],
	}

	marshal, _ := json.Marshal(body)
	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8090/api/classes/class/addCode", bytes.NewBuffer(marshal))
	resp, err := client.Do(req)
	require.Nil(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode, resp)

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	var v openapi.ClassWithMembers
	_ = decoder.Decode(&v)
	assert.Equal(t, 6, len(v.AddCode))
}
