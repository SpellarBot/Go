// Memcache Client
// @Author: Golion
// @Date: 2017.4

package memcache

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	mc "github.com/bradfitz/gomemcache/memcache"
)

const (
	defaultTimeout     int = 1
	defaultInitConn    int = 10
	defaultMaxConn     int = 1000
	defaultIdleSeconds int = 60
)

type MCClient struct {
	Hosts       string
	Timeout     int
	InitConn    int // 初始化连接数 < MaxConn
	MaxConn     int // 最大连接数
	IdleSeconds int // 空闲连接有效期（秒）
	Debug       bool
	Logger      func(string)

	conn      *mc.Client
	ticker    *time.Ticker
	getQPS    int32
	getQPSCnt int32
	setQPS    int32
	setQPSCnt int32
}

func (m *MCClient) Init() error {
	if m.Logger == nil {
		m.Logger = func(msg string) {
			fmt.Printf(msg)
		}
	}
	if m.Timeout <= 0 {
		m.Timeout = defaultTimeout
	}
	if m.InitConn <= 0 {
		m.InitConn = defaultInitConn
	}
	if m.MaxConn <= 0 {
		m.MaxConn = defaultMaxConn
	}
	if m.IdleSeconds <= 0 {
		m.IdleSeconds = defaultIdleSeconds
	}
	if m.conn == nil {
		hosts := strings.Split(m.Hosts, ",")
		m.conn = mc.New(hosts...)
	}
	if m.ticker == nil {
		go m.countQPS()
	}
	return nil
}

func (m *MCClient) Set(key string, value string, expireSeconds int32) (bool, error) {
	if m.conn == nil {
		return false, fmt.Errorf("[MCClient][Set] Error! Connect to Memcache Failed! Have You Init()?")
	}
	atomic.AddInt32(&m.setQPSCnt, 1)
	if err := m.conn.Set(&mc.Item{
		Key: key,
		Value: []byte(value),
		Expiration: expireSeconds,
	}); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (m *MCClient) Get(key string) (value string, err error) {
	if m.conn == nil {
		return "", fmt.Errorf("[MCClient][Get] Error! Connect to Memcache Failed! Have You Init()?")
	}
	atomic.AddInt32(&m.getQPSCnt, 1)
	if item, err := m.conn.Get(key); err != nil {
		return "", err
	} else {
		return string(item.Value), nil
	}
}

func (m *MCClient) SetQPS() int32 {
	return m.setQPS
}

func (m *MCClient) GetQPS() int32 {
	return m.getQPS
}

func (m *MCClient) countQPS() {
	m.ticker = time.NewTicker(1 * time.Second)
	for _ = range m.ticker.C {
		m.getQPS = m.getQPSCnt
		m.getQPSCnt = 0
		m.setQPS = m.setQPSCnt
		m.setQPSCnt = 0
	}
}
