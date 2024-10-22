package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hzx/models"
	"hzx/utils"
	"io/ioutil"
	"path"
	"strings"

	"github.com/astaxie/beego"
)

type FilesInfoController struct {
	beego.Controller
}

var FileInfos = make(map[string]interface{})

type DNameValPair struct {
	Name string `json:"name"`
	Val  string `json:"val"`
}

func (c *FilesInfoController) GetDisks() {

	jsonPath := ""
	if utils.IsWindowRunTime() { //os windows
		jsonPath = "diskconf1.json"
	} else { //os linux
		jsonPath = "diskconf2.json"
	}

	apiUrl := "./static/diskdata/" + jsonPath
	data, err := ioutil.ReadFile(apiUrl)
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf")) //防止出现编码问题：invalid character 'ï' looking for beginning of value
	if err != nil {
		return
	}

	//var str []interface{} //切片
	var str []DNameValPair //切片

	errr := json.Unmarshal(data, &str)
	if errr == nil {
		fmt.Printf("JSON转换有误!")
	}
	c.Data["json"] = str
	c.ServeJSON() //返回JSON对象,供前端使用
}

type pth struct {
	Mypath string `json:"mypth"`
}

// 获取指定目录下的所有文件夹和文件信息
func (c *FilesInfoController) GetAllFolderAndFiles() {
	result := []models.FilePathInfo{} //自定义类FilePathInfo

	acceptedParam := c.Ctx.Input.RequestBody //注意：从前端传过来的JSON格式要用data:'{}',一定要用单引号，.net core中不用单引号
	var pathparam pth
	merr := json.Unmarshal(acceptedParam, &pathparam)
	if merr != nil {
		fmt.Printf("JSON参数解析有误" + merr.Error())
	}

	fis, err := ioutil.ReadDir(pathparam.Mypath)
	if err != nil {
		fmt.Printf("读取文件目录失败!pathname=%v, err=%v \n", pathparam, err)
		//return result, err
		c.Data["json"] = result
		c.ServeJSON() //返回JSON对象,供前端使用
	}

	// 所有文件/文件夹 ；GetExtensionIcon对扩展名进行处理，以适应前端icon的扩展名
	for index, fi := range fis {
		fullname := pathparam.Mypath + "/" + fi.Name()
		fullname = utils.ReplaceSlashStr(fullname)
		fileNameWithSuffix := path.Base(fullname)
		fileType := path.Ext(fileNameWithSuffix)
		fileType = utils.GetExtensionIcon(fileType)                                                                                                              //获取文件的后缀(文件类型)
		fileNameOnly := strings.TrimSuffix(fileNameWithSuffix, fileType)                                                                                         //获取文件名称(不带后缀)
		fpi := models.FilePathInfo{ID: index, Name: fileNameOnly, FullName: fullname, IsDir: fi.IsDir(), CreateTime: fi.ModTime().String(), Extension: fileType} //一个对象赋值;一个字段对应一个值，用冒号隔开，否则会报错--composite literal uses unkeyed fields
		result = append(result, fpi)                                                                                                                             //将对象放入数组
	}

	//return result, nil
	c.Data["json"] = result
	c.ServeJSON() //返回JSON对象,供前端使用
}

func (c *FilesInfoController) GetAllFilesToFront(pathname string) {
	result, _ := GetAllFiles(pathname)
	c.Data["json"] = result
	c.ServeJSON() //返回JSON对象,供前端使用
}

// 递归获取指定目录下的所有文件名
func GetAllFiles(pathname string) ([]string, error) {
	result := []string{}

	fis, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Printf("读取文件目录失败!pathname=%v, err=%v \n", pathname, err)
		return result, err
	}

	// 所有文件/文件夹
	for _, fi := range fis {
		fullname := pathname + "/" + fi.Name()
		// 是文件夹则递归进入获取;是文件，则压入数组
		if fi.IsDir() {
			temp, err := GetAllFiles(fullname) //递归
			if err != nil {
				fmt.Printf("读取文件目录失败,fullname=%v, err=%v", fullname, err)
				//return result, err
				continue
			}
			result = append(result, temp...)
		} else {
			result = append(result, fullname)
		}
	}

	return result, nil
}
