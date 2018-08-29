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
	defer client.Close()
	for i := 0; i < 100; i++ {
		go client.Send([]byte("abcdefghjiklmnopkrstuvwxyzabcdefghjiklmnopkrstuvwxyz"))
		time.Sleep(100 * time.Nanosecond)
	}
	time.Sleep(20 * time.Second)
}
