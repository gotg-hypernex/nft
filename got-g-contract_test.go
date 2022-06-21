/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const getStateError = "world state get error"

type MockStub struct {
	shim.ChaincodeStubInterface
	mock.Mock
}

func (ms *MockStub) GetState(key string) ([]byte, error) {
	args := ms.Called(key)

	return args.Get(0).([]byte), args.Error(1)
}

func (ms *MockStub) PutState(key string, value []byte) error {
	args := ms.Called(key, value)

	return args.Error(0)
}

func (ms *MockStub) DelState(key string) error {
	args := ms.Called(key)

	return args.Error(0)
}

type MockContext struct {
	contractapi.TransactionContextInterface
	mock.Mock
}

func (mc *MockContext) GetStub() shim.ChaincodeStubInterface {
	args := mc.Called()

	return args.Get(0).(*MockStub)
}

func configureStub() (*MockContext, *MockStub) {
	var nilBytes []byte

	testGotG := new(GotG)
	testGotG.Value = "set value"
	gotGBytes, _ := json.Marshal(testGotG)

	ms := new(MockStub)
	ms.On("GetState", "statebad").Return(nilBytes, errors.New(getStateError))
	ms.On("GetState", "missingkey").Return(nilBytes, nil)
	ms.On("GetState", "existingkey").Return([]byte("some value"), nil)
	ms.On("GetState", "gotGkey").Return(gotGBytes, nil)
	ms.On("PutState", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
	ms.On("DelState", mock.AnythingOfType("string")).Return(nil)

	mc := new(MockContext)
	mc.On("GetStub").Return(ms)

	return mc, ms
}

func TestGotGExists(t *testing.T) {
	var exists bool
	var err error

	ctx, _ := configureStub()
	c := new(GotGContract)

	exists, err = c.GotGExists(ctx, "statebad")
	assert.EqualError(t, err, getStateError)
	assert.False(t, exists, "should return false on error")

	exists, err = c.GotGExists(ctx, "missingkey")
	assert.Nil(t, err, "should not return error when can read from world state but no value for key")
	assert.False(t, exists, "should return false when no value for key in world state")

	exists, err = c.GotGExists(ctx, "existingkey")
	assert.Nil(t, err, "should not return error when can read from world state and value exists for key")
	assert.True(t, exists, "should return true when value for key in world state")
}

func TestCreateGotG(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(GotGContract)

	err = c.CreateGotG(ctx, "statebad", "some value")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.CreateGotG(ctx, "existingkey", "some value")
	assert.EqualError(t, err, "The asset existingkey already exists", "should error when exists returns true")

	err = c.CreateGotG(ctx, "missingkey", "some value")
	stub.AssertCalled(t, "PutState", "missingkey", []byte("{\"value\":\"some value\"}"))
}

func TestReadGotG(t *testing.T) {
	var gotG *GotG
	var err error

	ctx, _ := configureStub()
	c := new(GotGContract)

	gotG, err = c.ReadGotG(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when reading")
	assert.Nil(t, gotG, "should not return GotG when exists errors when reading")

	gotG, err = c.ReadGotG(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when reading")
	assert.Nil(t, gotG, "should not return GotG when key does not exist in world state when reading")

	gotG, err = c.ReadGotG(ctx, "existingkey")
	assert.EqualError(t, err, "Could not unmarshal world state data to type GotG", "should error when data in key is not GotG")
	assert.Nil(t, gotG, "should not return GotG when data in key is not of type GotG")

	gotG, err = c.ReadGotG(ctx, "gotGkey")
	expectedGotG := new(GotG)
	expectedGotG.Value = "set value"
	assert.Nil(t, err, "should not return error when GotG exists in world state when reading")
	assert.Equal(t, expectedGotG, gotG, "should return deserialized GotG from world state")
}

func TestUpdateGotG(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(GotGContract)

	err = c.UpdateGotG(ctx, "statebad", "new value")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors when updating")

	err = c.UpdateGotG(ctx, "missingkey", "new value")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when updating")

	err = c.UpdateGotG(ctx, "gotGkey", "new value")
	expectedGotG := new(GotG)
	expectedGotG.Value = "new value"
	expectedGotGBytes, _ := json.Marshal(expectedGotG)
	assert.Nil(t, err, "should not return error when GotG exists in world state when updating")
	stub.AssertCalled(t, "PutState", "gotGkey", expectedGotGBytes)
}

func TestDeleteGotG(t *testing.T) {
	var err error

	ctx, stub := configureStub()
	c := new(GotGContract)

	err = c.DeleteGotG(ctx, "statebad")
	assert.EqualError(t, err, fmt.Sprintf("Could not read from world state. %s", getStateError), "should error when exists errors")

	err = c.DeleteGotG(ctx, "missingkey")
	assert.EqualError(t, err, "The asset missingkey does not exist", "should error when exists returns true when deleting")

	err = c.DeleteGotG(ctx, "gotGkey")
	assert.Nil(t, err, "should not return error when GotG exists in world state when deleting")
	stub.AssertCalled(t, "DelState", "gotGkey")
}
