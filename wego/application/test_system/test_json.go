package main
import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)
// 解析Log_content
type LogContents struct{
	Zippers string `json:"zippers"`
	Video_ids string `json:"video_ids"`
}
type LogContent struct{
	Zipper string `json:"zipper"`
	Video_id string `json:"video_id"`
}

// ugc_video_show 解析
func GetVideosZippers( S string)([]string,[]string, string, float64){
	A := strings.Split(S,"`")
	var video_ids,zippers []string
	var log_time string
	var L LogContents
	var app_ver float64
	for _,i := range A{
		if strings.Contains(i,"video_ids"){
			i = i[12:]
			json.Unmarshal([]byte(i), &L)
			Video_ids := L.Video_ids
			Zippers := L.Zippers
			video_ids = strings.Split(Video_ids,",")
			zippers = strings.Split(Zippers,",")
		}
		if strings.Contains(i,"log_time="){
			log_time = i[9:]
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
func GetVideoZipper( S string)(string,string,string, float64){
	A := strings.Split(S,"`")
	var video_id,zipper,log_time string
	var app_ver float64
	video_id = ""
	zipper = ""
	var L LogContent
	for _,i := range A{
		if strings.Contains(i,"video_id"){
			i = i[12:]
			json.Unmarshal([]byte(i), &L)
			video_id = L.Video_id
			zipper = L.Zipper
		}
		if strings.Contains(i,"log_time="){
			log_time = i[9:]
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
	return video_id,zipper,log_time,app_ver
}



func main(){
	S := "uid=10371973`utdid=WXti+HCEBPgDADrs9fcKI1TO`network_type=HSPA`action_code=ugc_video_show`app_id=com.uc.vmate.app.gp`log_content={\"sr\":\"0\",\"refer\":\"\",\"scene\":\"UGCVideoDiscover_New\",\"abtag\":\"abtag-3\",\"zippers\":\"new,new,new\",\"video_ids\":\"a4ll1d1s8s0,ahbheag9680,aww8l66dpau\"}`client_id=72834995-1e4d-4b01-af5f-0300fe14bdc9`log_id=46afbeb6-9f13-4ba3-b8a6-edd8b8140afa_3`app_ver=1.43`country=BD`log_time=2017-10-24 17:48:14`ip=202.134.9.130"
	if strings.Contains(S, "ugc_video_show") {
		fmt.Println(GetVideosZippers(S))

		}else{
			if strings.Contains(S, "video_play_start") {
				fmt.Println(GetVideoZipper(S))
			}
		}
}