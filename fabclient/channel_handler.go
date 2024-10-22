package fabclient

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type channelHandler interface {
	invoke(request *ChaincodeRequest, opts ...Option) (*TransactionResponse, error)
	query(request *ChaincodeRequest, opts ...Option) (*TransactionResponse, error)
	queryBlock(blockNumber uint64) (*Block, error)
	queryBlockByTxID(txID string) (*Block, error)
	queryBlockByHash(blockHash []byte) (*Block, error)
	queryInfo() (*BlockchainInfo, error)
	registerChaincodeEvent(chaincodeID, eventFilter string) (<-chan *ChaincodeEvent, error)
	unregisterChaincodeEvent(eventFilter string)

	//added on 20200501
	QueryLatestBlocksInfo() ([]*Block2, error)
	QueryLatestBlocksInfoJsonStr() (string, error)
	QueryBlockByBlockNumber(num int64) (*Block2, error)
	QueryTransactionByTxId(txId string) (*Transaction, error)
	QueryTransactionByTxIdJsonStr(txId string) (string, error)
}

type ongoingEvent struct {
	registration fab.Registration
	stopChan     chan chan struct{}
	wrapChan     chan *ChaincodeEvent
}

type channelHandlerClient struct { //channelHandlerClient realize the function interface of channelHandler 20220501
	client           *channel.Client
	eventManager     *event.Client
	underlyingLedger *ledger.Client

	chaincodeEvents map[string]*ongoingEvent
	mutex           sync.Mutex
}

func newChannelHandler(ctx context.ChannelProvider) (channelHandler, error) {
	channelClient, err := channel.New(ctx)
	if err != nil {
		return nil, err
	}

	eventManager, err := event.New(ctx, event.WithBlockEvents())
	if err != nil {
		return nil, err
	}

	ledgerClient, err := ledger.New(ctx)
	if err != nil {
		return nil, err
	}

	client := &channelHandlerClient{
		client:           channelClient,
		eventManager:     eventManager,
		underlyingLedger: ledgerClient,
		chaincodeEvents:  make(map[string]*ongoingEvent),
		mutex:            sync.Mutex{},
	}

	return client, nil
}

var _ channelHandler = (*channelHandlerClient)(nil)

func (chn *channelHandlerClient) invoke(request *ChaincodeRequest, opts ...Option) (*TransactionResponse, error) {
	response, err := chn.client.Execute(convertChaincodeRequest(request), convertOptions(opts...)...)
	return convertChaincodeTransactionResponse(response), err
}

func (chn *channelHandlerClient) query(request *ChaincodeRequest, opts ...Option) (*TransactionResponse, error) {
	response, err := chn.client.Query(convertChaincodeRequest(request), convertOptions(opts...)...)
	return convertChaincodeTransactionResponse(response), err
}

func (chn *channelHandlerClient) queryBlock(blockNumber uint64) (*Block, error) {
	block, err := chn.underlyingLedger.QueryBlock(blockNumber)
	return convertBlock(block), err
}
func (chn *channelHandlerClient) queryBlockByHash(blockHash []byte) (*Block, error) {
	block, err := chn.underlyingLedger.QueryBlockByHash(blockHash)
	return convertBlock(block), err
}

func (chn *channelHandlerClient) queryBlockByTxID(txID string) (*Block, error) {
	block, err := chn.underlyingLedger.QueryBlockByTxID(fab.TransactionID(txID))
	return convertBlock(block), err
}

func (chn *channelHandlerClient) queryInfo() (*BlockchainInfo, error) {
	blockchainInfo, err := chn.underlyingLedger.QueryInfo()
	return convertBlockchainInfo(blockchainInfo), err
}

func (chn *channelHandlerClient) registerChaincodeEvent(chaincodeID, eventFilter string) (<-chan *ChaincodeEvent, error) {
	chn.mutex.Lock()
	defer chn.mutex.Unlock()

	if _, ok := chn.chaincodeEvents[eventFilter]; ok {
		return nil, fmt.Errorf("event filter '%s' already registered", eventFilter)
	}

	registration, ch, err := chn.eventManager.RegisterChaincodeEvent(chaincodeID, eventFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to register chaincode event '%s' for chaincode '%s': %w", eventFilter, chaincodeID, err)
	}

	stopChan := make(chan chan struct{})
	wrapChan := make(chan *ChaincodeEvent)
	chn.chaincodeEvents[eventFilter] = &ongoingEvent{
		registration: registration,
		stopChan:     stopChan,
		wrapChan:     wrapChan,
	}

	go func() {
		for {
			select {
			case event := <-ch:
				wrapChan <- convertChaincodeEvent(event)
			case witness := <-stopChan:
				witness <- struct{}{}
				return
			}
		}
	}()

	return wrapChan, nil
}

func (chn *channelHandlerClient) unregisterChaincodeEvent(eventFilter string) {
	chn.mutex.Lock()
	defer chn.mutex.Unlock()

	if _, ok := chn.chaincodeEvents[eventFilter]; ok {
		witness := make(chan struct{})
		ongoingEvent := chn.chaincodeEvents[eventFilter]
		ongoingEvent.stopChan <- witness
		<-witness
		chn.eventManager.Unregister(ongoingEvent.registration)
		close(ongoingEvent.stopChan)
		close(ongoingEvent.wrapChan)
		delete(chn.chaincodeEvents, eventFilter)
	}

	return
}

func convertBlock(b *common.Block) *Block {
	if b == nil {
		return nil
	}

	header := &BlockHeader{
		Number:       b.GetHeader().Number,
		PreviousHash: b.GetHeader().PreviousHash,
		DataHash:     b.GetHeader().DataHash,
	}

	data := &BlockData{
		Data: b.GetData().Data,
	}

	metadata := &BlockMetadata{
		Metadata: b.GetMetadata().Metadata,
	}

	block := &Block{
		Header:   header,
		Data:     data,
		Metadata: metadata,
	}

	return block
}

func convertBlockchainInfo(info *fab.BlockchainInfoResponse) *BlockchainInfo {
	if info == nil {
		return nil
	}

	return &BlockchainInfo{
		info.BCI.Height,
		info.BCI.CurrentBlockHash,
		info.BCI.PreviousBlockHash,
	}
}

func convertChaincodeEvent(e *fab.CCEvent) *ChaincodeEvent {
	event := ChaincodeEvent(*e)
	return &event
}

func convertChaincodeRequest(request *ChaincodeRequest) channel.Request {
	if request == nil {
		return channel.Request{}
	}

	invocationChain := make([]*fab.ChaincodeCall, 0, len(request.InvocationChain))
	for _, invoc := range request.InvocationChain {
		invocationChain = append(invocationChain, &fab.ChaincodeCall{
			ID:          invoc.ID,
			Collections: invoc.Collections,
		})
	}

	return channel.Request{
		Args:            convertArrayOfStringsToArrayOfByteArrays(request.Args),
		Fcn:             request.Function,
		ChaincodeID:     request.ChaincodeID,
		TransientMap:    request.TransientMap,
		InvocationChain: invocationChain,
		IsInit:          request.IsInit,
	}
}

func convertChaincodeTransactionResponse(response channel.Response) *TransactionResponse {
	return &TransactionResponse{
		Payload:       response.Payload,
		Status:        response.ChaincodeStatus,
		TransactionID: string(response.TransactionID),
	}
}

func convertOptions(opts ...Option) []channel.RequestOption {
	convertedOpts := make([]channel.RequestOption, 0, len(opts))

	o := &options{
		ordererResponseTimeout: -1,
	}

	for _, opt := range opts {
		opt.apply(o)
	}

	if o.ordererResponseTimeout != -1 {
		convertedOpts = append(convertedOpts, channel.WithTimeout(fab.OrdererResponse, o.ordererResponseTimeout))
	}

	return convertedOpts
}

//////////////////////////////////20220501///////////////////////////////////////
//查询最新5个区块信息
func (chn *channelHandlerClient) QueryLatestBlocksInfo() ([]*Block2, error) {
	ledgerInfo, err := chn.underlyingLedger.QueryInfo()
	if err != nil {
		fmt.Printf("QueryLatestBlocksInfo return error: %s\n", err)
		return nil, err
	}
	latestBlockList := []*Block2{}
	lastetBlockNum := ledgerInfo.BCI.Height - 1

	for i := lastetBlockNum; i > 0 && i > (lastetBlockNum-5); i-- {
		block, err := chn.QueryBlockByBlockNumber(int64(i))
		if err != nil {
			fmt.Printf("QueryLatestBlocksInfo return error: %s", err)
			return latestBlockList, err
		}
		latestBlockList = append(latestBlockList, block)
	}
	return latestBlockList, nil
}

func (chn *channelHandlerClient) QueryLatestBlocksInfoJsonStr() (string, error) {
	blockList, err := chn.QueryLatestBlocksInfo()
	jsonStr, err := json.Marshal(blockList)
	return string(jsonStr), err
}

//查询指定区块信息
func (chn *channelHandlerClient) QueryBlockByBlockNumber(num int64) (*Block2, error) {
	rawBlock, err := chn.underlyingLedger.QueryBlock(uint64(num))
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

	block := Block2{
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
func (chn *channelHandlerClient) QueryTransactionByTxId(txId string) (*Transaction, error) {
	rawTx, err := chn.underlyingLedger.QueryTransaction(fab.TransactionID(txId))
	if err != nil {
		fmt.Printf("QueryBlock return error: %s", err)
		return nil, err
	}

	transaction, err := GetTransactionFromEnvelopeDeep(rawTx.TransactionEnvelope)
	if err != nil {
		fmt.Printf("QueryBlock return error: %s", err)
		return nil, err
	}
	block, err := chn.underlyingLedger.QueryBlockByTxID(fab.TransactionID(txId))
	if err != nil {
		fmt.Printf("QueryBlock return error: %s", err)
		return nil, err
	}
	for i := range transaction.TransactionActionList {
		transaction.TransactionActionList[i].BlockNum = block.Header.Number
	}
	return transaction, nil
}

func (chn *channelHandlerClient) QueryTransactionByTxIdJsonStr(txId string) (string, error) {
	transaction, err := chn.QueryTransactionByTxId(txId)
	if err != nil {
		return "", err
	}
	jsonStr, err := json.Marshal(transaction)
	return string(jsonStr), err
}

//////////////////////////////////////////////////////////////////////////////
////////////////////////////////codes below are BLOCK utils//////20220501/////////////////////
//json Marshal/Unmarshal解析工具
//自定义解析类型

//获取区块体数据item
func GetEnvelopeFromBlock(data []byte) (*common.Envelope, error) {
	var err error
	env := &common.Envelope{}
	if err = proto.Unmarshal(data, env); err != nil {
		fmt.Printf("block unmarshal err: %s", err)
	}

	return env, nil
}

//解析Transaction深度递归
func GetTransactionFromEnvelopeDeep(rawEnvelope *common.Envelope) (*Transaction, error) {
	//解析payload
	rawPayload := &common.Payload{}
	err := proto.Unmarshal(rawEnvelope.Payload, rawPayload)
	if err != nil {
		fmt.Printf("block unmarshal err: %s", err)
	}
	//解析channelHeader
	channelHeader := &common.ChannelHeader{}
	err = proto.Unmarshal(rawPayload.Header.ChannelHeader, channelHeader)
	if err != nil {
		fmt.Printf("block unmarshal err: %s\n", err)
	}
	//解析Transaction
	transactionObj := &peer.Transaction{}
	err = proto.Unmarshal(rawPayload.Data, transactionObj)
	if err != nil {
		fmt.Printf("block unmarshal err: %s\n", err)
	}
	transactionActionList := []*TransactionAction{}
	//解析transactionAction
	for i := range transactionObj.Actions {
		transactionAction, err := GetTransactionActionFromTransactionDeep(transactionObj.Actions[i])
		if err != nil {
			fmt.Printf("block unmarshal err: %s\n", err)
		}
		transactionAction.TxId = channelHeader.TxId
		transactionAction.Type = string(channelHeader.Type)
		transactionAction.Timestamp = time.Unix(channelHeader.Timestamp.Seconds, 0).Format("2006-01-02 15:04:05")
		transactionAction.ChannelId = channelHeader.ChannelId
		transactionActionList = append(transactionActionList, transactionAction)
	}
	transaction := Transaction{transactionActionList}

	return &transaction, nil
}

func GetTransactionActionFromTransactionDeep(transactionAction *peer.TransactionAction) (*TransactionAction, error) {

	//解析ChaincodeActionPayload 1
	ChaincodeActionPayload := &peer.ChaincodeActionPayload{}
	err := proto.Unmarshal(transactionAction.Payload, ChaincodeActionPayload)
	if err != nil {
		fmt.Printf("block unmarshal err: %s\n", err)
	}

	//解析ProposalResponsePayload 1.2
	ProposalResponsePayload := &peer.ProposalResponsePayload{}
	ChaincodeAction := &peer.ChaincodeAction{}
	chaincodeId := ""
	NsReadWriteSetList := []*rwset.NsReadWriteSet{}
	ReadWriteSetList := []*kvrwset.KVRWSet{}
	readSetList := []string{}
	writeSetList := []string{}
	if ChaincodeActionPayload.GetAction() != nil {
		err = proto.Unmarshal(ChaincodeActionPayload.Action.ProposalResponsePayload, ProposalResponsePayload)
		if err != nil {
			fmt.Printf("block unmarshal err: %s", err)
		}
		//解析ChaincodeAction 1.2.1
		err = proto.Unmarshal(ProposalResponsePayload.Extension, ChaincodeAction)
		if err != nil {
			fmt.Printf("block unmarshal err: %s", err)
		}

		chaincodeId = ChaincodeAction.ChaincodeId.Name
		//解析TxReadWriteSet	1.2.1.1
		TxReadWriteSet := &rwset.TxReadWriteSet{}
		err = proto.Unmarshal(ChaincodeAction.Results, TxReadWriteSet)
		if err != nil {
			fmt.Printf("block unmarshal err: %s", err)
		}
		//解析TxReadWriteSet	1.2.1.1.1
		for i := range TxReadWriteSet.NsRwset {
			ReadWriteSet := &kvrwset.KVRWSet{}
			//解析ReadWriteSet	1.2.1.1.1.1
			err = proto.Unmarshal(TxReadWriteSet.NsRwset[i].Rwset, ReadWriteSet)
			if err != nil {
				fmt.Printf("block unmarshal err: %s", err)
			}

			//解析读集
			for i := range ReadWriteSet.Reads {
				readSetJsonStr, err := json.Marshal(ReadWriteSet.Reads[i])
				if err != nil {
					fmt.Printf("block unmarshal err: %s", err)
				}
				readSetList = append(readSetList, string(readSetJsonStr))
			}

			//解析写集
			for i := range ReadWriteSet.Writes {
				writeSetItem := map[string]interface{}{
					"Key":      ReadWriteSet.Writes[i].GetKey(),
					"Value":    string(ReadWriteSet.Writes[i].GetValue()),
					"IsDelete": ReadWriteSet.Writes[i].GetIsDelete(),
				}

				writeSetJsonStr, err := json.Marshal(writeSetItem)
				if err != nil {
					fmt.Printf("block unmarshal err: %s", err)
				}
				writeSetList = append(writeSetList, string(writeSetJsonStr))

			}

			ReadWriteSetList = append(ReadWriteSetList, ReadWriteSet)
			NsReadWriteSetList = append(NsReadWriteSetList, TxReadWriteSet.NsRwset[i])
		}

	} else {
		chaincodeId = "没有交易数据"
	}

	//log.Println("数据:"+fmt.Sprintf("%s\n",ChaincodeActionPayload.Action.GetEndorsements()[0] ))
	//解析Endorsements 1.3
	endorsements := []string{}
	if ChaincodeActionPayload.Action.GetEndorsements() != nil {
		for i := range ChaincodeActionPayload.Action.GetEndorsements() {
			endorser := &msp.SerializedIdentity{}
			err = proto.Unmarshal(ChaincodeActionPayload.Action.Endorsements[i].Endorser, endorser)
			if err != nil {
				fmt.Printf("block unmarshal err: %s", err)
			}

			endorsements = append(endorsements, string(endorser.Mspid))
		}
	}

	transactionActionObj := TransactionAction{
		Endorsements: endorsements,
		ChaincodeId:  chaincodeId,
		ReadSetList:  readSetList,
		WriteSetList: writeSetList,
	}

	return &transactionActionObj, nil
}
