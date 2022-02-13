package main

import (
	"bytes"
	"encoding/json"
	openapi "github.com/acceleratedlife/backend/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestAllSchoolApiServiceImpl_AddCodeClass(t *testing.T) {
	_, tearDown := FullStartTestServer("addCode", 8090, "test@admin.com")
	defer tearDown()

	// initialize http client
	client := &http.Client{}

	req, _ := http.NewRequest(http.MethodPut, "http://127.0.0.1:8090/api/classes/class/addCode", bytes.NewBuffer([]byte("{}")))
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
