package fabclient

// BlockData holds the transactions.
type BlockData struct {
	Data [][]byte
}

// BlockHeader is the element of the block which forms the blockchain.
type BlockHeader struct {
	Number       uint64
	PreviousHash []byte
	DataHash     []byte
}

// BlockMetadata defines metadata of the block.
type BlockMetadata struct {
	Metadata [][]byte
}

// Block is finalized block structure to be shared among the orderer and peer.
type Block struct {
	Header   *BlockHeader
	Data     *BlockData
	Metadata *BlockMetadata
}

// BlockchainInfo contains information about the blockchain ledger such as height, current block hash, and previous block hash.
type BlockchainInfo struct {
	Height            uint64
	CurrentBlockHash  []byte
	PreviousBlockHash []byte
}

// Chaincode describes info of a chaincode.
type Chaincode struct {
	Collections          []ChaincodeCollection `json:"collections,omitempty" yaml:"collections,omitempty"`
	InitRequired         bool                  `json:"initRequired" yaml:"initRequired"`
	MustBeApprovedByOrgs []string              `json:"mustBeApprovedByOrgs" yaml:"mustBeApprovedByOrgs"`
	Name                 string                `json:"name" yaml:"name"`
	Path                 string                `json:"path" yaml:"path"`
	Role                 string                `json:"role" yaml:"role"`
	Sequence             int64                 `json:"sequence" yaml:"sequence"`
	Version              string                `json:"version" yaml:"version"`
}

// ChaincodeCall contains the ID of the chaincode as well as an optional set of private data collections that may be accessed by the chaincode.
type ChaincodeCall struct {
	ID          string
	Collections []string
}

// ChaincodeCollection defines the configuration of a collection.
type ChaincodeCollection struct {
	BlockToLive       uint64 `json:"blockToLive" yaml:"blockToLive"`
	MaxPeerCount      int32  `json:"maxPeerCount" yaml:"maxPeerCount"`
	MemberOnlyRead    bool   `json:"memberOnlyRead" yaml:"memberOnlyRead"`
	Name              string `json:"name" yaml:"name"`
	Policy            string `json:"policy" yaml:"policy"`
	RequiredPeerCount int32  `json:"requiredPeerCount" yaml:"requiredPeerCount"`
}

// ChaincodeEvent contains the data for a chaincode event.
type ChaincodeEvent struct {
	TxID        string
	ChaincodeID string
	EventName   string
	Payload     []byte
	BlockNumber uint64
	SourceURL   string
}

// ChaincodeRequest contains the parameters to query and execute an invocation transaction.
type ChaincodeRequest struct {
	ChaincodeID     string
	Function        string
	Args            []string
	TransientMap    map[string][]byte
	InvocationChain []*ChaincodeCall
	IsInit          bool
}

// Channel describes a channel configuration.
type Channel struct {
	AnchorPeerConfigPath string `json:"anchorPeerConfigPath,omitempty" yaml:"anchorPeerConfigPath,omitempty"`
	ConfigPath           string `json:"configPath" yaml:"configPath"`
	Name                 string `json:"name" yaml:"name"`
}

// Identity holds crypto material for creating a signing identity.
type Identity struct {
	Certificate string `json:"certificate" yaml:"certificate"`
	PrivateKey  string `json:"privateKey" yaml:"privateKey"`
	Username    string `json:"username" yaml:"username"`
}

// TransactionResponse  contains response parameters for query and execute an invocation transaction.
type TransactionResponse struct {
	Payload       []byte
	Status        int32
	TransactionID string
}

//自定义区块结构体 codes below added on 20220501
//注意：DataHash并不是当前区块哈希，
//当前区块哈希的计算方式为区块头的三个字段（即number、previous_hash、data_hash）
//首先使用ASN.1中的DER编码规则进行编码，而后进行SHA256哈希值计算得出。
type Block2 struct {
	Number          uint64         `json:"number"`          //区块号
	PreviousHash    []byte         `json:"previousHash"`    //前区块Hash
	DataHash        []byte         `json:"dataHash"`        //交易体Hash
	BlockHash       []byte         `json:"blockHash"`       //区块Hash
	TxNum           int            `json:"txNum"`           //区块内交易个数
	TransactionList []*Transaction `json:"transactionList"` //交易列表
	CreateTime      string         `json:"createTime"`      //区块生成时间
}
type Transaction struct {
	TransactionActionList []*TransactionAction `json:"transactionActionList"` //交易列表
}

type TransactionAction struct {
	TxId         string   `json:"txId"`         //交易ID
	BlockNum     uint64   `json:"blockNum"`     //区块号
	Type         string   `json:"type"`         //交易类型
	Timestamp    string   `json:"timestamp"`    //交易创建时间
	ChannelId    string   `json:"channelId"`    //通道ID
	Endorsements []string `json:"endorsements"` //背书组织ID列表
	ChaincodeId  string   `json:"chaincodeId"`  //链代码名称
	ReadSetList  []string `json:"readSetList"`  //读集
	WriteSetList []string `json:"writeSetList"` //写集
}
