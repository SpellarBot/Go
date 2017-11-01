// Elasticsearch Client
// @Author: Golion
// @Date: 2017.3

package elastic

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"vidmate.com/common/http"
	"vidmate.com/common/utils"
)

const (
	esDefaultMaxConn int = 100
	esDefaultTimeout int = 10
)

type ESClient struct {
	Timeout int
	MaxConn int
	Hosts   string
	Auth    string
	Debug   bool
	Logger  func(string)

	httpClients []*http.KeepAliveHttpClient
}

func (e *ESClient) Init() {
	if e.Logger == nil {
		e.Logger = func(msg string) {
			fmt.Printf(msg)
		}
	}
	if e.Timeout <= 0 {
		e.Timeout = esDefaultTimeout
	}
	if e.MaxConn <= 0 {
		e.MaxConn = esDefaultMaxConn
	}
	hosts := strings.Split(e.Hosts, ",")
	e.httpClients = []*http.KeepAliveHttpClient{}
	if len(hosts) <= 0 {
		e.Logger("[ESClient][Init] Error! Please Set ESClient.Hosts Before Init()!")
		os.Exit(0)
	} else {
		for i := 0; i < len(hosts); i++ {
			httpClient := &http.KeepAliveHttpClient{
				Host:    fmt.Sprintf("http://%s", hosts[i]),
				Timeout: e.Timeout,
				MaxConn: e.MaxConn,
				Auth:    e.Auth,
				Debug:   e.Debug,
				Logger:  e.Logger,
			}
			httpClient.Init()
			e.httpClients = append(e.httpClients, httpClient)
		}
	}
}

func (e *ESClient) checkInit() {
	if (e.httpClients == nil) || (len(e.httpClients) <= 0) {
		e.Logger("[ESClient][checkInit] Warning! No Init Before Calling Methods!")
		e.Init()
	}
}

func (e *ESClient) getRandomClient() *http.KeepAliveHttpClient {
	e.checkInit()
	i := rand.Intn(len(e.httpClients))
	return e.httpClients[i]
}

func (e *ESClient) Do(method string, url string, data string) (string, error) {
	e.checkInit()
	switch strings.ToUpper(method) {
	case "GET":
		return e.getRandomClient().Get(url)
	case "POST":
		return e.getRandomClient().Post(url, data)
	case "PUT":
		return e.getRandomClient().Put(url, data)
	case "DELETE":
		return e.getRandomClient().Delete(url)
	default:
		return "", fmt.Errorf("[ESClient][Do] Error! Only Support GET/POST/PUT/DELETE Method!")
	}
}

func (e *ESClient) QPS() int32 {
	e.checkInit()
	var qps int32 = 0
	for i := 0; i < len(e.httpClients); i++ {
		qps += e.httpClients[i].QPS()
	}
	return qps
}

func (e *ESClient) Get(esindex string, estype string, id string) (string, error) {
	e.checkInit()
	url := "/" + esindex + "/" + estype + "/" + id
	resp, err := e.Do("GET", url, "")
	if err != nil {
		return "", fmt.Errorf("[ESClient][Get] error=[%v]", err.Error())
	}
	return resp, nil
}

func (e *ESClient) Scroll(esindex string, estype string, firstQuery string, parseResult func(ESScrollResult), signal *string) (bool, error) {
	e.checkInit()
	firstUrl := "/" + esindex + "/" + estype + "/_search?scroll=1m"
	client := e.getRandomClient()
	resp, err := client.Post(firstUrl, firstQuery)
	if err != nil {
		return false, fmt.Errorf("[ESClient][Scroll] error=[%v]", err.Error())
	}
	var firstScrollResult ESScrollResult
	err = json.Unmarshal([]byte(resp), &firstScrollResult)
	if err != nil {
		return false, fmt.Errorf("[ESClient][Scroll] error=[%v]", err.Error())
	}
	if len(firstScrollResult.ScrollId) > 0 {
		parseResult(firstScrollResult)
		if signal != nil && *signal == "stop" {
			e.Logger("[ESClient][Scroll][afterFirstScroll] Received Stop Signal!")
			return true, nil
		}
		scrollQuery := ESScrollQuery{
			ScrollTime: "1m",
			ScrollId:   firstScrollResult.ScrollId,
		}
		for {
			if signal != nil && *signal == "stop" {
				e.Logger("[ESClient][Scroll] Received Stop Signal!")
				break
			}
			tryCnt := 3
			hit := false
			finished := false
			for tryCnt > 0 {
				tryCnt--
				jsonStr, err := utils.JSONEncode(scrollQuery)
				if err != nil {
					e.Logger(fmt.Sprintf("[ESClient][Scroll] error=[%v]", err.Error()))
					continue
				}
				resp, err := client.Post("/_search/scroll", jsonStr)
				if err != nil {
					e.Logger(fmt.Sprintf("[ESClient][Scroll] error=[%v]", err.Error()))
					continue
				}
				var scrollResult ESScrollResult
				err = json.Unmarshal([]byte(resp), &scrollResult)
				if err != nil {
					e.Logger(fmt.Sprintf("[ESClient][Scroll] error=[%v]", err.Error()))
					continue
				}
				if len(scrollResult.Hits.Hits) <= 0 {
					finished = true
					break
				}
				hit = true
				parseResult(scrollResult)
				break
			}
			if finished {
				break
			}
			if hit {
				continue
			} else {
				return false, fmt.Errorf("[ESClient][Scroll] Error! No Response From ES When Scolling!")
			}
		}
		return true, nil
	} else {
		return false, fmt.Errorf("[ESClient][Scroll] Error! Can Not Get `_scorll_id`!")
	}
}

func (e *ESClient) Put(esindex string, estype string, id string, query string) (bool, error) {
	e.checkInit()
	url := "/" + esindex + "/" + estype + "/" + id
	resp, err := e.Do("PUT", url, query)
	if err != nil {
		return false, fmt.Errorf("[ESClient][Put] error=[%v]", err.Error())
	}
	var putResult esPutResult
	putResult.Init()
	json.Unmarshal([]byte(resp), &putResult)
	if !putResult.CheckPutSucc() {
		return false, fmt.Errorf("[ESClient][Put] Error! Check Valid Failed!")
	}
	return true, nil
}

func (e *ESClient) Put2AllHosts(esindex string, estype string, id string, query string) (bool, error) {
	e.checkInit()
	for i := 0; i < len(e.httpClients); i++ {
		url := "/" + esindex + "/" + estype + "/" + id
		resp, err := e.httpClients[i].Put(url, query)
		if err != nil {
			return false, fmt.Errorf("[ESClient][Put2AllHosts] error=[%v]", err.Error())
		}
		var putResult esPutResult
		putResult.Init()
		json.Unmarshal([]byte(resp), &putResult)
		if !putResult.CheckPutSucc() {
			return false, fmt.Errorf("[ESClient][Put2AllHosts] Error! Check Valid Failed!")
		}
	}
	return true, nil
}
