package controllers

import (
	"fmt"
)

func (w *WebService) SetStatus() {
	statusset := w.GetString("status")
	pass := w.GetString("password")
	if pass == *password {
		switch statusset {
		case "Stop":
			if *status == "Work" {
				*status = "Stop"
				w.Ctx.WriteString(fmt.Sprintf("%s success!",statusset))
			}else{
				w.Ctx.WriteString("You should work it first!")
			}
		case "Restore":
			if *status == "Pause" {
				*status = "Restore"
				w.Ctx.WriteString(fmt.Sprintf("%s success!",statusset))
			}else{
				w.Ctx.WriteString("You should Stop it first!")
			}
		default:
			w.Ctx.WriteString("Unknow status")
		}
	}else{
		w.Ctx.WriteString("Wrong Password!")
	}
}
