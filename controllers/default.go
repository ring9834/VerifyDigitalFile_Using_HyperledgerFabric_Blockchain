package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	// c.Data["Website"] = "beego.me"
	// c.Data["Email"] = "astaxie@gmail.com"
	// c.TplName = "index.tpl"
	//c.TplName = "index.html"
	//c.TplName = "pays.html"
	//c.TplName = "extracthash.html"
	//c.TplName = "connArchChain.html"
	c.TplName = "firstpage.html"
	//c.TplName = "PickFilePath.html"
}
