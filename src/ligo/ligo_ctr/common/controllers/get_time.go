package controllers

import (
	"time"
	"encoding/json"
	"fmt"
)

type TimeResult struct {
	Starttime int64`json:"starttime"`
	Endtime int64 `json:"endtime"`
	Status string `json:"status"`
}


// 输出show play
func (w *WebService) GetTime() {
	t1 := time.Now().UnixNano()
	videoId  := w.GetString("video_id")
	var Tr TimeResult
	Tr.Status = "error"
	if videoId !="" {
		start,end := ctrModel.GetTime(videoId)
		Tr.Status = "normal"
		Tr.Starttime = start
		Tr.Endtime = end
	}
	J,err := json.Marshal(Tr)
	if err == nil{
		w.Ctx.WriteString(string(J))
	}
	t2 := time.Now().UnixNano()
	t := float64(t2-t1)/1e6
	logClient.Info(fmt.Sprintf("Get time cost %f ms",t))
}
