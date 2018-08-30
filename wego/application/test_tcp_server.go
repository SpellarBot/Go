package main

import (
	"sync"
	"time"
	"wego/common/easyserver"
)

func main() {
	wait := sync.WaitGroup{}
	wait.Add(1)
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
	server.Init()
	defer server.Close()
	wait.Wait()
}
