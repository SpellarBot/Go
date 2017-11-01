// DingTalk Robot Client
// @Author: Golion
// @Date: 2017.5

package utils

import (
	"fmt"

	"vidmate.com/common/http"
)

type DingTalkRobotClient struct {
	RobotHost string
	RobotName string
	Logger    func(string)

	httpClient *http.HttpClient
}

func (d *DingTalkRobotClient) Init() {
	if d.Logger == nil {
		d.Logger = func(msg string) {
			fmt.Printf(msg)
		}
	}
	d.httpClient = &http.HttpClient{
		Headers: "Content-Type^application/json",
		Timeout: 10,
		MaxConn: 1,
		Logger:  d.Logger,
	}
	d.httpClient.Init()
}

func (d *DingTalkRobotClient) Send(msg string) {
	if d.httpClient == nil {
		d.Logger("[DingdingRobotClient][Send] Error! No Init Before Calling Methods!")
		return
	}
	if d.RobotHost == "" {
		d.Logger("[DingdingRobotClient][Send] Error! Please Set RobotHost!")
		return
	}
	jsonStr := []byte(`
	{
		"msgtype": "text",
		"text": {
			"content": "` + msg + `"
		}
	}`)
	if jsonObj, err := JSONDecode(string(jsonStr)); err == nil {
		if output, err := JSONEncode(jsonObj); err == nil {
			if resp, err := d.httpClient.Post(d.RobotHost, output); err == nil {
				d.Logger(fmt.Sprintf("[DingdingRobotClient][Send][Post] Succeed! name=[%v] resp=[%v]", d.RobotName, resp))
			} else {
				d.Logger(fmt.Sprintf("[DingdingRobotClient][Send][Post] Error! name=[%v] error=[%v]", d.RobotName, err.Error()))
			}
		} else {
			d.Logger(fmt.Sprintf("[DingdingRobotClient][Send][JSONEncode] Error! name=[%v] error=[%v]", d.RobotName, err.Error()))
		}
	} else {
		d.Logger(fmt.Sprintf("[DingdingRobotClient][Send][JSONDecode] Error! name=[%v] error=[%v]", d.RobotName, err.Error()))
	}
}