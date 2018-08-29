package easyserver

import (
	"net"
	"strconv"
)

type EasyTcpClient struct {
	TType  TcpType
	Host   string
	Port   int
	Logger func(string)
	conn   net.Conn
}

func NewEasyTcpClient(ttype TcpType,
	host string,
	port int,
	logger func(string)) (*EasyTcpClient, error) {
	server := EasyTcpClient{
		TType:  ttype,
		Host:   host,
		Port:   port,
		Logger: logger,
	}
	err := server.Init()
	return &server, err
}

func (u *EasyTcpClient) Init() (err error) {
	u.conn, err = net.Dial(string(u.TType), u.Host+":"+strconv.Itoa(u.Port))
	if err != nil {
		u.Logger("TCP Client Conn Fail")
	} else {
		u.Logger("TCP CLient Conn Succ")
	}
	return err
}

func (u *EasyTcpClient) Close() {
	u.conn.Close()
}

func (u *EasyTcpClient) Send(msg []byte) (s []byte, err error) {
	var read int
	get := make([]byte, 4096)
	_, err = u.conn.Write(msg)
	if err == nil {
		u.Logger("Send Msg Succ: " + string(msg))
		read, err = u.conn.Read(get)
		s = get[0:read]
	} else {
		u.Logger("Send Msg Fail: " + err.Error())
	}
	return s, err
}
