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
		Port:         8082,
		TType:        easyserver.Tcp4,
		Threads:      4,
		WriteBuffer:  4096,
		ReadBuffer:   4096,
		Timeout:      2 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}
	server.Init()
	defer server.Close()
	wait.Wait()
}
