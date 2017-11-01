// @Author: Bottle.lwt
// @Date: 2017.8

package core

import (
	"encoding/json"
)

///////////////////////////////////////////////////
// Video Submmit Return structure

type VideoSubmitRet struct {
	Code      int      `json:"code"`
	Msg       string   `json:"msg"`
	RequestId string   `json:"requestId"`
	Data      VideoSubmitRetData `json:"data"`
}

func (result *VideoSubmitRet) ToJSONStr() string {
	output, err := json.Marshal(result)
	if err != nil {
		return ""
	}
	return string(output)
}

type VideoSubmitRetData struct {
	TaskId string `json:"taskId"`
}

///////////////////////////////////////////////////
// Video Green Task OSS Result

type VideoPornOSSResult struct {
	Code      int       `json:"code"`
	Msg       string    `json:"msg"`
	RequestId string    `json:"requestId"`
	Data      VideoPornData `json:"data"`
}

func (result *VideoPornOSSResult) ToJSONStr() string {
	output, err := json.Marshal(result)
	if err != nil {
		return ""
	}
	return string(output)
}

type VideoPornData struct {
	TaskId  string `json:"taskId"`
	Results []RetVideoPornResult `json:"results"`
}

type RetVideoPornResult struct {
	Scene      string           `json:"scene"`
	Label      string           `json:"label"`
	Rate       float64          `json:"rate"`
	Suggestion string           `json:"suggestion"`
	Details    []VideoFrameData   `json:"details"`
}

type VideoFrameData struct {
	Offset int     `json:"offset"`
	Rate   float64 `json:"rate"`
}

type VideoPornResult struct {
	Code       int     `json:"code"`
	TaskId     string  `json:"task_id"`
	Ctime      int64   `json:"ctime"`
	Scene      string  `json:"scene"`
	Label      string  `json:"label"`
	Rate       float64 `json:"rate"`
	Suggestion string  `json:"suggestion"`
}

func (result *VideoPornResult) Init() {
	result.TaskId = ""
	result.Ctime = int64(0)
	result.Scene = "porn"
	result.Label = "reload"
	result.Rate = float64(-1)
	result.Suggestion = "reload"
}

func (result *VideoPornResult) ToJSONStr() string {
	output, err := json.Marshal(result)
	if err != nil {
		return ""
	}
	return string(output)
}

///////////////////////////////////////////////////
// 绿网视频鉴黄接口

// 同步鉴黄接口
/* Result Example:

Result结构体：
字段	类型	是否必须	备注
scene	字符串	必须	风险场景
label	字符串	必须	图片的可能分类。取值范围为:[“normal”,  “porn”]
rate	浮点数	必须	判定为该分类（label）的概率。值越大，越趋向于该分类。取值范围为[0.0, 100.0]
suggestion	字符串	必须	绿网提供的操作建议，取值范围为[“pass”, “review”, “block”]

{
  "code": 200,            //错误码，和HTTP的status code一致
  "msg": "OK",	          //错误码的文本描述
  "requestId": "ali4YGIkY4rr9K54s@9qT8BuJ-1mpFNV",       //该请求的requestId
  "data": {              //当错误发生时，该字段可能不存在。
    "taskId": "ali2YuMPmiAQpB7l36gzz7OLq-1mpFNV"      //该视频检测的任务ID
    "results": [           //检测结果。当错误发生时，该字段可能不存在
      {
        "scene": "porn",   //风险场景
        "label": "normal",    //视频的可能分类。取值范围为:[“normal”, “porn”]。该取值范围随着算法的优化，可能会增加
        "rate": 99.9,     //判定为该分类（label）的概率。值越大，越趋向于该分类。取值范围为[0.0, 100.0]
        "suggestion": "pass"       //绿网提供的操作建议，取值范围为[“pass”, “review”, “block”]
        "details": [         //详情，判断为该分类的某些截帧信息
        	{
        		"offset": 10                //截帧的offset，距离片头的时间，单位秒
        		"rate":                     //判断为该分类的概率
        	}
        ]
      }
    ]
  }
}
*/


