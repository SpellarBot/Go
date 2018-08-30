package main

import (
	"fmt"
	"strconv"
	"time"
	"wego/common/easyserver"
)

func main() {
	logger := func(s string) {
		fmt.Println(s)
	}
	client := easyserver.EasyUdpClient{
		UType:  easyserver.UDP4,
		Host:   "127.0.0.1",
		Port:   8082,
		Logger: logger,
	}
	err := client.Init()
	defer client.Close()
	if err == nil {
		for i := 0; i < 100; i++ {
			go client.Send([]byte(strconv.Itoa(i) + "abcdefghjiklmnopkrstuvwxyzabcdefghjiklmnopkrstuvwxyz"))
			time.Sleep(100 * time.Nanosecond)
		}
		time.Sleep(20 * time.Second)
	}
}
