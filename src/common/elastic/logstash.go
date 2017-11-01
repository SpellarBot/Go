// Logstash Client
// @Author: Golion
// @Date: 2017.3

package elastic

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"vidmate.com/common/http"
	"vidmate.com/common/utils"
)

const (
	logstashDefaultMaxConn int = 100
	logstashDefaultTimeout int = 1
)

type logstashPutData struct {
	Type  string `json:"_type"`
	Id    string `json:"_id"`
	Query string `json:"_query"`
}

type LogstashClient struct {
	Timeout int
	MaxConn int
	Hosts   string
	Auth    string
	Debug   bool
	Logger  func(string)

	httpClients []*http.KeepAliveHttpClient
}

func (l *LogstashClient) Init() {
	if l.Logger == nil {
		l.Logger = func(msg string) {
			fmt.Printf(msg)
		}
	}
	if l.Timeout <= 0 {
		l.Timeout = logstashDefaultTimeout
	}
	if l.MaxConn <= 0 {
		l.MaxConn = logstashDefaultMaxConn
	}
	hosts := strings.Split(l.Hosts, ",")
	l.httpClients = []*http.KeepAliveHttpClient{}
	if len(hosts) <= 0 {
		l.Logger("[LogstashClient][Init] Error! Please Set LogstashClient.Hosts Before Init()!")
		os.Exit(0)
	} else {
		for i := 0; i < len(hosts); i++ {
			httpClient := &http.KeepAliveHttpClient{
				Host:    fmt.Sprintf("http://%s", hosts[i]),
				Timeout: l.Timeout,
				MaxConn: l.MaxConn,
				Auth:    l.Auth,
				Debug:   l.Debug,
				Logger:  l.Logger,
			}
			httpClient.Init()
			l.httpClients = append(l.httpClients, httpClient)
		}
	}
}

func (l *LogstashClient) checkInit() {
	if (l.httpClients == nil) || (len(l.httpClients) <= 0) {
		l.Logger("[LogstashClient][checkInit] Warning! No Init Before Calling Methods!")
		l.Init()
	}
}

func (l *LogstashClient) getRandomClient() *http.KeepAliveHttpClient {
	i := rand.Intn(len(l.httpClients))
	return l.httpClients[i]
}

func (l *LogstashClient) Do(method string, url string, data string) (string, error) {
	l.checkInit()
	switch strings.ToUpper(method) {
	case "GET":
		return l.getRandomClient().Get(url)
	case "POST":
		return l.getRandomClient().Post(url, data)
	case "PUT":
		return l.getRandomClient().Put(url, data)
	case "DELETE":
		return l.getRandomClient().Delete(url)
	default:
		return "", fmt.Errorf("[LogstashClient][Do] Error! Only Support GET/POST/PUT/DELETE Method!")
	}
}

func (l *LogstashClient) QPS() int32 {
	l.checkInit()
	var qps int32 = 0
	for i := 0; i < len(l.httpClients); i++ {
		qps += l.httpClients[i].QPS()
	}
	return qps
}

func (l *LogstashClient) Put(estype string, esid string, query string) (bool, error) {
	l.checkInit()
	data := logstashPutData{
		Type:  estype,
		Id:    esid,
		Query: query,
	}
	jsonStr, err := utils.JSONEncode(data)
	if err != nil {
		return false, fmt.Errorf("[LogstashClient][Put] error=[%v]", err.Error())
	}
	resp, err := l.Do("POST", "", jsonStr)
	if err != nil {
		return false, fmt.Errorf("[LogstashClient][Put] error=[%v]", err.Error())
	}
	if strings.Contains(strings.ToLower(resp), "ok") {
		return true, nil
	} else {
		return false, fmt.Errorf("[LogstashClient][Put] Error! Logstash DO NOT Return `OK`!")
	}
}
