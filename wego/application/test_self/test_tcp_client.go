package main

import (
	"fmt"
	"strconv"
	"time"
	"wego/common/easyserver"
)

func testTcp() {
	logger := func(s string) {
		fmt.Println(s)
	}
	client := easyserver.EasyTcpClient{
		TType:  easyserver.TCP4,
		Host:   "127.0.0.1",
		Port:   8082,
		Logger: logger,
	}
	err := client.Init()
	//defer client.Close()
	if err == nil {
		for i := 0; i < 100; i++ {
			go client.Send([]byte(strconv.Itoa(i) + "abcdefghjiklmnopkrstuvwxyzabcdefghjiklmnopkrstuvwxyz"))
		}
	}
}

func main() {
	for i := 0; i < 20; i++ {
		go testTcp()
	}
	time.Sleep(30 * time.Second)
}
