package main
import (
	"fmt"
	"time"
)

func Str2timestamp(s string) (int64){
	loc, _ := time.LoadLocation("Asia/Chongqing")
	tm, _ := time.ParseInLocation("2006-01-02 15:04:05", s,loc)
	t := tm.Unix()
	return t
}
func Timestamp2str(x int64)(string){
	loc, _ := time.LoadLocation("Asia/Chongqing")
	tm := time.Unix(x,0).In(loc)
	return tm.Format("2006-01-02 15:04:05")
}

type CTRUnit struct{
	TimeStamp int64
	ShowNum int
	PlayNum int
	Next *CTRUnit
	Last *CTRUnit
}

type CTRLine struct {
	ExistTime int64
	Step int64
	ExtendTime int64
	Start *CTRUnit
}

// Init
func (c *CTRLine) Init(starttime int64, existtime int64, step int64, extendtime int64){
	c.ExistTime = existtime
	c.Step = step
	c.ExtendTime = extendtime
	c.Start = &(CTRUnit{TimeStamp:starttime})
	c.Start.Next = c.Start
	c.Start.Last = c.Start
}


//扩张和缩减都只能向后，扩张只需要改参数，缩减才需要操作
func (c *CTRLine) BackClip(allexisttime int64){
	L := c.Start.Last
	N := c.Start
	nowexisttime := L.TimeStamp + c.Step - N.TimeStamp
	if nowexisttime > allexisttime{
		shouldstarttime := L.TimeStamp + c.Step - allexisttime
		for{
			if N.TimeStamp >= shouldstarttime{
				break
			}
			N = N.Next
		}
		c.Start = N
		c.Start.Last = L
		L.Next = c.Start
	}
}

//向后平移到可以包括actiontime的位置
func (c *CTRLine) BackMove(actiontime int64){
	pointstart := c.Start.TimeStamp + c.Step*( (actiontime-c.Start.TimeStamp)/c.Step )
	if actiontime >= c.Start.TimeStamp + c.ExistTime + c.ExtendTime{
		if actiontime >= c.Start.Last.TimeStamp + c.ExistTime + c.ExtendTime{
			exist,step,extend := c.ExistTime, c.Step, c.ExtendTime
			c.Init(pointstart, exist,step,extend)
		}else{
			N := c.Start
			L := c.Start.Last
			for {
				if (N.TimeStamp + c.ExistTime < actiontime) &&
					(N.TimeStamp+ c.ExistTime + c.ExtendTime > actiontime){
					break
				}
				N = N.Next
			}
			L.Next = N
			N.Last = L
			c.Start = N
		}
	}
}

// 检查保存时间长度并更正
func (c *CTRLine) CheckLen(existtime int64, extendtime int64){
	allexisttime := existtime + extendtime
	nowexisttime := c.Start.Last.TimeStamp+c.Step - c.Start.TimeStamp
	if nowexisttime > allexisttime{
		c.BackClip(allexisttime)
	}
	c.ExistTime = existtime
	c.ExtendTime = extendtime
	return
}

//TODO:检查步长并更正:不科学（只能改大不能改小）
func (c *CTRLine) CheckStep( step int64){
	if step!= c.Step{
		if step > c.Step{
		}
	}
	return
}


//增加，不超出最长时限
func (c *CTRLine) AddLimited(t int64,action string ){
	N := c.Start.Last
	pointstart := c.Start.TimeStamp + c.Step*( (t-c.Start.TimeStamp)/c.Step )
	for {
		if N.TimeStamp <= pointstart {
			break
		}
		N = N.Last
	}
	if N.TimeStamp == pointstart {
		switch action {
		case "ugc_video_show":
			N.ShowNum ++
		case "video_play_start":
			N.PlayNum ++
		default:
			return
		}
	} else {
		P := CTRUnit{TimeStamp: pointstart}
		P.Next = N.Next
		N.Next = &P
		P.Last = N
		P.Next.Last = &P
		switch action {
		case "ugc_video_show":
			P.ShowNum ++
		case "video_play_start":
			P.PlayNum ++
		default:
			return
		}
	}
	return
}



//插入一条action(show or play)
func (c *CTRLine) Add(actiontime int64,action string ) {
	t := actiontime
	N := c.Start.Last
	if (t < N.TimeStamp+c.Step) && (t >= c.Start.TimeStamp) {
		c.AddLimited(t, action)
	}else {
		if t >= (N.TimeStamp+c.Step) {
			if t >= (c.Start.TimeStamp+c.ExtendTime+c.ExistTime) {
				c.BackMove(t)
			}
			c.AddLimited(t, action)
		}
	}
}

// 返回总show/play，向两边扩展
func (c *CTRLine) GetCtr(starttime int64, endtime int64) (int, int){
	show := 0
	play := 0
	if endtime < c.Start.TimeStamp ||
		starttime > c.Start.Last.TimeStamp + c.Step ||
		starttime > endtime{
		return show,play
	}
	if starttime <= c.Start.TimeStamp{
		starttime = c.Start.TimeStamp
	}
	if endtime >= c.Start.Last.TimeStamp{
		endtime = c.Start.Last.TimeStamp
	}
	S := c.Start
	for {
		if (S.TimeStamp+c.Step >= starttime && S.TimeStamp <= endtime){
			show += S.ShowNum
			play += S.PlayNum
			fmt.Println(show,play,Timestamp2str(S.TimeStamp))
		}
		if S.TimeStamp+c.Step > endtime{
			fmt.Println(Timestamp2str(S.TimeStamp),Timestamp2str(endtime))
			break
		}
		S = S.Next
	}
	return show,play
}

//返回show/play趋势
func (c *CTRLine) GetCtrs(starttime int64, endtime int64)([]int, []int){
	var shows,plays []int
	var show,play int
	if endtime<=starttime{
		return shows,plays
	}
	start := starttime
	end := start + c.Step
	S := c.Start
	for{
		show,play = 0,0
		if(start >= endtime){
			break
		}
		for {
			if (S.TimeStamp+c.Step >= start && S.TimeStamp < end){
				show += S.ShowNum
				play += S.PlayNum
			}
			if S.TimeStamp+c.Step > end{
				break
			}
			S = S.Next
		}
		start += c.Step
		end += c.Step
		shows = append(shows,show)
		plays = append(plays,play)
	}
	return shows,plays
}
func Print(C CTRLine){
	fmt.Println(fmt.Sprintf("Starttime  : %s",Timestamp2str(C.Start.TimeStamp)))
	fmt.Println(fmt.Sprintf("Start      : %s",Timestamp2str(C.Start.TimeStamp)))
	fmt.Println(fmt.Sprintf("End        : %s",Timestamp2str(C.Start.Last.TimeStamp)))
	fmt.Println(fmt.Sprintf("Step       : %d",C.Step))
	fmt.Println(fmt.Sprintf("Existtime  : %d",C.ExistTime))
	fmt.Println(fmt.Sprintf("Extendtime : %d",C.ExtendTime))

	S := C.Start
	L := S.Last
	for {
		fmt.Println(fmt.Sprintf("Time:%20s; Show:%d; Play:%d",Timestamp2str(S.TimeStamp),S.ShowNum,S.PlayNum))
		if S.TimeStamp==L.TimeStamp{
			break
		}
		S = S.Next
	}
}

func main(){
	var C CTRLine
	var start,exist,extend,step int64
	start = Str2timestamp("2017-10-27 15:01:00")
	exist = 3600
	step = 60
	extend = 3600
	C.Init(start,exist,step,extend)
	Print(C)
	C.Add(Str2timestamp("2017-10-27 16:02:30"),"ugc_video_show")
	Print(C)
	C.Add(Str2timestamp("2017-10-27 16:03:44"),"video_play_start")
	C.Add(Str2timestamp("2017-10-27 17:05:30"),"ugc_video_show")
	C.Add(Str2timestamp("2017-10-27 17:09:30"),"ugc_video_show")
	C.Add(Str2timestamp("2017-10-27 16:08:30"),"video_play_start")
	C.Add(Str2timestamp("2017-10-27 17:11:30"),"video_play_start")
	C.Add(Str2timestamp("2017-10-27 16:39:30"),"video_play_start")
	Print(C)
	//C.CheckLen(1800,1800)
	Print(C)
	fmt.Println(C.GetCtrs(Str2timestamp("2017-10-27 16:00:00"),Str2timestamp("2017-10-27 16:08:13")))
}