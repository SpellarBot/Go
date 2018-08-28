package main

import (
	"fmt"
	"time"
	"wego/common/easyserver"
)

func main() {
	logger := func(s string) {
		fmt.Println(s)
	}
	client, _ := easyserver.NewEasyUdpClient(easyserver.Udp4, "127.0.0.1", 8082, logger)
	for i := 0; i < 100; i++ {
		go client.Send("I am a test")
	}
	time.Sleep(20 * time.Second)
}
