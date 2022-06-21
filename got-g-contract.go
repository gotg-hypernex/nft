/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// GotGContract contract for managing CRUD for GotG
type GotGContract struct {
	contractapi.Contract
}

// GotGExists returns true when asset with given ID exists in world state
func (c *GotGContract) GotGExists(ctx contractapi.TransactionContextInterface, gotGID string) (bool, error) {
	data, err := ctx.GetStub().GetState(gotGID)

	if err != nil {
		return false, err
	}

	return data != nil, nil
}

// CreateGotG creates a new instance of GotG
func (c *GotGContract) CreateGotG(ctx contractapi.TransactionContextInterface, gotGID string, value string) error {
	exists, err := c.GotGExists(ctx, gotGID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if exists {
		return fmt.Errorf("The asset %s already exists", gotGID)
	}

	gotG := new(GotG)
	gotG.Value = value

	bytes, _ := json.Marshal(gotG)

	return ctx.GetStub().PutState(gotGID, bytes)
}

// ReadGotG retrieves an instance of GotG from the world state
func (c *GotGContract) ReadGotG(ctx contractapi.TransactionContextInterface, gotGID string) (*GotG, error) {
	exists, err := c.GotGExists(ctx, gotGID)
	if err != nil {
		return nil, fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return nil, fmt.Errorf("The asset %s does not exist", gotGID)
	}

	bytes, _ := ctx.GetStub().GetState(gotGID)

	gotG := new(GotG)

	err = json.Unmarshal(bytes, gotG)

	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal world state data to type GotG")
	}

	return gotG, nil
}

// UpdateGotG retrieves an instance of GotG from the world state and updates its value
func (c *GotGContract) UpdateGotG(ctx contractapi.TransactionContextInterface, gotGID string, newValue string) error {
	exists, err := c.GotGExists(ctx, gotGID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The asset %s does not exist", gotGID)
	}

	gotG := new(GotG)
	gotG.Value = newValue

	bytes, _ := json.Marshal(gotG)

	return ctx.GetStub().PutState(gotGID, bytes)
}

// DeleteGotG deletes an instance of GotG from the world state
func (c *GotGContract) DeleteGotG(ctx contractapi.TransactionContextInterface, gotGID string) error {
	exists, err := c.GotGExists(ctx, gotGID)
	if err != nil {
		return fmt.Errorf("Could not read from world state. %s", err)
	} else if !exists {
		return fmt.Errorf("The asset %s does not exist", gotGID)
	}

	return ctx.GetStub().DelState(gotGID)
}
