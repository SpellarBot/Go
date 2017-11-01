package services

import (
	"strings"
	"encoding/json"
	"fmt"
	"time"
	"os"
	"errors"
	"strconv"
	"encoding/gob"
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
		if actiontime >= c.Start.Last.TimeStamp + c.ExistTime {
			exist,step,extend := c.ExistTime, c.Step, c.ExtendTime
			c.Init(pointstart, exist,step,extend)
		}else{
			N := c.Start
			L := c.Start.Last
			for {
				if (N.TimeStamp+ c.ExistTime > actiontime){
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
		}
		if S.TimeStamp+c.Step > endtime{
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
			if (S.TimeStamp == c.Start.Last.TimeStamp){
				break
			}
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

// CTR计算内核的数据结构
// video_id:zipper:ctrline
type CTRModel struct {
	ExistTime int64
	Step int64
	ExtendTime int64
	Country string
	VideosCtr map[string]map[string]*CTRLine
}

// Init
func (c *CTRModel) Init(existtime int64, step int64, extendtime int64, country string){
	c.ExistTime = existtime
	c.Step = step
	c.ExtendTime = extendtime
	c.Country = country
	c.VideosCtr = make(map[string]map[string]*CTRLine)
}

// 更改Existtime
func (c *CTRModel) ChangeExisttime(existtime int64){
	c.ExistTime = existtime
}
// 更改Extendtime
func (c *CTRModel) ChangeExtendtime(extendtime int64){
	c.ExtendTime = extendtime
}


// 保存到本地
func (c *CTRModel) Save(filename string)(error){
	_,IsExist := os.Stat(filename)
	if IsExist == nil {
		os.Remove(filename)
	}
	file,err := os.Create(filename)
	defer file.Close()
	C := make(map[string]map[string]CTRLineList)
	for id,id_ctr := range c.VideosCtr{
		C[id] = make(map[string]CTRLineList)
		for zipper,zipper_ctr := range id_ctr{
			C0 := CTRLineList{}
			Node2List(zipper_ctr,&C0)
			C[id][zipper] = C0
		}
	}
	if err == nil{
		encoder := gob.NewEncoder(file)
		err = encoder.Encode(C)
		if err == nil{
			return nil
		}
	}
	return Error("Fail to save CTRModel")
}

// 从本地恢复
func (c *CTRModel) Restore(filename string)(error){
	file,IsExist := os.Open(filename)
	defer file.Close()
	C := make(map[string]map[string]CTRLineList)
	if IsExist == nil{
		Decoder := gob.NewDecoder(file)
		err := Decoder.Decode(&C)
		if err == nil{
			for id,id_ctr := range C{
				c.VideosCtr[id] = make(map[string]*CTRLine)
				for zipper,zipper_ctr := range id_ctr{
					C0 := &(CTRLine{})
					List2Node(C0,&zipper_ctr)
					c.VideosCtr[id][zipper] = C0
				}
			}
			return nil
		}
	}
	return Error("Fail to restore CTRModel")
}



// 检查CTRLine与CTRModel的参数是否保持一致并更正
func (c *CTRModel) Check(video_id string,zipper string){
	c.VideosCtr[video_id][zipper].CheckLen(c.ExistTime,c.ExtendTime)
}

//添加一条show/play
//video 和 zipper 都存在时直接添加，不存在时先初始化
//同一个video的不同zipper时间保持一致
func (c *CTRModel) AddUnit(video_id string, zipper string, log_time int64,action string){
	_,err1 := c.VideosCtr[video_id]
	if err1 {
		_,err2 := c.VideosCtr[video_id][zipper]
		if err2 {
			c.Check(video_id, zipper)
		}else{
			c.VideosCtr[video_id][zipper] = &CTRLine{}
			c.VideosCtr[video_id][zipper].Init(log_time, c.ExistTime, c.Step,c.ExtendTime)
		}
	}else{
		c.VideosCtr[video_id] = make(map[string]*CTRLine)
		c.VideosCtr[video_id][zipper] = &CTRLine{}
		c.VideosCtr[video_id][zipper].Init(log_time,c.ExistTime,c.Step,c.ExtendTime)
	}
	c.VideosCtr[video_id][zipper].Add(log_time, action)
}


// 添加一条记录
func (c *CTRModel) AddLog(S string){
	if strings.Contains(S, fmt.Sprintf("country=%s",c.Country)) {
		if strings.Contains(S, "ugc_video_show") {
			video_ids, zippers,log_time,app_ver := GetVideosZippers(S)
			N := len(video_ids)
			if N == len(zippers) && log_time >0 && app_ver >1.40 {
				for i:=0;i<N;i++{
					id := video_ids[i]
					zi := zippers[i]
					if id != ""{
						if zi == ""{
							zi = "others"
						}
						c.AddUnit(id,zi,log_time,"ugc_video_show")
					}
				}
			}

		}else{
			if strings.Contains(S, "video_play_start") {
				video_id, zipper,log_time,app_ver := GetVideoZipper(S)
				if video_id != "" && app_ver > 1.40 {
					if zipper == ""{
						zipper = "others"
					}
					c.AddUnit(video_id,zipper,log_time,"video_play_start")
				}
			}
		}
	}
}

//返回time区间
func (c *CTRModel) GetTime(videoId string) (int64, int64){
	_,err1 := c.VideosCtr[videoId]
	var start,end int64
	if err1{
		_,err2 := c.VideosCtr[videoId]
		if err2{
			for _,ctr := range c.VideosCtr[videoId]{
				s := ctr.Start.TimeStamp
				e := ctr.Start.Last.TimeStamp
				if s < start{
					start = s
				}
				if e > end{
					end = e
				}
			}
		}
	}
	return start, end
}

//返回zipper
func (c *CTRModel) GetZipper(videoId string) ([]string) {
	_,err1 := c.VideosCtr[videoId]
	var zippers []string
	zippers = append(zippers, "all")
	if err1{
		_,err2 := c.VideosCtr[videoId]
		if err2{
			for zipper,_ := range c.VideosCtr[videoId]{
				zippers = append(zippers,zipper)
			}
		}
	}
	return zippers
}

//获取CTR
func (c *CTRModel) GetCTR(videoId string, zipper string,starttime int64, endtime int64) (int,int) {
	_,err1 := c.VideosCtr[videoId]
	Show,Play := 0,0
	if err1{
		_,err2 := c.VideosCtr[videoId][zipper]
		if err2{
			Show,Play = c.VideosCtr[videoId][zipper].GetCtr(starttime, endtime)
		}
		if zipper == "all"{
			for _,value := range c.VideosCtr[videoId]{
				show,play := value.GetCtr(starttime, endtime)
				Show += show
				Play += play
			}
		}
		return Show,Play
	}
	return 0,0
}


//获取CTRS
func (c *CTRModel) GetCTRs(videoId string, zipper string,starttime int64, endtime int64) ([]int,[]int, int64, int64) {
	_,err1 := c.VideosCtr[videoId]
	var start, step,end int64
	step = c.Step
	if starttime<(time.Now().Unix()-c.ExistTime){
		start = time.Now().Unix()-c.ExistTime
	}else{
		start = starttime
	}
	if endtime>=time.Now().Unix(){
		end = time.Now().Unix()
	}else{
		end = endtime
	}
	var Show,Play []int
	if err1{
		_,err2 := c.VideosCtr[videoId][zipper]
		if err2{
			Show,Play = c.VideosCtr[videoId][zipper].GetCtrs(start, end)
		}
		if zipper == "all"{
			k := 0
			for _,value := range c.VideosCtr[videoId]{
				shows,plays := value.GetCtrs(start, end)
				if k == 0{
					Show = shows
					Play = plays
				}else{
					for i,_ := range shows{
						Show[i] += shows[i]
						Play[i] += plays[i]
					}
				}
				k ++
			}
		}
		return Show,Play,start,step
	}
	return Show,Play,start,step
}

//去除太久时间之前的video ，结束时间太早则予以删除
func (c *CTRModel)Delete(){
	var max_endtime int64
	endtime := time.Now().Unix() - c.ExistTime - c.ExtendTime
	for id,ctr_video:=range c.VideosCtr{
		max_endtime = 0
		for zipper,ctr_zipper:=range ctr_video{
			real_endtime := ctr_zipper.Start.Last.TimeStamp
			if real_endtime>max_endtime{
				max_endtime = real_endtime
			}
			if real_endtime<endtime{
				delete(c.VideosCtr[id],zipper)
			}
		}
		if max_endtime<endtime{
			delete(c.VideosCtr,id)
		}
	}
}

// 数组后移n个位置
func BackExtend(s1 []CTRUnit,n int){
	L := len(s1)
	if n<L && n>0{
		copy(s1,s1[n:])
		N := len(s1)-n
		s2 := make([]CTRUnit,n,n)
		copy(s1[N:],s2)
	}else{
		s2 := make([]CTRUnit,L,L)
		copy(s1,s2)
	}
}



// 解析Log_content
type LogContents struct{
	Zippers string `json:"zippers"`
	Video_ids string `json:"video_ids"`
}
type LogContent struct{
	Zipper string `json:"zipper"`
	Video_id string `json:"video_id"`
}
type KafkaMessage struct{
	Message string `json:"message"`
}

// ugc_video_show 解析
func GetVideosZippers( kafkamessage string)([]string,[]string, int64, float64){
	var K KafkaMessage
	var video_ids,zippers []string
	var log_time int64
	var L LogContents
	var app_ver float64
	json.Unmarshal([]byte(kafkamessage), &K)
	if K.Message == ""{
		return video_ids,zippers,log_time,app_ver
	}
	S := K.Message
	A := strings.Split(S,"`")
	for _,i := range A{
		if strings.Contains(i,"video_ids"){
			i = i[12:]
			json.Unmarshal([]byte(i), &L)
			Video_ids := L.Video_ids
			Zippers := L.Zippers
			video_ids = strings.Split(Video_ids,",")
			zippers = strings.Split(Zippers,",")
			if len(video_ids) >0 && len(zippers)<len(video_ids){
				for t:=0;t<(len(video_ids)-len(zippers));t++{
					zippers = append(zippers,"others")
				}
			}
		}
		if strings.Contains(i,"log_time="){
			log_time = Str2timestamp(i[9:])
		}
		if strings.Contains(i, "app_ver="){
			_,err := strconv.ParseFloat(i[8:], 32)
			if err != nil{
				app_ver = 0.0
			}else{
				app_ver,_ = strconv.ParseFloat(i[8:], 32)
			}
		}
	}
	return video_ids,zippers,log_time,app_ver
}

// video_play_start 解析
func GetVideoZipper( kafkamessage string)(string,string,int64, float64){
	var video_id,zipper string
	var log_time int64
	var app_ver float64
	var K KafkaMessage
	json.Unmarshal([]byte(kafkamessage), &K)
	if K.Message == ""{
		return video_id,zipper,log_time,app_ver
	}
	A := strings.Split(K.Message,"`")
	var L LogContent
	for _,i := range A{
		if strings.Contains(i,"video_id"){
			i = i[12:]
			json.Unmarshal([]byte(i), &L)
			video_id = L.Video_id
			zipper = L.Zipper
		}
		if strings.Contains(i,"log_time="){
			log_time = Str2timestamp(i[9:])
		}
		if strings.Contains(i, "app_ver="){
			_,err := strconv.ParseFloat(i[8:], 32)
			if err != nil{
				app_ver = 0.0
			}else{
				app_ver,_ = strconv.ParseFloat(i[8:], 32)
			}
		}
	}
	if video_id != "" && zipper == ""{
		zipper = "others"
	}
	return video_id,zipper,log_time,app_ver
}


func Error(Msg string) error {
	var e error = errors.New(Msg)
	return e
}


func Node2List (C1 *CTRLine, C2 *CTRLineList){
	C2.ExistTime = C1.ExistTime
	C2.ExtendTime = C1.ExtendTime
	C2.Step = C1.Step
	S := C1.Start
	end := C1.Start.Last.TimeStamp
	for{
		C2.Start = append(C2.Start,CTRUnitL{ShowNum:S.ShowNum,PlayNum:S.PlayNum,TimeStamp:S.TimeStamp})
		if S.TimeStamp == end{
			break
		}
		S = S.Next
	}
}

func List2Node (C1 *CTRLine, C2 *CTRLineList){
	C1.ExistTime = C2.ExistTime
	C1.ExtendTime = C2.ExtendTime
	C1.Step = C2.Step
	var S *CTRUnit
	N := len(C2.Start)
	for i:=0;i<N;i++{
		value := C2.Start[i]
		if i != 0{
			S.Next = &(CTRUnit{TimeStamp:value.TimeStamp,ShowNum:value.ShowNum,PlayNum:value.PlayNum})
			S.Next.Last = S
			S = S.Next
		}else {
			C1.Start = &(CTRUnit{TimeStamp: value.TimeStamp, ShowNum: value.ShowNum, PlayNum: value.PlayNum})
			S = C1.Start
		}
	}
	C1.Start.Last = S
	S.Next = C1.Start
	return
}

type CTRLineList struct {
	ExistTime int64
	Step int64
	ExtendTime int64
	Start []CTRUnitL
}

type CTRUnitL struct {
	TimeStamp int64
	ShowNum int
	PlayNum int
}