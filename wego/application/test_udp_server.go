package main

import (
	"fmt"
	"os"
	"sync"
	"wego/common/easyserver"
)

func main() {
	wait := sync.WaitGroup{}
	wait.Add(1)
	logger := func(s string) {
		fmt.Println(s)
	}
	responser := func(s []byte) []byte {
		return []byte("OK")
	}
	server := easyserver.EasyUdpServer{
		UType:       easyserver.UDP4,
		Port:        8082,
		Threads:     4,
		Logger:      logger,
		Responser:   responser,
		WriteBuffer: 64,
		ReadBuffer:  64,
	}
	err := server.Init()
	if err == nil {
		defer server.Close()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(server)
	}
	wait.Wait()
}
