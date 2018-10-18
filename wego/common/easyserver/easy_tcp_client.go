package easyserver

import (
	"net"
	"strconv"
	"sync"
	"time"
)

type EasyTcpClient struct {
	TType        TcpType
	Host         string
	Port         int
	Timeout      time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	ReadBuffer   int
	Logger       func(string)

	lock sync.RWMutex
	conn net.Conn
}

func (t *EasyTcpClient) Init() (err error) {
	t.conn, err = net.Dial(string(t.TType), t.Host+":"+strconv.Itoa(t.Port))
	if err != nil {
		t.Logger("TCP Client Conn Fail")
	} else {
		t.Logger("TCP CLient Conn Succ")
	}
	t.lock = sync.RWMutex{}

	if t.Timeout > 0 {
		t.conn.SetDeadline(time.Now().Add(t.Timeout))
	}
	if t.WriteTimeout > 0 {
		t.conn.SetReadDeadline(time.Now().Add(t.ReadTimeout))
	}
	if t.ReadTimeout > 0 {
		t.conn.SetWriteDeadline(time.Now().Add(t.WriteTimeout))
	}

	return err
}

func (u *EasyTcpClient) Close() {
	u.conn.Close()
}

func (u *EasyTcpClient) Send(msg []byte) (s []byte, err error) {
	u.lock.Lock()
	defer u.lock.Unlock()
	var read int
	get := make([]byte, u.ReadBuffer)
	_, err = u.conn.Write(msg)
	if err == nil {
		u.Logger("Send Msg Succ: " + string(msg))
		read, err = u.conn.Read(get)
		if err == nil {
			s = get[0:read]
			u.Logger("Get Msg Succ: " + string(s))
		} else {
			u.Logger("Get Msg Fail: " + err.Error())
		}
	} else {
		u.Logger("Send Msg Fail: " + err.Error())
	}
	return s, err
}
