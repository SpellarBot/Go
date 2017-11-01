// KeepAlive HTTP Client
// @Author: Golion
// @Date: 2017.3

package http

import "fmt"

type KeepAliveHttpClient struct {
	Timeout int
	MaxConn int
	Auth    string
	Host    string
	Debug   bool
	Logger  func(string)

	conn    *HttpClient
}

func (h *KeepAliveHttpClient) Init() {
	if h.Logger == nil {
		h.Logger = func(msg string) {
			fmt.Printf(msg)
		}
	}
	h.conn = &HttpClient{
		Timeout: h.Timeout,
		MaxConn: h.MaxConn,
		Auth:    h.Auth,
		Debug:   h.Debug,
		Logger:  h.Logger,
	}
	h.conn.Init()
}

func (h *KeepAliveHttpClient) checkInit() {
	if h.conn == nil {
		h.Logger("[KeepAliveHttpClient][checkInit] Warning! No Init Before Calling Methods!")
		h.Init()
	}
}

func (h *KeepAliveHttpClient) Get(url string) (string, error) {
	h.checkInit()
	return h.conn.Get(h.Host + url)
}

func (h *KeepAliveHttpClient) Post(url string, data string) (string, error) {
	h.checkInit()
	return h.conn.Post(h.Host + url, data)
}

func (h *KeepAliveHttpClient) Put(url string, data string) (string, error) {
	h.checkInit()
	return h.conn.Put(h.Host + url, data)
}

func (h *KeepAliveHttpClient) Delete(url string) (string, error) {
	h.checkInit()
	return h.conn.Delete(h.Host + url)
}

func (h *KeepAliveHttpClient) QPS() int32 {
	h.checkInit()
	return h.conn.QPS()
}