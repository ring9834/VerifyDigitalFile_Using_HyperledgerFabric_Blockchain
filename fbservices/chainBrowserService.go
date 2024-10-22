package fbServices

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/util/pathvar"
)

//区块链浏览器服务

var mainSDK *fabsdk.FabricSDK
var ledgerClient *ledger.Client

const (
	org1Name         = "Org1"
	org2Name         = "Org2"
	org1Peer0        = "peer0.org1.example.com"
	org1AdminUser    = "Admin"
	org2AdminUser    = "Admin"
	org1User         = "User1"
	org2User         = "User1"
	channelID        = "channel1"
	windowConfigPath = "D:/GoProject/src/github.com/hyperledger/fabric/3ArchGoChain/config/config.yaml"
	linuxConfigPath  = "/opt/gopath/src/github.com/hyperledger/fabric/3ArchGoChain/config/config.yaml"
)

var chainBrowserConfigPath = linuxConfigPath

//初始化区块浏览器SDK
func InitChainBrowserService() {
	log.Println("============ 初始化区块浏览器服务 ============")
	//获取fabsdk
	var err error
	ConfigBackend := config.FromFile(pathvar.Subst(chainBrowserConfigPath))
	mainSDK, err = fabsdk.New(ConfigBackend)
	if err != nil {
		panic(fmt.Sprintf("Failed to create new SDK: %s", err))
	}
	//获取context
	org1AdminChannelContext := mainSDK.ChannelContext(channelID, fabsdk.WithUser(org1AdminUser), fabsdk.WithOrg(org1Name))
	//Ledger client
	ledgerClient, err = ledger.New(org1AdminChannelContext)
	if err != nil {
		fmt.Printf("Failed to create new resource management client: %s", err)
	}
}

//查询账本信息
func QueryLedgerInfo() (*fab.BlockchainInfoResponse, error) {
	ledgerInfo, err := ledgerClient.QueryInfo()
	if err != nil {
		fmt.Printf("QueryInfo return error: %s", err)
		return nil, err
	}
	return ledgerInfo, nil
}

//查询最新5个区块信息
func QueryLatestBlocksInfo() ([]*Block, error) {
	ledgerInfo, err := ledgerClient.QueryInfo()
	if err != nil {
		fmt.Printf("QueryLatestBlocksInfo return error: %s\n", err)
		return nil, err
	}
	latestBlockList := []*Block{}
	lastetBlockNum := ledgerInfo.BCI.Height - 1

	for i := lastetBlockNum; i > 0 && i > (lastetBlockNum-5); i-- {
		block, err := QueryBlockByBlockNumber(int64(i))
		if err != nil {
			fmt.Printf("QueryLatestBlocksInfo return error: %s", err)
			return latestBlockList, err
		}
		latestBlockList = append(latestBlockList, block)
	}
	return latestBlockList, nil
}

func QueryLatestBlocksInfoJsonStr() (string, error) {
	blockList, err := QueryLatestBlocksInfo()
	jsonStr, err := json.Marshal(blockList)
	return string(jsonStr), err
}

//查询指定区块信息
func QueryBlockByBlockNumber(num int64) (*Block, error) {
	rawBlock, err := ledgerClient.QueryBlock(uint64(num))
	if err != nil {
		fmt.Printf("QueryBlock return error: %s", err)
		return nil, err
	}

	//解析区块体
	txList := []*Transaction{}
	for i := range rawBlock.Data.Data {
		rawEnvelope, err := GetEnvelopeFromBlock(rawBlock.Data.Data[i])
		if err != nil {
			fmt.Printf("QueryBlock return error: %s", err)
			return nil, err
		}
		transaction, err := GetTransactionFromEnvelopeDeep(rawEnvelope)
		if err != nil {
			fmt.Printf("QueryBlock return error: %s", err)
			return nil, err
		}
		for i := range transaction.TransactionActionList {
			transaction.TransactionActionList[i].BlockNum = rawBlock.Header.Number
		}
		txList = append(txList, transaction)
	}

	block := Block{
		Number:          rawBlock.Header.Number,
		PreviousHash:    rawBlock.Header.PreviousHash,
		DataHash:        rawBlock.Header.DataHash,
		BlockHash:       rawBlock.Header.DataHash, //需要计算
		TxNum:           len(rawBlock.Data.Data),
		TransactionList: txList,
		CreateTime:      txList[0].TransactionActionList[0].Timestamp,
	}

	return &block, nil
}

//查询交易信息
func QueryTransactionByTxId(txId string) (*Transaction, error) {
	rawTx, err := ledgerClient.QueryTransaction(fab.TransactionID(txId))
	if err != nil {
		fmt.Printf("QueryBlock return error: %s", err)
		return nil, err
	}

	transaction, err := GetTransactionFromEnvelopeDeep(rawTx.TransactionEnvelope)
	if err != nil {
		fmt.Printf("QueryBlock return error: %s", err)
		return nil, err
	}
	block, err := ledgerClient.QueryBlockByTxID(fab.TransactionID(txId))
	if err != nil {
		fmt.Printf("QueryBlock return error: %s", err)
		return nil, err
	}
	for i := range transaction.TransactionActionList {
		transaction.TransactionActionList[i].BlockNum = block.Header.Number
	}
	return transaction, nil
}

func QueryTransactionByTxIdJsonStr(txId string) (string, error) {
	transaction, err := QueryTransactionByTxId(txId)
	if err != nil {
		return "", err
	}
	jsonStr, err := json.Marshal(transaction)
	return string(jsonStr), err
}
