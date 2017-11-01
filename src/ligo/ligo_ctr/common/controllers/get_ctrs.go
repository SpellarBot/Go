package controllers

import (
	"encoding/json"
	"time"
	"fmt"
)

type CtrResult struct {
	Show int `json:"show"`
	Play int `json:"play"`
	Ctr float64 `json:"ctr"`
	Timestamp int64 `json:"timestamp"`
}

type CtrResults struct {
	Id string `json:"id"`
	Zipper string `json:"zipper"`
	Data []CtrResult `json:"data"`
	Status string `json:"status"`
}


// 输出show play
func (w *WebService) GetCTRs() {
	t1 := time.Now().UnixNano()
	videoId  := w.GetString("video_id")
	zipper   := w.GetString("zipper")
	start, err1 := w.GetInt64("start")
	end, err2   := w.GetInt64("end")
	var C2 CtrResults
	var C0 []CtrResult
	C2.Status = "normal"
	C2.Id = videoId
	C2.Zipper = zipper
	C2.Data = make([]CtrResult,0)
	var shows,plays []int
	var starttime,step int64
	if videoId =="" || err1 !=nil || err2 != nil{
		C2.Status = "error"
	}else{
		shows,plays,starttime,step = ctrModel.GetCTRs(videoId, zipper,start, end)
		C0 = make([]CtrResult,0)
		for i,j := range shows{
			C0 = append(C0, CtrResult{Show:j,Play:plays[i],Timestamp:starttime+int64(i)*step,Ctr:CalCtr(j,plays[i])})
		}
		C2.Data = C0
	}
	J,err := json.Marshal(C2)
	if err == nil{
		w.Ctx.WriteString(string(J))
	}
	t2 := time.Now().UnixNano()
	t := float64(t2 - t1)/1e6
	logClient.Info(fmt.Sprintf("Get ctrs cost %f ms",t))
}