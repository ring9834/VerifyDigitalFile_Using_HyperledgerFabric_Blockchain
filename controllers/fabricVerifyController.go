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

type FabricVerifyController struct {
	beego.Controller
}

//for return use
type VerifyResult struct {
	Rst  int    `json:"rst"`
	Info string `json:"info"`
}

func (c *FabricVerifyController) VerifyArchReality() {
	f, h, _ := c.GetFile("xfile")
	fn := h.Filename
	strs := strings.Split(fn, ".")
	dh := strs[0]
	f.Close()
	topth := "static/uploads/forveriy/" + dh + "_" + time.Now().Format("2006_01_02_15_04_05") + "." + strs[1] //2006_01_02_15_04_05为go语言诞生的时间
	c.SaveToFile("xfile", topth)                                                                              //保存到的位置
	hashstr := utils.GetHash(topth)

	hashObj := ReadAssetByDhForVerify(dh)
	if hashstr == hashObj.ElectContracts {
		rslt := CreateAssetsResult{Rst: 1, Info: "恭喜！经区块链验证档案是真的！"} //flg=1 fail,1 success
		jsondata, _ := json.Marshal(rslt)                           // obj to json
		c.Data["json"] = string(jsondata)                           // return json data
		c.ServeJSON()                                               //返回JSON对象,供前端使用
	}
	rslt := VerifyResult{Rst: 0, Info: "对不起，经与区块链数据对比，档案有问题，请核查！"} //flg=0 fail,1 success
	jsondata, _ := json.Marshal(rslt)                              // obj to json
	c.Data["json"] = string(jsondata)                              // return json data
	c.ServeJSON()
}

type ElectArch struct {
	DH             string `json:"dh"`
	ElectContracts string `json:"hashstr"`
}

func ReadAssetByDhForVerify(dh string) ElectArch {
	InitialClient()
	searchRequest := &fabclient.ChaincodeRequest{
		ChaincodeID: chncode.Name,
		Function:    "ReadAsset",
		Args:        []string{dh},
	}
	storeResult, err := clnt.Invoke(searchRequest, fabclient.WithOrdererResponseTimeout(2*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	//var str = storeResult.Payload
	var archInfo ElectArch
	err2 := json.Unmarshal(storeResult.Payload, &archInfo)
	if err2 != nil {
		log.Printf("err info is : %s", err2.Error())
	}
	return archInfo
}
