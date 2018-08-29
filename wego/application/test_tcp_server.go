package main

import (
	"sync"
	"wego/common/easyserver"
)

func main() {
	wait := sync.WaitGroup{}
	wait.Add(1)
	server := easyserver.EasyTcpServer{
		Port:    8082,
		TType:   easyserver.Tcp4,
		Threads: 4,
	}
	server.Init()
	defer server.Close()
	wait.Wait()
}
