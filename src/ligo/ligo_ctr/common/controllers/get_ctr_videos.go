package controllers


import (
	"strings"
	"encoding/json"
	"fmt"
	"time"
)
type ctrResult struct {
	Id string `json:"id"`
	Show int `json:"show"`
	Play int `json:"play"`
	Ctr float64`json:"ctr"`
}

type ctrResults struct {
	Data []ctrResult `json:"data"`
	Status string `json:"status"`
}

// 输出多个CTR
func (w *WebService) GetCTRVideos() {
	t1 := time.Now().UnixNano()
	videos  := w.GetString("video_ids")
	start, err1 := w.GetInt64("start")
	end, err2   := w.GetInt64("end")
	video_Ids := strings.Split(videos,",")
	var C ctrResults
	C.Data = make([]ctrResult,0)
	C.Status = "normal"
	if !(len(video_Ids) == 0 || err1 !=nil || err2 != nil){
		C.Status = "error"
		for _,value:=range video_Ids{
			show,play := ctrModel.GetCTR(value, "all", start, end)
			c := ctrResult{Id:value,Show:show,Play:play,Ctr:CalCtr(show,play)}
			C.Data = append(C.Data,c)
		}
	}
	J,err := json.Marshal(C)
	if err == nil{
		w.Ctx.WriteString(string(J))
	}
	t2 := time.Now().UnixNano()
	t := float64(t2-t1)/1e6
	logClient.Info(fmt.Sprintf("Get ctr mult cost %f ms",t))

}
