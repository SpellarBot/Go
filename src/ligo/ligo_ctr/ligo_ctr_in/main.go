package main

import (
	_ "ligo/ligo_ctr/common/routers"
	"ligo/ligo_ctr/common/controllers"
	"github.com/astaxie/beego"
)

func main() {
	beego.AddAPPStartHook(controllers.GlobalInit)
	go controllers.Consumer()
	beego.Run()
	select{}
}