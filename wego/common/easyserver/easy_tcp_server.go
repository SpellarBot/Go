package easyserver

import (
	"fmt"
	"net"
	"runtime"
	"strconv"
)

type TcpType string

const (
	Tcp  = TcpType("tcp")
	Tcp4 = TcpType("tcp4")
	Tcp6 = TcpType("tcp6")
)

type EasyTcpServer struct {
	TType     TcpType
	Port      int
	Threads   int
	Responser func([]byte) []byte
	Logger    func(string)

	addr     *net.TCPAddr
	listener *net.TCPListener
}

func (t *EasyTcpServer) Init() {
	var err error

	if t.TType == TcpType("") {
		t.TType = Tcp4
	}
	if t.Logger == nil {
		t.Logger = func(s string) {
			fmt.Println(s)
		}
	}
	if t.Responser == nil {
		t.Responser = func(s []byte) []byte {
			return []byte("OK")
		}
	}
	if t.Port < 0 {
		t.Port = 8080
	}
	if t.Threads <= 0 {
		t.Threads = runtime.NumCPU()
	}
	t.addr, err = net.ResolveTCPAddr(string(t.TType), "0.0.0.0:"+strconv.Itoa(t.Port))
	if err == nil {
		t.listener, err = net.ListenTCP(string(t.TType), t.addr)
		if err == nil {
			for i := 0; i < t.Threads; i++ {
				go t.listen()
			}
			t.Logger(fmt.Sprintf("TCP Serve Start Succ At Port %d", t.Port))
			t.Logger(fmt.Sprintf("TCP Serve Start %d Threads", t.Threads))
		} else {
			t.Logger(fmt.Sprintf("TCP Serve Start Fail: %s", err.Error()))
		}
	} else {
		t.Logger(fmt.Sprintf("TCP Serve Start Fail: %s", err.Error()))
	}

}

func (t *EasyTcpServer) Close() {
	t.listener.Close()
}

func (t *EasyTcpServer) readFromTcp(conn *net.TCPConn) ([]byte, error) {
	readdata := make([]byte, 128)
	read, err := conn.Read(readdata)
	if err == nil {
		t.Logger(fmt.Sprintf("Read %d From %s Succ: %s", read, conn.RemoteAddr().String(), string(readdata[0:read])))
	} else {
		t.Logger(fmt.Sprintf("Read Fail: %s", err.Error()))
	}
	return readdata, err
}

func (t *EasyTcpServer) writeToTcp(conn *net.TCPConn, writedata []byte) {
	write, err := conn.Write(writedata)
	if err == nil {
		t.Logger(fmt.Sprintf("Write %d To %s Succ: %s", write, conn.RemoteAddr().String(), string(writedata[0:write])))
	} else {
		t.Logger(fmt.Sprintf("Write Fail: %s", err.Error()))
	}
}

func (t *EasyTcpServer) listen() {
	var readdata, writedata []byte
	for {
		conn, err := t.listener.AcceptTCP()
		if err == nil {
			readdata, err = t.readFromTcp(conn)
			if err == nil {
				writedata = t.Responser(readdata)
				t.writeToTcp(conn, writedata)
			} else {
				t.Logger(fmt.Sprintf("Read From TCP Fail: %s", err.Error()))
			}
		} else {
			t.Logger(fmt.Sprintf("Accept Conn Fail: %s", err.Error()))
		}
	}
}
