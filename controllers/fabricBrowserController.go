package controllers

import (
	"encoding/json"
	"hzx/fabclient"
	"log"
	"time"

	"github.com/astaxie/beego"
)

type FabricBrowserController struct {
	beego.Controller
}

func (c *FabricBrowserController) ChainBrowserView() {
	c.TplName = "BlockInfoBrowse.html"
}

//get previouse block and current block hash info, and block height 20220501
func (c *FabricBrowserController) QueryInfo() {
	InitialClient()
	info, err := clnt.QueryInfo()
	if err != nil {
		log.Fatal(err)
	}
	c.Data["json"] = info
	c.ServeJSON() //返回JSON对象,供前端使用
}

// get latest 5 blocks
func (c *FabricBrowserController) QueryLatestBlocksInfo(expected int64) {
	InitialClient()
	blocks, err := clnt.QueryLatestBlocksInfo()
	if err != nil {
		log.Fatal(err)
	}
	c.Data["json"] = blocks
	c.ServeJSON() //返回JSON对象,供前端使用
}

//get the block by Num
func (c *FabricBrowserController) QueryBlockByBlockNumber(blocknum int64) {
	InitialClient()
	blockbynum, err := clnt.QueryBlockByBlockNumber(blocknum)
	if err != nil {
		log.Fatal(err)
	}
	c.Data["json"] = blockbynum
	c.ServeJSON() //返回JSON对象,供前端使用
}

// get all assets, can be invoked successfully
func (c *FabricBrowserController) ReadAllAssets() {
	InitialClient()
	searchRequest := &fabclient.ChaincodeRequest{
		ChaincodeID: chncode.Name,
		Function:    "GetAllAssets",
		Args:        []string{},
	}
	storeResult, err := clnt.Invoke(searchRequest, fabclient.WithOrdererResponseTimeout(2*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("store txID: %s", storeResult.TransactionID)
	var archInfos []ElectArch
	json.Unmarshal(storeResult.Payload, &archInfos)
	c.Data["json"] = archInfos
	//c.Data["json"] = string(storeResult.Payload)
	c.ServeJSON() //返回JSON对象,供前端使用
}

type ElectArchsWithBookmark struct {
	RecordsCount int32        `json:"total"`
	Bookmark     string       `json:"bookmark"`
	ElcArches    []*ElectArch `json:"rows"`
}

type TableArgs struct {
	PageSize string `json:"pageSize"`
	Bookmark string `json:"bookmark"`
}

// get assets with pagination20220503 pageSize string, bookmark string
func (c *FabricBrowserController) ReadPagedAssets() {
	acceptedParam := c.Ctx.Input.RequestBody //注意：从前端传过来的JSON格式要用data:'{}',一定要用单引号，.net core中不用单引号
	var param TableArgs
	merr := json.Unmarshal(acceptedParam, &param)
	if merr != nil {
		log.Printf("JSON参数解析有误" + merr.Error())
	}

	InitialClient()
	searchRequest := &fabclient.ChaincodeRequest{
		ChaincodeID: chncode.Name,
		Function:    "GetAssetsWithPagination",
		Args: []string{
			"",
			"",
			param.PageSize,
			param.Bookmark,
		},
	}
	storeResult, err := clnt.Invoke(searchRequest, fabclient.WithOrdererResponseTimeout(2*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	var archInfos ElectArchsWithBookmark
	json.Unmarshal(storeResult.Payload, &archInfos)
	log.Printf("store txID: %s", storeResult.TransactionID)
	//jsondata, _ := json.Marshal(archInfos)
	c.Data["json"] = archInfos
	//c.Data["json"] = string(storeResult.Payload)
	c.ServeJSON() //返回JSON对象,供前端使用
}

//get one  asset by dh
func (c *FabricBrowserController) ReadAssetByDh(dh string) {
	InitialClient()
	searchRequest := &fabclient.ChaincodeRequest{
		ChaincodeID: chncode.Name,
		Function:    "ReadAsset",
		Args:        []string{dh},
	}
	storeResult, err := clnt.Invoke(searchRequest, fabclient.WithOrdererResponseTimeout(2*time.Second)) //调用智能合约里的函数CreateAsset
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("store txID: %s", storeResult.TransactionID)
	c.Data["json"] = storeResult.Payload
	c.ServeJSON() //返回JSON对象,供前端使用
}
