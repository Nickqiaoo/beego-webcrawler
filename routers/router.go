package routers

import (
	"beego-webcrawler/controllers"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/evaluate", &controllers.MainController{}, "post:Evaluate")
	beego.Router("/login", &controllers.MainController{}, "get:Login;post:Craw")
	beego.Router("/checkcode", &controllers.MainController{}, "get:CheckCode")
	beego.Router("/grade", &controllers.MainController{}, "post:QueryGrade")
	beego.Router("/credit", &controllers.MainController{}, "post:QueryCredit")
	beego.Handler("/metrics", promhttp.Handler())

}
