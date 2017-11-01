// HTTP Client
// @Author: Golion
// @Date: 2017.3

package http

/**
Example:

import "common/http"

h := http.HttpClient{
	Timeout: 5,
	MaxConn: 200,
	Auth: "HTTP_BASE_AUTH_USER:HTTP_BASE_AUTH_PASS",
}
h.Init()
resp, err := h.Get("http://www.google.com/")

*/

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

const (
	defaultMaxConn int = 100
	defaultTimeout int = 3
)

type HttpClient struct {
	Timeout int
	MaxConn int
	Headers string
	Auth    string
	Debug   bool
	Logger  func(string)

	httpHeaders      [][]string
	httpBaseAuthUser string
	httpBaseAuthPass string
	ticker           *time.Ticker
	conn             *http.Client
	qps              int32
	qpsCnt           int32
}

func (h *HttpClient) Init() {
	if h.Logger == nil {
		h.Logger = func(msg string) {
			fmt.Printf(msg)
		}
	}
	if h.Timeout <= 0 {
		h.Timeout = defaultTimeout
	}
	if h.MaxConn <= 0 {
		h.MaxConn = defaultMaxConn
	}
	if len(h.Auth) > 0 {
		auth := strings.Split(h.Auth, ":")
		if len(auth) == 2 {
			h.httpBaseAuthUser = auth[0]
			h.httpBaseAuthPass = auth[1]
		}
	}
	if len(h.Headers) > 0 {
		h.httpHeaders = [][]string{}
		rawHeaders := strings.Split(h.Headers, "|")
		for i := 0; i < len(rawHeaders); i++ {
			rawHeader := strings.Split(rawHeaders[i], "^")
			if len(rawHeader) == 2 {
				header := []string{rawHeader[0], rawHeader[1]}
				h.httpHeaders = append(h.httpHeaders, header)
			}
		}
	}
	h.conn = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: h.MaxConn,
		},
		Timeout: time.Duration(h.Timeout) * time.Second,
	}
	if h.ticker == nil {
		go h.countQPS()
	}
}

func (h *HttpClient) Do(method string, url string, data string) (string, error) {
	if h.conn == nil || h.Logger == nil {
		h.Logger("[HttpClient][Do] Warning! No Init Before Calling Methods!")
		h.Init()
	}
	atomic.AddInt32(&h.qpsCnt, 1)
	t0 := time.Now()
	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		h.printDebugLog(err, method, url, t0)
		return "", fmt.Errorf("[HttpClient][Do][NewRequest] error=[%v]", err.Error())
	} else {
		if len(h.httpBaseAuthUser) > 0 && len(h.httpBaseAuthPass) > 0 {
			req.SetBasicAuth(h.httpBaseAuthUser, h.httpBaseAuthPass)
		}
		if len(h.Headers) > 0 {
			for i := 0; i < len(h.httpHeaders); i++ {
				req.Header.Add(h.httpHeaders[i][0], h.httpHeaders[i][1])
			}
		}
		resp, err := h.conn.Do(req)
		if err != nil && resp == nil && strings.Contains(strings.ToLower(err.Error()), "connection reset by peer") {
			resp, err = h.conn.Do(req)
		}
		if err != nil && resp == nil {
			h.printDebugLog(err, method, url, t0)
			return "", fmt.Errorf("[HttpClient][Do][Sending2Endpoint] error=[%v]", err.Error())
		} else {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				h.printDebugLog(err, method, url, t0)
				return "", fmt.Errorf("[HttpClient][Do][ParsingBody] error=[%v]", err.Error())
			}
			h.printDebugLog(err, method, url, t0)
			return string(body), nil
		}
	}
}

func (h *HttpClient) Get(url string) (string, error) {
	return h.Do("GET", url, "")
}

func (h *HttpClient) Post(url string, data string) (string, error) {
	return h.Do("POST", url, data)
}

func (h *HttpClient) Put(url string, data string) (string, error) {
	return h.Do("PUT", url, data)
}

func (h *HttpClient) Delete(url string) (string, error) {
	return h.Do("DELETE", url, "")
}

func (h *HttpClient) QPS() int32 {
	return h.qps
}

func (h *HttpClient) countQPS() {
	h.ticker = time.NewTicker(1 * time.Second)
	for _ = range h.ticker.C {
		h.qps = h.qpsCnt
		h.qpsCnt = 0
	}
}

func (h *HttpClient) printDebugLog(err error, method string, url string, t0 time.Time) {
	if h.Debug {
		t1 := time.Now()
		if err == nil {
			h.Logger(fmt.Sprintf("[HttpClient][Debug] qps=[%v] method=[%s] url=[%s] costtime=[%v] succ",
				h.qps, method, url, t1.Sub(t0).Seconds()*1000))
		} else {
			h.Logger(fmt.Sprintf("[HttpClient][Debug] qps=[%v] method=[%s] url=[%s] costtime=[%v] error=[%s]",
				h.qps, method, url, t1.Sub(t0).Seconds()*1000, err.Error()))
		}
	}
}
