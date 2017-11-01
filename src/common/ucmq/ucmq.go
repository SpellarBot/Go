// UCMQ Client
// @Author: Golion
// @Date: 2017.3

package ucmq

import (
	"strings"
	"fmt"

	"vidmate.com/common/http"
)

const (
	defaultMaxConn int = 100
	defaultTimeout int = 10
)

type UCMQClient struct {
	Host       string
	TaskName   string
	Timeout    int
	MaxConn    int
	Debug      bool
	Logger     func(string)
	httpClient *http.HttpClient
}

func (u *UCMQClient) Init() {
	if u.Logger == nil {
		u.Logger = func(msg string) {
			fmt.Printf(msg)
		}
	}
	if u.Timeout <= 0 {
		u.Timeout = defaultTimeout
	}
	if u.MaxConn <= 0 {
		u.MaxConn = defaultMaxConn
	}
	u.httpClient = &http.HttpClient{
		Timeout: u.Timeout,
		MaxConn: u.MaxConn,
		Debug:   u.Debug,
		Logger:  u.Logger,
	}
	u.httpClient.Init()
}

func (u *UCMQClient) checkInit() {
	if u.httpClient == nil {
		u.Logger("[UCMQClient][checkInit] Warning! No Init Before Calling Methods!")
		u.Init()
	}
}

func (u *UCMQClient) Put(data string) bool {
	u.checkInit()
	resp, err := u.httpClient.Post("http://" + u.Host + "/?name=" + u.TaskName + "&opt=put&ver=2", data)
	if err != nil {
		u.Logger(fmt.Sprintf("[UCMQClient][Put][FirstTry] error=[%v]", err.Error()))
		return false
	}
	if strings.Contains(resp, "UCMQ_HTTP_OK") {
		return true
	} else {
		// 再试一次
		resp, err = u.httpClient.Post("http://" + u.Host + "/?name=" + u.TaskName + "&opt=put&ver=2", data)
		if err != nil {
			u.Logger(fmt.Sprintf("[UCMQClient][Put][SecondTry] error=[%v]", err.Error()))
			return false
		}
		if strings.Contains(resp, "UCMQ_HTTP_OK") {
			return true
		}
		// 重试失败
		return false
	}
}

func (u *UCMQClient) Get() string {
	u.checkInit()
	resp, err := u.httpClient.Get("http://" + u.Host + "/?name=" + u.TaskName + "&opt=get&ver=2")
	if err != nil {
		u.Logger(fmt.Sprintf("[UCMQClient][Get] error=[%v]", err.Error()))
		return ""
	}
	var _resp string = ""
	if strings.Contains(resp, "UCMQ_HTTP_OK") {
		_resp = resp[14:]
	}
	return _resp
}

func (u *UCMQClient) QPS() int32 {
	u.checkInit()
	return u.httpClient.QPS()
}