// Transwarp Client
// @Author: Golion
// @Date: 2017.5

package transwarp

import (
	"fmt"
	"strings"
	"encoding/json"

	"vidmate.com/common/http"
	"vidmate.com/common/utils"
)

const (
	defaultMaxConn int = 10
	defaultTimeout int = 10
)

type TranswarpClient struct {
	Host       string
	Auth       string
	Timeout    int
	MaxConn    int
	Debug      bool // 为true时将会开启所有日志（包括底层http的日志）
	SetLog     bool // 默认是false，为true时输出SetValue和SetConf的普通日志
	GetLog     bool // 默认是false，为true时输出GetValue和GetConf的普通日志
	Logger     func(string)

	httpClient *http.HttpClient
}

type Conf struct {
	Id         string `json:"id"`
	Step       string `json:"step"`
	Status     string `json:"status"`
	Comment    string `json:"comment"`
	UpdateTime int    `json:"update_time"`
}

func (t *TranswarpClient) Init() {
	if t.Logger == nil {
		t.Logger = func(msg string) {
			fmt.Printf(msg)
		}
	}
	if t.Timeout <= 0 {
		t.Timeout = defaultTimeout
	}
	if t.MaxConn <= 0 {
		t.MaxConn = defaultMaxConn
	}
	t.httpClient = &http.HttpClient{
		Timeout: t.Timeout,
		MaxConn: t.MaxConn,
		Debug:   t.Debug,
		Logger:  t.Logger,
	}
	t.httpClient.Init()
}

func (t *TranswarpClient) checkInit() {
	if t.httpClient == nil {
		t.Logger("[TranswarpClient][checkInit] Warning! No Init Before Calling Methods!")
		t.Init()
	}
}

func (t *TranswarpClient) CheckSert(id string, by string) (bool, error) {
	t.checkInit()
	url := t.Host + "/checksert" +
		"?id=" + utils.URLEncode(id) +
		"&by=" + utils.URLEncode(by)
	if resp, err := t.httpClient.Get(url); err != nil {
		t.Logger(fmt.Sprintf("[TranswarpClient][CheckSert] Error! id=[%v] error=[%v]",
			id, err.Error()))
		return false, err
	} else {
		if strings.Contains(strings.ToLower(resp), "already exist") {
			if t.Debug {
				t.Logger(fmt.Sprintf("[TranswarpClient][CheckSert] Already Exist! id=[%v]", id))
			}
			return true, nil
		} else if strings.Contains(strings.ToLower(resp), "insert succ") {
			if t.Debug || t.SetLog {
				t.Logger(fmt.Sprintf("[TranswarpClient][CheckSert] Insert Succeed! id=[%v]", id))
			}
			return true, nil
		} else {
			errorMsg := fmt.Sprintf("[TranswarpClient][CheckSert] Error! id=[%v] resp=[%v]", id, resp)
			t.Logger(errorMsg)
			return false, fmt.Errorf(errorMsg)
		}
	}
}

func (t *TranswarpClient) SetValue(id string, col string, value string, by string) (bool, error) {
	t.checkInit()
	url := t.Host + "/setvalue" +
		"?id=" + utils.URLEncode(id) +
		"&col=" + utils.URLEncode(col) +
		"&value=" + utils.URLEncode(value) +
		"&by=" + utils.URLEncode(by)
	if resp, err := t.httpClient.Get(url); err != nil {
		t.Logger(fmt.Sprintf("[TranswarpClient][SetValue] Error! id=[%v] col=[%v] value=[%v] error=[%v]",
			id, col, value, err.Error()))
		return false, err
	} else {
		if strings.Contains(strings.ToLower(resp), "update succ") {
			if t.Debug || t.SetLog {
				t.Logger(fmt.Sprintf("[TranswarpClient][SetValue] Succeed! id=[%v] col=[%v] value=[%v]",
					id, col, value))
			}
			return true, nil
		} else {
			errorMsg := fmt.Sprintf("[TranswarpClient][SetValue] Error! Transwarp Do Not Return `update succ`! id=[%v] col=[%v] value=[%v] resp=[%v]",
				id, col, value, resp)
			t.Logger(errorMsg)
			return false, fmt.Errorf(errorMsg)
		}
	}
}

func (t *TranswarpClient) SetConf(id string, step string, status string, comment string, by string) (bool, error) {
	t.checkInit()
	url := t.Host + "/setconf" +
		"?id=" + utils.URLEncode(id) +
		"&step=" + utils.URLEncode(step) +
		"&status=" + utils.URLEncode(status) +
		"&comment=" + utils.URLEncode(comment) +
		"&by=" + utils.URLEncode(by)
	if resp, err := t.httpClient.Get(url); err != nil {
		t.Logger(fmt.Sprintf("[TranswarpClient][SetConf] Error! id=[%v] step=[%v] status=[%v] comment=[%v] error=[%v]",
			id, step, status, comment, err.Error()))
		return false, err
	} else {
		if strings.Contains(strings.ToLower(resp), "update succ") {
			if t.Debug || t.SetLog {
				t.Logger(fmt.Sprintf("[TranswarpClient][SetConf] Succeed! id=[%v] step=[%v] status=[%v] comment=[%v]",
					id, step, status, comment))
			}
			return true, nil
		} else {
			errorMsg := fmt.Sprintf("[TranswarpClient][SetConf] Error! Transwarp Do Not Return `update succ`! id=[%v] step=[%v] status=[%v] comment=[%v] resp=[%v]",
				id, step, status, comment, resp)
			t.Logger(errorMsg)
			return false, fmt.Errorf(errorMsg)
		}
	}
}

func (t *TranswarpClient) GetValue(id string, col string, by string) (string, error) {
	t.checkInit()
	url := t.Host + "/getvalue" +
		"?id=" + utils.URLEncode(id) +
		"&col=" + utils.URLEncode(col) +
		"&by=" + utils.URLEncode(by)
	if resp, err := t.httpClient.Get(url); err != nil {
		t.Logger(fmt.Sprintf("[TranswarpClient][GetValue] Error! id=[%v] col=[%v] error=[%v]",
			id, col, err.Error()))
		return "", err
	} else {
		if t.Debug || t.GetLog {
			t.Logger(fmt.Sprintf("[TranswarpClient][GetValue] id=[%v] col=[%v] Got:value=[%v]",
				id, col, resp))
		}
		return resp, nil
	}
}

func (t *TranswarpClient) GetConf(id string, by string) (Conf, error) {
	t.checkInit()
	var conf Conf
	url := t.Host + "/getconf" +
		"?id=" + utils.URLEncode(id) +
		"&by=" + utils.URLEncode(by)
	if resp, err := t.httpClient.Get(url); err != nil {
		t.Logger(fmt.Sprintf("[TranswarpClient][GetConf] Error! id=[%v] error=[%v]", id, err.Error()))
		return conf, err
	} else {
		if err := json.Unmarshal([]byte(resp), &conf); err != nil {
			t.Logger(fmt.Sprintf("[TranswarpClient][GetConf] Error! id=[%v] Got:conf=[%v] error=[%v]",
				id, resp, err.Error()))
			return conf, err
		} else {
			if t.Debug || t.GetLog {
				t.Logger(fmt.Sprintf("[TranswarpClient][GetConf] id=[%v] Got:conf=[%v]",
					id, resp))
			}
			return conf, nil
		}
	}
}

func (t *TranswarpClient) QPS() int32 {
	t.checkInit()
	return t.httpClient.QPS()
}