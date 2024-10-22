package archchaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type ElectArch struct {
	DH             string `json:"dh"`
	ElectContracts string `json:"hashstr"`
}

type ElectArchsWithBookmark struct {
	RecordsCount int32        `json:"total"`
	Bookmark     string       `json:"bookmark"`
	ElcArches    []*ElectArch `json:"rows"`
}

// InitLedger adds a base set of ElectArch to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	ela := []ElectArch{
		{DH: "0024-001-025-003", ElectContracts: "0f33e0db32d345ba2b60b61e8fda48026db826040cc65069e87c6bb61295aa61"},
	}

	for _, asset := range ela {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.DH, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, dh string, electContracts string) error {
	exists, err := s.AssetExists(ctx, dh)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", dh)
	}

	asset := ElectArch{
		DH:             dh,
		ElectContracts: electContracts,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(dh, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, dh string) (*ElectArch, error) {
	assetJSON, err := ctx.GetStub().GetState(dh)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", dh)
	}

	var asset ElectArch
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, dh string, electContracts string) error {
	exists, err := s.AssetExists(ctx, dh)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", dh)
	}

	// overwriting original asset with new asset
	asset := ElectArch{
		DH:             dh,
		ElectContracts: electContracts,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(dh, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, dh string) error {
	exists, err := s.AssetExists(ctx, dh)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", dh)
	}

	return ctx.GetStub().DelState(dh)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, dh string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(dh)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*ElectArch, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*ElectArch
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset ElectArch
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

// GetAssets returns assets found in world state with pagination
func (s *SmartContract) GetAssetsWithPagination(ctx contractapi.TransactionContextInterface, startKey string, endKey string, pageSize int32, bookmark string) (*ElectArchsWithBookmark, error) {
	// pgs, err := strconv.Atoi(pageSize) //return int64
	// if err != nil {
	// 	return nil, err
	// }
	// pgs2 := int32(pgs) //convert to int32
	resultsIterator, queryResponseMetadata, err := ctx.GetStub().GetStateByRangeWithPagination(startKey, endKey, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var bm = queryResponseMetadata.Bookmark
	var cnt = queryResponseMetadata.FetchedRecordsCount
	var assets []*ElectArch
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset ElectArch
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	rst := ElectArchsWithBookmark{
		cnt,
		bm,
		assets,
	}
	return &rst, nil
}
