package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/sensu/sensu-go/testing/mockstore"
	"github.com/sensu/sensu-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHttpApiChecksGet(t *testing.T) {
	store := &mockstore.MockStore{}

	c := &ChecksController{
		Store: store,
	}

	checks := []*types.CheckConfig{
		types.FixtureCheckConfig("check1"),
		types.FixtureCheckConfig("check2"),
	}
	store.On("GetCheckConfigs", mock.Anything).Return(checks, nil)
	req, _ := http.NewRequest("GET", "/checks", nil)
	res := processRequest(c, req)

	assert.Equal(t, http.StatusOK, res.Code)

	body := res.Body.Bytes()

	returnedChecks := []*types.CheckConfig{}
	err := json.Unmarshal(body, &returnedChecks)

	assert.NoError(t, err)
	assert.EqualValues(t, checks, returnedChecks)
}

func TestHttpApiChecksGetError(t *testing.T) {
	store := &mockstore.MockStore{}

	c := &ChecksController{
		Store: store,
	}

	var nilChecks []*types.CheckConfig
	store.On("GetCheckConfigs", mock.Anything).Return(nilChecks, errors.New("error"))
	req, _ := http.NewRequest("GET", "/checks", nil)
	res := processRequest(c, req)

	body := res.Body.Bytes()

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Equal(t, "error\n", string(body))
}

func TestHttpApiCheckGet(t *testing.T) {
	store := &mockstore.MockStore{}

	c := &ChecksController{
		Store: store,
	}

	var nilCheck *types.CheckConfig
	store.On("GetCheckConfigByName", mock.Anything, "somecheck").Return(nilCheck, nil)
	notFoundReq, _ := http.NewRequest("GET", "/checks/somecheck", nil)
	notFoundRes := processRequest(c, notFoundReq)

	assert.Equal(t, http.StatusNotFound, notFoundRes.Code)

	check1 := types.FixtureCheckConfig("check1")
	store.On("GetCheckConfigByName", mock.Anything, "check1").Return(check1, nil)
	foundReq, _ := http.NewRequest("GET", "/checks/check1", nil)
	foundRes := processRequest(c, foundReq)

	assert.Equal(t, http.StatusOK, foundRes.Code)

	body := foundRes.Body.Bytes()

	check := &types.CheckConfig{}
	err := json.Unmarshal(body, &check)

	assert.NoError(t, err)
	assert.NotNil(t, check.Name)
	assert.NotNil(t, check.Command)
	assert.NotEqual(t, check.Name, "")
	assert.NotEqual(t, check.Command, "")
}

func TestHttpApiCheckPut(t *testing.T) {
	store := &mockstore.MockStore{}

	c := &ChecksController{
		Store: store,
	}

	check := types.FixtureCheckConfig("check1")
	updatedCheckJSON, _ := json.Marshal(check)

	store.On("UpdateCheckConfig", mock.Anything, mock.AnythingOfType("*types.CheckConfig")).Return(nil).Run(func(args mock.Arguments) {
		receivedCheck := args.Get(1).(*types.CheckConfig)
		assert.NoError(t, receivedCheck.Validate())
		assert.EqualValues(t, check, receivedCheck)
	})
	putReq, _ := http.NewRequest("PUT", fmt.Sprintf("/checks/%s", "check1"), bytes.NewBuffer(updatedCheckJSON))
	putRes := processRequest(c, putReq)

	assert.Equal(t, http.StatusOK, putRes.Code)
}

func TestHttpApiCheckPost(t *testing.T) {
	store := &mockstore.MockStore{}

	c := &ChecksController{
		Store: store,
	}

	check := types.FixtureCheckConfig("check1")
	updatedCheckJSON, _ := json.Marshal(check)

	store.On("UpdateCheckConfig", mock.Anything, mock.AnythingOfType("*types.CheckConfig")).Return(nil).Run(func(args mock.Arguments) {
		receivedCheck := args.Get(1).(*types.CheckConfig)
		assert.NoError(t, receivedCheck.Validate())
		assert.EqualValues(t, check, receivedCheck)
	})
	putReq, _ := http.NewRequest("POST", fmt.Sprintf("/checks/check1"), bytes.NewBuffer(updatedCheckJSON))
	putRes := processRequest(c, putReq)

	assert.Equal(t, http.StatusOK, putRes.Code)
}

func TestHttpApiCheckDelete(t *testing.T) {
	store := &mockstore.MockStore{}

	c := &ChecksController{
		Store: store,
	}

	check := types.FixtureCheckConfig("check1")
	store.On("GetCheckConfigByName", mock.Anything, "check1").Return(check, nil)
	store.On("DeleteCheckConfigByName", mock.Anything, "check1").Return(nil)
	deleteReq, _ := http.NewRequest("DELETE", fmt.Sprintf("/checks/check1"), nil)
	deleteRes := processRequest(c, deleteReq)

	assert.Equal(t, http.StatusOK, deleteRes.Code)
}
