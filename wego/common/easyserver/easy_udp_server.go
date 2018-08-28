package easyserver

import (
	"fmt"
	"net"
	"runtime"
)

type UdpType string

const (
	Udp  = UdpType("udp")
	Udp4 = UdpType("udp4")
	Udp6 = UdpType("udp6")
)

type EasyUdpServer struct {
	UType       UdpType
	Port        int
	Threads     int
	Responser   func(string) string
	Logger      func(string)
	WriteBuffer int
	ReadBuffer  int
	listener    *net.UDPConn
}

func NewEasyUdpServer(utype UdpType,
	port int,
	threads int,
	writebuffer int,
	readbuffer int,
	responser func(string) string,
	logger func(string)) (*EasyUdpServer, error) {
	server := EasyUdpServer{
		UType:       utype,
		Port:        port,
		Threads:     threads,
		Logger:      logger,
		Responser:   responser,
		WriteBuffer: writebuffer,
		ReadBuffer:  readbuffer,
	}
	err := server.Init()
	return &server, err
}

func (u *EasyUdpServer) Init() error {
	if u.Logger == nil {
		u.Logger = func(s string) {
			fmt.Println(s)
		}
	}
	if u.Responser == nil {
		u.Responser = func(s string) string {
			return "OK"
		}
	}
	if u.Port < 0 {
		u.Port = 8080
	}
	if u.Threads <= 0 {
		u.Threads = runtime.NumCPU()
	}
	if u.WriteBuffer < 64 {
		u.WriteBuffer = 64
	}
	if u.ReadBuffer < 64 {
		u.ReadBuffer = 64
	}
	var err error
	u.listener, err = getUdpListener(string(u.UType), u.Port)
	if err == nil {
		u.listener.SetReadBuffer(u.ReadBuffer)
		u.listener.SetWriteBuffer(u.WriteBuffer)
		u.Logger(fmt.Sprintf("Serve Start Succ At Port %d", u.Port))
		u.Logger(fmt.Sprintf("Serve Start %d Threads", u.Threads))
		for i := 0; i < u.Threads; i++ {
			go u.listen()
		}
	}

	return err
}

func (u *EasyUdpServer) listen() {
	for {
		readdata, remoteAddr, err := u.readFromUdp()
		if err == nil {
			senddata := u.Responser(string(readdata))
			u.writeToUdp(remoteAddr, senddata)
		}
	}
}

func (u *EasyUdpServer) readFromUdp() ([]byte, *net.UDPAddr, error) {
	var msg string
	readdata := make([]byte, u.ReadBuffer)
	read, remoteAddr, err := u.listener.ReadFromUDP(readdata)
	if err == nil {
		msg = fmt.Sprintf("Read %d From %s Succ: %s", read, udpaddr2str(remoteAddr), string(readdata[0:read]))
		u.Logger(msg)
		return readdata[0:read], remoteAddr, nil
	}
	msg = fmt.Sprintf("Read 0 From %s Fail:", udpaddr2str(remoteAddr))
	u.Logger(msg)
	return readdata, remoteAddr, err
}

func (u *EasyUdpServer) writeToUdp(remoteAddr *net.UDPAddr, send string) {
	var msg string
	data := []byte(send)
	write, err := u.listener.WriteToUDP(data, remoteAddr)
	if err == nil {
		msg = fmt.Sprintf("Write %d To %s SUcc: %s", write, udpaddr2str(remoteAddr), string(data[0:write]))
	} else {
		msg = fmt.Sprintf("Write 0 To %s Fail:", udpaddr2str(remoteAddr))
	}
	u.Logger(msg)
}

func (u *EasyUdpServer) Close() {
	u.listener.Close()
}

func getUdpListener(proto string, port int) (*net.UDPConn, error) {
	listener, err := net.ListenUDP(proto, &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
	})
	return listener, err
}

func udpaddr2str(add *net.UDPAddr) (s string) {
	return fmt.Sprintf("%v:%v", add.IP, add.Port)
}
