// @Author: Golion
// @Date: 2017.2

package core

import (
	"fmt"
	"encoding/json"
	"time"
)

///////////////////////////////////////////////////
// Porn Task OSS Result

type PornOSSResult struct {
	Code      int       `json:"code"`
	Msg       string    `json:"msg"`
	RequestId string    `json:"requestId"`
	Data      PornTasks `json:"data"`
}

func (result *PornOSSResult) ToJSONStr() string {
	output, err := json.Marshal(result)
	if err != nil {
		return ""
	}
	return string(output)
}

type PornTasks struct {
	Results []RetPornResult `json:"results"`
}

type RetPornResult struct {
	Scene      string  `json:"scene"`
	Label      string  `json:"label"`
	Rate       float64 `json:"rate"`
	Suggestion string  `json:"suggestion"`
}

type PornResult struct {
	TaskId     string  `json:"task_id"`
	Ctime      int64   `json:"ctime"`
	Scene      string  `json:"scene"`
	Label      string  `json:"label"`
	Rate       float64 `json:"rate"`
	Suggestion string  `json:"suggestion"`
}

type PornScannerResult struct {
	Status int             `json:"status"`
	Msg    string          `json:"msg"`
	Result PornResult      `json:"result"`
}

func (result *PornResult) Init() {
	result.TaskId = ""
	result.Ctime = int64(0)
	result.Scene = "porn"
	result.Label = "reload"
	result.Rate = float64(-1)
	result.Suggestion = "reload"
}

func (result *PornResult) ToJSONStr() string {
	output, err := json.Marshal(result)
	if err != nil {
		return ""
	}
	return string(output)
}

///////////////////////////////////////////////////
// 绿网鉴黄接口

// 同步鉴黄接口
/* Result Example:

Result结构体：
字段	类型	是否必须	备注
scene	字符串	必须	风险场景
label	字符串	必须	图片的可能分类。取值范围为:[“normal”,  “porn”]
rate	浮点数	必须	判定为该分类（label）的概率。值越大，越趋向于该分类。取值范围为[0.0, 100.0]
suggestion	字符串	必须	绿网提供的操作建议，取值范围为[“pass”, “review”, “block”]

{
  "code": 200,
  "msg": "OK",
  "requestId": "ali4YGIkY4rr9K54s@9qT8BuJ-1mpFNV",
  "data": {
    "taskId": "ali2YuMPmiAQpB7l36gzz7OLq-1mpFNV"
    "results": [
      {
        "scene": "porn",
        "label": "normal",
        "rate": 99.9,
        "suggestion": "pass"
      }
    ]
  }
}

 */
func GetPornScore(bucketName string, imageName string) (PornResult, error) {
	startTime := time.Now().Unix()
	ret := GetOSSObjectWithParams(bucketName, imageName, "x-oss-process=udf/green/image/scan,porn")
	var result PornOSSResult
	json.Unmarshal(ret, &result)
	fmt.Printf("Get OSS Result - %v\n", result.ToJSONStr())
	var pornResult PornResult
	if result.Msg != "OK" {
		pornResult.Init()
		return pornResult, fmt.Errorf(result.Msg)
	} else if len(result.Data.Results) > 0 {
		pornResult.Ctime      = startTime
		pornResult.TaskId     = result.RequestId
		pornResult.Scene      = result.Data.Results[0].Scene
		pornResult.Label      = result.Data.Results[0].Label
		pornResult.Rate       = result.Data.Results[0].Rate
		pornResult.Suggestion = result.Data.Results[0].Suggestion
		return pornResult, nil
	} else {
		pornResult.Init()
		return pornResult, fmt.Errorf("No Result")
	}
}