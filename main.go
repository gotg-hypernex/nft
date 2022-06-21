/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-contract-api-go/metadata"
)

func main() {
	gotGContract := new(GotGContract)
	gotGContract.Info.Version = "0.0.1"
	gotGContract.Info.Description = "My Smart Contract"
	gotGContract.Info.License = new(metadata.LicenseMetadata)
	gotGContract.Info.License.Name = "Apache-2.0"
	gotGContract.Info.Contact = new(metadata.ContactMetadata)
	gotGContract.Info.Contact.Name = "John Doe"

	chaincode, err := contractapi.NewChaincode(gotGContract)
	chaincode.Info.Title = "src chaincode"
	chaincode.Info.Version = "0.0.1"

	if err != nil {
		panic("Could not create chaincode from GotGContract." + err.Error())
	}

	err = chaincode.Start()

	if err != nil {
		panic("Failed to start chaincode. " + err.Error())
	}
}
