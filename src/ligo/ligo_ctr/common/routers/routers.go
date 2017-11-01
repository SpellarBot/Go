package routers

import (
	"ligo/ligo_ctr/common/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/set_time", &controllers.WebService{}, "*:SetTime")
	beego.Router("/set_status",&controllers.WebService{},"*:SetStatus")
	beego.Router("/get_detail", &controllers.WebService{}, "*:GetDetail")
	beego.Router("/get_zippers", &controllers.WebService{}, "*:GetZippers")
	beego.Router("/get_time", &controllers.WebService{}, "*:GetTime")
	beego.Router("/get_ctr", &controllers.WebService{}, "*:GetCTR")
	beego.Router("/get_ctrs", &controllers.WebService{}, "*:GetCTRs")
	beego.Router("/get_ctr_videos", &controllers.WebService{}, "*:GetCTRVideos")
}
