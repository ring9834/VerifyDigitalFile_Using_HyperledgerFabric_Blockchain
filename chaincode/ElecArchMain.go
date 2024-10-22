package main

import (
	"hzx/chaincode/archchaincode"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	assetChaincode, err := contractapi.NewChaincode(&archchaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating electronic-archive chaincode: %v", err)
	}

	if err := assetChaincode.Start(); err != nil {
		log.Panicf("Error starting electronic-archive chaincode: %v", err)
	}
}
