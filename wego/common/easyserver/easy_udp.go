package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"time"
)

type UdpType string

const (
	Udp  = UdpType("udp")
	Udp4 = UdpType("udp4")
	Udp6 = UdpType("udp6")
)

type EasyUdpServer struct {
	UType    UdpType
	Port     int
	Threads  int
	Logger   func(string)
	listener *net.UDPConn
}

func NewEasyUdpServer(utype UdpType, port int, threads int, logger func(string)) *EasyUdpServer {
	return &EasyUdpServer{
		UType:   utype,
		Port:    port,
		Threads: threads,
		Logger:  logger,
	}
}

func (u *EasyUdpServer) Init() error {
	if u.Logger == nil {
		u.Logger = func(s string) {
			fmt.Println(s)
		}
	}
	if u.Port == 0 {
		u.Port = 8080
	}
	if u.Threads <= 0 {
		u.Threads = runtime.NumCPU()
	}
	var err error
	u.listener, err = getUdpListener(string(u.UType), u.Port)
	if err == nil {
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
		u.readFromUdp()
	}
}

func (u *EasyUdpServer) readFromUdp() {
	var msg string
	data := make([]byte, 4096)
	read, remoteAddr, err := u.listener.ReadFromUDP(data)
	if err == nil {
		msg = fmt.Sprintf("Read %d From %s Succ: %s", read, udpaddr2str(remoteAddr), string(data[0:read]))
		u.Logger(msg)
		u.writeToUdp(remoteAddr)
	} else {
		msg = fmt.Sprintf("Read 0 From %s Fail:", udpaddr2str(remoteAddr))
		u.Logger(msg)
	}

}

func (u *EasyUdpServer) writeToUdp(remoteAddr *net.UDPAddr) {
	var msg string
	data := []byte("OK")
	write, err := u.listener.WriteToUDP(data, remoteAddr)
	if err == nil {
		msg = fmt.Sprintf("Write %d To %s SUcc: %s", write, udpaddr2str(remoteAddr), string(data[0:write]))
	} else {
		msg = fmt.Sprintf("Write 0 To %s Fail:", udpaddr2str(remoteAddr))
	}
	u.Logger(msg)
}

func (u *EasyUdpServer) close() {
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

func main() {
	logger := func(s string) {
		fmt.Println(s)
	}
	server := NewEasyUdpServer(Udp4, 8082, 4, logger)
	err := server.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	server.listen()
	time.Sleep(2 * time.Minute)
}
