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
	server, err := easyserver.NewEasyUdpServer(easyserver.Udp4, 8082, 4, 64, 64, responser, logger)
	defer server.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(server)
	wait.Wait()
}
