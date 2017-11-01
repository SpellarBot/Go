package controllers

import (
	"encoding/json"
	"time"
	"fmt"
)

type Result struct {
	Show int`json:"show"`
	Play int `json:"play"`
	Ctr float64 `json:"ctr"`
	Status string `json:"status"`
}


// 输出show play
func (w *WebService) GetCTR() {
	t1 := time.Now().UnixNano()
	zipper   := w.GetString("zipper")
	videoId  := w.GetString("video_id")
	start, err1 := w.GetInt64("start")
	end, err2   := w.GetInt64("end")
	var C Result
	C.Status = "error"
	if zipper == "" || videoId =="" || err1 !=nil || err2 != nil{
		show,play := 0,0
		C = Result{Show:show,Play:play,Ctr:0.0}
	}else{
		show,play := ctrModel.GetCTR(videoId, zipper, start, end)
		C = Result{Show:show,Play:play,Ctr:CalCtr(show,play),Status:"normal"}
	}
	J,err := json.Marshal(C)
	if err == nil{
		w.Ctx.WriteString(string(J))
	}
	t2 := time.Now().UnixNano()
	t := float64(t2-t1)/1e6
	logClient.Info(fmt.Sprintf("Get ctr cost %f ms",t))
}