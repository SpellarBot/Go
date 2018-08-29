package main

import (
	"fmt"
	"time"
	"wego/common/easyserver"
)

func testTcp() {
	logger := func(s string) {
		fmt.Println(s)
	}
	client, _ := easyserver.NewEasyTcpClient(easyserver.Tcp4, "127.0.0.1", 8082, logger)
	defer client.Close()
	client.Send([]byte("abcdefghjiklmnopkrstuvwxyzabcdefghjiklmnopkrstuvwxyz"))

}

func main() {
	for i := 0; i < 100; i++ {
		testTcp()
	}
	time.Sleep(20 * time.Second)
}

// func main() {
// 	var localaddr, remoteaddr *net.TCPAddr
// 	localaddr, _ = net.ResolveTCPAddr("tcp4", "127.0.0.1:8083")
// 	remoteaddr, _ = net.ResolveTCPAddr("tcp4", "127.0.0.1:8082")
// 	conn, _ := net.DialTCP("tcp4", localaddr, remoteaddr)
// 	net.Dial()
// 	for {
// 		writemsg := "question"
// 		_, err1 := conn.Write([]byte(writemsg))
// 		if err1 == nil {
// 			fmt.Println("Write Succ")
// 			buff := make([]byte, 128)
// 			read, err2 := conn.Read(buff)
// 			if err2 != nil {
// 				fmt.Println("Read Fail")
// 				break
// 			}
// 			readmsg := string(buff[0:read])
// 			fmt.Println(fmt.Sprintf("Get %s From %s", readmsg, conn.RemoteAddr().String()))
// 		} else {
// 			fmt.Println("Write Fail:" + err1.Error())
// 		}

// 	}
// }
