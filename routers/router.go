package routers

import (
	"hzx/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/query", &controllers.MainController{}, "Post:PayQuery") // 新增PayQuery路由
	beego.Router("/add", &controllers.MainController{}, "Post:PayAdd")     // 新增PayAdd路由
	beego.Router("/paysview", &controllers.MainController{}, "Post:GetPaysView")
	//beego.Router("/pays", &controllers.MainController{}, "Post:GetPays")
	beego.Router("/pays", &controllers.MainController{}, "Post:GetPagedPays")
	//beego.AutoRouter(&controllers.MainController{})

	beego.Router("/GetDsks", &controllers.FilesInfoController{}, "Post:GetDisks")
	beego.Router("/GetDirectories", &controllers.FilesInfoController{}, "Post:GetAllFolderAndFiles")

	beego.Router("/ServerEndUpHash", &controllers.FabricHashController{}, "Get:ExtracthashSrvView") //页面用Get
	beego.Router("/ClientEndUpHash", &controllers.FabricHashController{}, "Get:ExtracthashClntView")
	beego.Router("/VerifyReality", &controllers.FabricHashController{}, "Get:VerifyRealityView")
	beego.Router("/extractHash", &controllers.FabricHashController{}, "Post:CreateArchAssets") //从客户端把一个或多个档案原文上区块链

	beego.Router("/createAssetsFromSvr", &controllers.FabricSvrHashController{}, "Post:CreateAssetsFromServer")
	beego.Router("/verifyarch", &controllers.FabricVerifyController{}, "Post:VerifyArchReality")

	beego.Router("/fbbrowse", &controllers.FabricBrowserController{}, "Get:ChainBrowserView")
	beego.Router("/getPagedAssets", &controllers.FabricBrowserController{}, "Post:ReadPagedAssets")
}
