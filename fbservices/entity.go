package fbServices

//自定义区块结构体
//注意：DataHash并不是当前区块哈希，
//当前区块哈希的计算方式为区块头的三个字段（即number、previous_hash、data_hash）
//首先使用ASN.1中的DER编码规则进行编码，而后进行SHA256哈希值计算得出。
type Block struct {
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
