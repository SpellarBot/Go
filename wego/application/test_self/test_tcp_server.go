package main

import (
	"time"
	"wego/common/easyserver"
)

func main() {
	server := easyserver.EasyTcpServer{
		Port:          8082,
		TType:         easyserver.TCP4,
		Threads:       4,
		WriteBuffer:   4096,
		ReadBuffer:    4096,
		Timeout:       100 * time.Second,
		ReadTimeout:   10 * time.Second,
		WriteTimeout:  10 * time.Second,
		KeepAliveTime: 1 * time.Second,
	}
	err := server.Init()
	if err == nil {
		defer server.Close()
	}
}
