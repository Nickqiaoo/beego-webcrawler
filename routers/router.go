package routers

import (
	"beego-webcrawler/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/login", &controllers.MainController{},"get:Login;post:Craw")
	beego.Router("/checkcode", &controllers.MainController{},"get:Checkcode")
	beego.Router("/grade", &controllers.MainController{},"post:Querygrade")
	beego.Router("/credit", &controllers.MainController{},"post:Querycredit")
}
