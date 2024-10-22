package controllers

import (
	"encoding/json"
	"hzx/utils"
	"io/ioutil"
	"log"
	"path"
	"strings"

	"github.com/astaxie/beego"
)

type FabricSvrHashController struct {
	beego.Controller
}

//档案原文从服务端入链
func (c *FabricSvrHashController) CreateAssetsFromServer() {
	acceptedParam := c.Ctx.Input.RequestBody //注意：从前端传过来的JSON格式要用data:'{}',一定要用单引号，.net core中不用单引号
	var pathparam pth
	merr := json.Unmarshal(acceptedParam, &pathparam)
	if merr != nil {
		log.Printf("JSON参数解析有误" + merr.Error())
	}

	results := []CreateAssetsResult{}

	f, _ := GetAllFile(pathparam.Mypath)
	for i := 0; i < len(f); i++ {
		hashstr, filename := utils.GetHash2(f[i])
		strs := strings.Split(filename, ".")
		dh := strs[0]
		flg, info := AddArchInofToChain(dh, hashstr)
		rslt := CreateAssetsResult{Rst: flg, Info: info} //flg=0 fail,1 success
		results = append(results, rslt)
	}

	jsondata, _ := json.Marshal(results) // obj to json
	c.Data["json"] = string(jsondata)    // return json data
	c.ServeJSON()
}

var extarray = []string{
	".pdf", ".ofd", ".jpg", ".jpeg", ".png", ".bmp", ".tif", ".tiff", ".doc", ".docx", ".wps", ".xls", ".xlsx", ".ppt", ".pptx",
}

// 递归获取指定目录下的所有文件名
func GetAllFile(pathname string) ([]string, error) {
	result := []string{}

	fis, err := ioutil.ReadDir(pathname)
	if err != nil {
		log.Printf("读取文件目录失败，pathname=%v, err=%v \n", pathname, err)
		return result, err
	}

	// 所有文件/文件夹
	for _, fi := range fis {
		fullname := pathname + "/" + fi.Name()
		// 是文件夹则递归进入获取;是文件，则压入数组
		if fi.IsDir() {
			temp, err := GetAllFile(fullname)
			if err != nil {
				log.Printf("读取文件目录失败,fullname=%v, err=%v", fullname, err)
				return result, err
			}
			result = append(result, temp...)
		} else {
			ext := path.Ext(fullname)
			ext = strings.ToLower(ext)
			if utils.IfIn(ext, extarray) { //判断是否属于特定格式的文件
				result = append(result, fullname)
			}
		}
	}

	return result, nil
}
