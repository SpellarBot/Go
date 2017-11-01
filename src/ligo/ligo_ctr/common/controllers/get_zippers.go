package controllers

import (
	"encoding/json"
	"time"
	"fmt"
)

type Zippers struct{
	Data []string `json:"data"`
	Status string `json:"status"`
}

func (w *WebService) GetZippers() {
	t1 := time.Now().UnixNano()
	videoId  := w.GetString("video_id")
	var Z Zippers
	Z.Status = "normal"
	if videoId == ""{
		Z.Status = "error"
	}else{
		Z.Data = ctrModel.GetZipper(videoId)
	}
	J,err := json.Marshal(Z)
	if err == nil{
		w.Ctx.WriteString(string(J))
	}
	t2 := time.Now().UnixNano()
	t := float64(t2 - t1)/1e6
	logClient.Info(fmt.Sprintf("Get ctrs cost %f ms",t))
}