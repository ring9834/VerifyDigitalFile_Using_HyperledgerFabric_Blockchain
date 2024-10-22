package controllers

import (
	"encoding/json"
	"hzx/fabclient"
	"hzx/utils"
	"log"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type FabricHashController struct {
	beego.Controller
}

var ( // globle var
	clnt    *fabclient.Client
	chnnl   fabclient.Channel
	chncode fabclient.Chaincode
)

func (c *FabricHashController) ExtracthashSrvView() {
	c.TplName = "ExtracthashSvr.html"
}

func (c *FabricHashController) ExtracthashClntView() {
	c.TplName = "ExtracthashClnt.html"
}

func (c *FabricHashController) VerifyRealityView() {
	c.TplName = "VerifyFileReality.html"
}

var HashArr = make(map[string]string)

//如果是一次上传多个文件，则每个文件请求一次controller
func (c *FabricHashController) GetExtractedHashStr() {
	f, h, _ := c.GetFile("xfile")
	fn := h.Filename
	strs := strings.Split(fn, ".")
	dh := strs[0]
	f.Close()
	//topth := "static/uploads/yws/" + fn
	topth := "static/uploads/yws/" + dh + "_" + time.Now().Format("2006_01_02_15_04_05") + "." + strs[1] //2006_01_02_15_04_05为go语言诞生的时间
	c.SaveToFile("xfile", topth)                                                                         //保存到的位置
	hashstr := utils.GetHash(topth)
	HashArr[dh] = hashstr
}

func InitialClient() {
	if clnt == nil {
		clnt, _ = fabclient.NewClientFromConfigFile("./fabclient/testdata/organizations/org1/client-config.yaml")
		chncode = clnt.Config().Chaincodes[0]
		chnnl = clnt.Config().Channels[0]
		clnt.CreateChannelHandler(chnnl.Name) // must be create channle handler before Invoke 20220501
	}
}

//for return use
type CreateAssetsResult struct {
	Rst  int    `json:"rst"`
	Info string `json:"info"`
}

//一次接收多个文件，仅请求一次controller
func (c *FabricHashController) CreateArchAssets() {
	f, h, _ := c.GetFile("xfile")
	fn := h.Filename
	strs := strings.Split(fn, ".")
	dh := strs[0]
	f.Close()
	topth := "static/uploads/yws/" + dh + "_" + time.Now().Format("2006_01_02_15_04_05_999999999") + "." + strs[1] //2006_01_02_15_04_05为go语言诞生的时间
	c.SaveToFile("xfile", topth)                                                                                   //保存到的位置
	hashstr := utils.GetHash(topth)
	//此处可以考虑删除已上传的档案文件，以节约服务器空间
	flg, info := AddArchInofToChain(dh, hashstr)
	rslt := CreateAssetsResult{Rst: flg, Info: info} //flg=0 fail,1 success
	jsondata, _ := json.Marshal(rslt)                // obj to json
	c.Data["json"] = string(jsondata)                // return json data
	c.ServeJSON()                                    //返回JSON对象,供前端使用
}

//档号和档案摘要入链20220425
func AddArchInofToChain(dh string, hashstr string) (int, string) {
	InitialClient()
	//一件档案一件档案地入链
	addRequest := &fabclient.ChaincodeRequest{
		ChaincodeID: chncode.Name,
		Function:    "CreateAsset",
		Args:        []string{dh, hashstr},
	}
	storeResult, err := clnt.Invoke(addRequest, fabclient.WithOrdererResponseTimeout(2*time.Second)) //调用智能合约里的函数CreateAsset
	if err != nil {
		//log.Fatal(err)
		var info = ""
		if strings.Index(err.Error(), "exist") == -1 {
			info = dh + ":档案存入区块链失败 " + time.Now().Format("2006-01-02 15:04:05")
		} else {
			info = dh + ":此档案已在区块链存在，存入失败 " + time.Now().Format("2006-01-02 15:04:05")
		}
		return 0, info // fail
	} else {
		log.Printf("store txID: %s", storeResult.TransactionID)
		return 1, dh + ":档案存入区块链成功 " + time.Now().Format("2006-01-02 15:04:05")
	}
}

func (c *FabricHashController) ConnectArchChainCode() {
	client, err := fabclient.NewClientFromConfigFile("./fabclient/testdata/organizations/org1/client-config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	channel := client.Config().Channels[0]
	chaincode := client.Config().Chaincodes[0]

	// if err := client.SaveChannel(channel.Name, channel.ConfigPath); err != nil {
	// 	log.Fatal(err)
	// }

	if err := client.JoinChannel(channel.Name); err != nil {
		log.Fatal(err)
	}

	chaincodePackageID, err := client.LifecycleInstallChaincode(chaincode)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.LifecycleApproveChaincode(channel.Name, chaincodePackageID, chaincode); err != nil {
		log.Fatal(err)
	}

	if err := client.LifecycleCommitChaincode(channel.Name, chaincode); err != nil {
		log.Fatal(err)
	}

	// Invoke/Query chaincode using default API
	storeRequest := &fabclient.ChaincodeRequest{
		ChaincodeID: chaincode.Name,
		Function:    "InitLedger",
		Args:        []string{},
	}

	storeResult, err := client.Invoke(storeRequest, fabclient.WithOrdererResponseTimeout(2*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("store txID: %s", storeResult.TransactionID)

	// queryRequest := &fabclient.ChaincodeRequest{
	// 	ChaincodeID: chaincode.Name,
	// 	Function:    "GetAllAssets",
	// 	Args:        []string{},
	// }

	// queryResult, err := client.Query(queryRequest)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Printf("query content: %s", string(queryResult.Payload))

	// // using Gateway
	// user := client.Config().Identities.Users[0]

	// cert, err := ioutil.ReadFile(user.Certificate)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// pk, err := ioutil.ReadFile(user.PrivateKey)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// id := fabclient.NewWalletX509Identity("Org1MSP", string(cert), string(pk))

	// w, err := fabclient.NewFileSystemWallet("./testdata/test-wallet")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if err := w.Put(user.Username, id); err != nil {
	// 	log.Fatal(err)
	// }

	// gtw, err := fabclient.Connect(fabclient.WithConfigFromFile(client.Config().ConnectionProfile), fabclient.WithIdentity(w, user.Username))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// network, err := gtw.GetNetwork(channel.Name)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// contract := network.GetContract(chaincode.Name)

	// if _, err := contract.SubmitTransaction("Store", []string{"gateway-asset-test", `{"content": "this is another content test"}`}); err != nil {
	// 	log.Fatal(err)
	// }

	// result, err := contract.EvaluateTransaction("Query", []string{"gateway-asset-test"})
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println(string(result))

}
