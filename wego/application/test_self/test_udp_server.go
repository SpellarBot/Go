package main

import (
	"wego/common/easyserver"
)

func main() {
	responser := func(s []byte) []byte {
		return []byte("OK")
	}
	server := easyserver.EasyUdpServer{
		UType:       easyserver.UDP4,
		Port:        8082,
		Threads:     4,
		Responser:   responser,
		WriteBuffer: 64,
		ReadBuffer:  64,
	}
	err := server.Init()
	if err == nil {
		defer server.Close()
	}
}
