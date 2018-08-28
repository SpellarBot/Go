package easyserver

import (
	"net"
	"strconv"
)

type EasyUdpClient struct {
	UType    UdpType
	Host     string
	Port     int
	Logger   func(string)
	listener net.Conn
}

func NewEasyUdpClient(utype UdpType,
	host string,
	port int,
	logger func(string)) (*EasyUdpClient, error) {
	server := EasyUdpClient{
		UType:  utype,
		Host:   host,
		Port:   port,
		Logger: logger,
	}
	err := server.Init()
	return &server, err
}

func (u *EasyUdpClient) Init() (err error) {
	u.listener, err = net.Dial(string(u.UType), u.Host+":"+strconv.Itoa(u.Port))
	if err != nil {
		u.Logger("UDP Client Conn Fail")
	} else {
		u.Logger("UDP CLient Conn Succ")
	}
	return err
}

func (u *EasyUdpClient) Close() {
	u.listener.Close()
}

func (u *EasyUdpClient) Send(msg string) (s string, err error) {
	var read int
	get := make([]byte, 4096)
	_, err = u.listener.Write([]byte(msg))
	if err == nil {
		u.Logger("Send Msg Succ: " + msg)
		read, err = u.listener.Read(get)
		s = string(get[0:read])
	} else {
		u.Logger("Send Msg Fail: " + msg)
	}
	return s, err
}
