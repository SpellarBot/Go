package logger

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type LoggerServer struct {
	port   string
	logger *FileLogger
}

var (
	port = flag.String("p", "12124", "Port number to listen on")
)

func e(err error) {
	if err != nil {
		fmt.Printf("Erroring out: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	bind, err := net.ResolveUDPAddr("udp", "0.0.0.0:"+*port)
	e(err)

	listener, err := net.ListenUDP("udp", bind)
	e(err)

	fmt.Printf("Listening to port %s...\n", *port)
	k := 0
	for {
		buffer := make([]byte, 128)
		_, remoteAddr, err := listener.ReadFromUDP(buffer)
		e(err)
		fmt.Println(strings.Replace("Read From UDP:"+string(buffer), "\n", "", -1))
		sendData := []byte(fmt.Sprintf("%s log", strconv.Itoa(k)))
		_, err = listener.WriteToUDP(sendData, remoteAddr)
		e(err)
		fmt.Println(strings.Replace("Send To UDP:"+string(sendData), "\n", "", -1))
		k++
	}
}
