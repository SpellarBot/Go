package services
import (
	"strings"
	"encoding/json"
	"fmt"
	"time"
	"os"
	"bufio"
	"errors"
	"strconv"
	"encoding/gob"
)
type CTRUnit struct{
	ShowNum int
	PlayNum int
}

type CTRLine struct {
	CTR []CTRUnit
	StartTime int64
	ExistTime int64
	Step int64
	ExtendTime int64
	Len int
}

// Init
func (c *CTRLine) Init(starttime int64, existtime int64, step int64, extendtime int64){
	c.StartTime = starttime
	c.ExistTime = existtime
	c.Step = step
	c.ExtendTime = extendtime
	c.Len = int((c.ExtendTime+c.ExistTime)/c.Step)
	c.CTR = make([]CTRUnit, c.Len, 2*c.Len)
}


// 检查保存时间长度并更正
func (c *CTRLine) CheckLen(existtime int64){
	if c.ExistTime != existtime{
		if c.ExistTime > existtime{
			N := int((existtime+c.ExtendTime)/c.Step)
			N1 := c.Len - N
			CTR_temp := make([]CTRUnit, N)
			copy(CTR_temp,c.CTR[N1:])
			c.CTR = CTR_temp
			c.StartTime = (c.StartTime+c.ExistTime) - existtime
			c.ExistTime = existtime
			c.Len = int((existtime+c.ExtendTime)/c.Step)
		} else{
			N := int((existtime+c.ExtendTime)/c.Step)
			N1 := N - c.Len
			c.CTR = append(c.CTR,make([]CTRUnit, N1)...)
			c.ExistTime = existtime
			c.Len = int((c.ExistTime+c.ExtendTime)/c.Step)
		}
	}
	return
}

// 检查扩张时间长度并更正
func (c *CTRLine) CheckExtend( extendtime int64){
	if extendtime != c.ExtendTime{
		if extendtime > c.ExtendTime{
			N1 := int((extendtime - c.ExtendTime)/c.Step)
			c.CTR = append(c.CTR,make([]CTRUnit, N1)...)
			c.ExtendTime = extendtime
			c.Len = int((c.ExistTime+c.ExtendTime)/c.Step)
		}else{
			N1 := int((c.ExtendTime- extendtime)/c.Step)
			ctr_temp := c.CTR[N1:]
			c.CTR = ctr_temp
			c.StartTime = c.StartTime+(c.ExtendTime-extendtime)
			c.ExtendTime = extendtime
			c.Len = int((c.ExistTime+c.ExtendTime)/c.Step)
		}
	}
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

//插入一条action(show or play)
func (c *CTRLine) Add(actiontime int64,action string ){
	t := actiontime
	if (t < c.StartTime + c.ExistTime + c.ExtendTime) && (t >= c.StartTime){
		n := int((t - c.StartTime)/c.Step)
		switch action {
		case "ugc_video_show":
			c.CTR[n].ShowNum += 1
		case "video_play_start":
			c.CTR[n].PlayNum += 1
		default:
			return
		}
	} else{
		if t >= (c.ExtendTime+c.ExistTime+c.StartTime){
			k := int((t-c.ExtendTime-c.ExtendTime-c.StartTime)/c.ExtendTime) + 1
			N := int(c.ExtendTime/c.Step)
			BackExtend(c.CTR, k*N)
			c.StartTime += int64(k)*c.ExtendTime
			n := int((t - c.StartTime)/c.Step)
			switch action {
			case "ugc_video_show":
				c.CTR[n].ShowNum += 1
			case "video_play_start":
				c.CTR[n].PlayNum += 1
			default:
				return
			}
		}
	}
}

// 返回总show/play，向两边扩展
func (c *CTRLine) GetCtr(starttime int64, endtime int64) (int, int){
	if (endtime<=starttime) ||
		(starttime>=c.StartTime+c.ExistTime+c.ExistTime) ||
		(endtime<=c.StartTime){
		return 0,0
	}
	var s,e int
	show := 0
	play := 0
	if starttime <= c.StartTime{
		s = 0
	}else{
		s = int((starttime - c.StartTime)/c.Step)
	}
	if endtime >= (c.StartTime+c.ExistTime+c.ExistTime){
		e = c.Len
	}else{
		e = int((endtime-c.StartTime)/c.Step)
	}
	for i:=s;i<e;i++{
		show += c.CTR[i].ShowNum
		play += c.CTR[i].PlayNum
	}
	return show,play
}

//返回show/play趋势
func (c *CTRLine) GetCtrs(starttime int64, endtime int64)([]int, []int){
	var shows,plays []int
	if (endtime<=starttime) ||
		(starttime>=c.StartTime+c.ExistTime+c.ExistTime) ||
		(endtime<=c.StartTime){
		return shows,plays
	}
	for{
		if(starttime>=endtime || starttime >= c.StartTime+c.ExtendTime+c.ExistTime){
			break
		}
		show,play := c.GetCtr(starttime,starttime+c.Step)
		starttime += c.Step
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
	if err == nil{
		encoder := gob.NewEncoder(file)
		err = encoder.Encode(c.VideosCtr)
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
	if IsExist == nil{
		Decoder := gob.NewDecoder(file)
		err := Decoder.Decode(&(c.VideosCtr))
		if err == nil{
			return nil
		}
	}
	return Error("Fail to restore CTRModel")
}

// 保存到本地
func (c *CTRModel) SaveFile(filename string)(error){
	file,err := os.Open(filename)
	defer file.Close()
	if err==nil{
		os.Remove(filename)
	}
	file,err = os.Create(filename)
	V := c.VideosCtr
	Writer := bufio.NewWriter(file)
	for video_id,id_ctr := range V{
		for zipper,zipper_ctr := range id_ctr{
			S := ctr2str(video_id,zipper, zipper_ctr)
			fmt.Println(S)
			Writer.WriteString(S+"\n")
		}
	}
	Writer.Flush()
	return nil
}

// 从本地恢复
func (c *CTRModel) RestoreFile(filename string)(error){
	file,err := os.Open(filename)
	if err != nil{
		return Error("Fail to save CTRModel")
	}
	defer file.Close()
	Reader := bufio.NewReader(file)
	for{
		S,er := Reader.ReadString('\n')
		if er != nil{
			break
		}
		ctrline,video_id,zipper := str2ctr(S)
		_,err1 := c.VideosCtr[video_id]
		if err1 == false{
			c.VideosCtr[video_id] = make(map[string]*CTRLine)
		}
		c.VideosCtr[video_id][zipper] = &ctrline
	}
	return nil
}

// 检查CTRLine与CTRModel的参数是否保持一致并更正
func (c *CTRModel) Check(video_id string,zipper string){
	c.VideosCtr[video_id][zipper].CheckLen(c.ExistTime)
	c.VideosCtr[video_id][zipper].CheckExtend(c.ExtendTime)
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
			c.VideosCtr[video_id][zipper].Add(log_time,action)
		}else{
			c.VideosCtr[video_id][zipper] = &CTRLine{}
			c.VideosCtr[video_id][zipper].Init(log_time, c.ExistTime, c.Step,c.ExtendTime)
			c.VideosCtr[video_id][zipper].Add(log_time, action)
		}
	}else{
		c.VideosCtr[video_id] = make(map[string]*CTRLine)
		c.VideosCtr[video_id][zipper] = &CTRLine{}
		c.VideosCtr[video_id][zipper].Init(log_time,c.ExistTime,c.Step,c.ExtendTime)
		c.VideosCtr[video_id][zipper].Add(log_time,action)
	}
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

func (c *CTRModel) GetZipper(videoId string) ([]string) {
	_,err1 := c.VideosCtr[videoId]
	var zippers []string
	if err1{
		_,err2 := c.VideosCtr[videoId]
		if err2{
			for zipper,_ := range c.VideosCtr[videoId]{
				zippers = append(zippers,zipper)
			}
		}
	}
	zippers = append(zippers, "all")
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

func (c *CTRModel) GetTime(videoId string) (int64, int64){
	_,err1 := c.VideosCtr[videoId]
	var start,end int64
	if err1{
		_,err2 := c.VideosCtr[videoId]
		if err2{
			for _,ctr := range c.VideosCtr[videoId]{
				s := ctr.StartTime
				e := time.Now().Unix()
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


//获取CTRS
func (c *CTRModel) GetCTRs(videoId string, zipper string,starttime int64, endtime int64) ([]int,[]int, int64, int64) {
	_,err1 := c.VideosCtr[videoId]
	var start, step int64
	step = c.Step
	if starttime<(time.Now().Unix()-c.ExistTime){
		start = time.Now().Unix()-c.ExistTime
	}else{
		start = starttime
	}
	var Show,Play []int
	if err1{
		_,err2 := c.VideosCtr[videoId][zipper]
		if err2{
			Show,Play = c.VideosCtr[videoId][zipper].GetCtrs(start, endtime)
		}
		if zipper == "all"{
			k := 0
			for _,value := range c.VideosCtr[videoId]{
				shows,plays := value.GetCtrs(start, endtime)
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
	endtime := time.Now().Unix() - c.ExtendTime - c.ExistTime
	for id,ctr_video:=range c.VideosCtr{
		max_endtime = 0
		for zipper,ctr_zipper:=range ctr_video{
			real_endtime := ctr_zipper.StartTime+ctr_zipper.ExtendTime+ctr_zipper.ExistTime
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

// str 2 timestamp
func Str2timestamp(s string) (int64){
	loc, _ := time.LoadLocation("Asia/Chongqing")
	tm, _ := time.ParseInLocation("2006-01-02 15:04:05", s,loc)
	t := tm.Unix()
	return t
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


func str2ctr(S string)(ctrLine CTRLine, id string, zipper string){
	var StartTime,ExistTime,Step,ExtendTime int64
	var Len int
	SS := strings.Split(S,":")
	S0 := SS[0]
	S1 := SS[1]
	fmt.Sscanf(S0, "%s%s%d%d%d%d%d",
		&id, &zipper, &StartTime, &Step, &ExistTime, &ExtendTime, &Len)
	S10 := strings.Split(S1,",")
	L := make([]CTRUnit,Len,Len*2)
	for i,s10 := range S10{
		show,play := 0,0
		fmt.Sscanf(s10, "%d;%d",&show, &play)
		L[i] = CTRUnit{ShowNum:show,PlayNum:play}
	}
	ctrLine.CTR = L
	ctrLine.ExtendTime = ExtendTime
	ctrLine.ExistTime = ExistTime
	ctrLine.Step = Step
	ctrLine.Len = Len
	ctrLine.StartTime = StartTime
	return
}

func ctr2str(id string, zipper string,ctrline *CTRLine)(S string){
	S1 := fmt.Sprintf("%s %s %d %d %d %d %d",
		id, zipper, ctrline.StartTime, ctrline.Step, ctrline.ExistTime, ctrline.ExtendTime, ctrline.Len)
	S2 := ""
	for _,j := range ctrline.CTR{
		S2 += fmt.Sprintf("%d;%d,",j.ShowNum,j.PlayNum)
	}
	S2 = S2[:len(S2)-1]
	S = fmt.Sprintf("%s:%s",S1,S2)
	return
}
