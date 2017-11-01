package controllers

import (
	"strconv"
	"time"
	"fmt"
)
// 输出show play
func (w *WebService) GetDetail() {
	var t1,t2 int64
	t1 = time.Now().UnixNano()
	w.Ctx.WriteString("================Current Videos======================"+"\n")
	w.Ctx.WriteString("ExistTime  : "+strconv.Itoa(int(ctrModel.ExistTime))+"\n")
	w.Ctx.WriteString("Step       : "+strconv.Itoa(int(ctrModel.Step))+"\n")
	w.Ctx.WriteString("ExtendTime : "+strconv.Itoa(int(ctrModel.ExtendTime))+"\n")
	N := 0
	w.Ctx.WriteString("====================================================\n")
	for id,ctr_video:=range ctrModel.VideosCtr{
		N ++
		for zipper,ctr_zipper := range ctr_video{
			w.Ctx.WriteString(fmt.Sprintf("  ID: %15s",id))
			show,play := ctrModel.GetCTR(id,zipper,0,time.Now().Unix()+60*60*48)
			w.Ctx.WriteString(fmt.Sprintf("    Zipper: %10s; StartTime: %20s; EndTime: %20s; Show: %5d; Play: %5d",
				zipper,timestamp2str(ctr_zipper.Start.TimeStamp),timestamp2str(ctr_zipper.Start.Last.TimeStamp),show,play))
			w.Ctx.WriteString("\n")
		}
		w.Ctx.WriteString("------------------------------------------------------------------------------\n")
	}
	t2 = time.Now().UnixNano()
	t := float64(t2 - t1)/1e6
	w.Ctx.WriteString(fmt.Sprintf("All %d Videos, Cost %f ms",N,t))
}
func timestamp2str(x int64)(string){
	loc, _ := time.LoadLocation("Asia/Chongqing")
	tm := time.Unix(x,0).In(loc)
	return tm.Format("2006-01-02 15:04:05")
}

func nowtime()(string){
	loc, _ := time.LoadLocation("Asia/Chongqing")
	tm := time.Unix(time.Now().Unix(),0).In(loc)
	return tm.Format("15:04")
}
