package controllers

import (
	"fmt"
)

func (w *WebService) SetTime() {
	extendtime,err1 := w.GetInt64("extendtime")
	existtime,err2 := w.GetInt64("existtime")
	pass := w.GetString("password")
	if pass == *password && err1 == nil && err2 == nil && extendtime>=3600 && existtime >=3600{
		ctrModel.ChangeExisttime(existtime)
		ctrModel.ChangeExtendtime(extendtime)
		w.Ctx.WriteString("Time Parameters Changed!")
		logClient.Info(fmt.Sprintf("ExtendTime changed to %d",ctrModel.ExtendTime))
		logClient.Info(fmt.Sprintf("ExistTime changed to %d",ctrModel.ExistTime))
	}else{
		w.Ctx.WriteString("Wrong Password Or Wrong Time Parameter!")
	}
}
