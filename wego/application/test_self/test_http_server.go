package main

import (
	"fmt"
	"net/http"
	"wego/common/easyserver"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println("Test For Writer")
	fmt.Println("URL:    ", r.URL)
	fmt.Println("Form:   ", r.Form)
	fmt.Println("Body:   ", r.Body)
	fmt.Println("Header: ", r.Header)
	fmt.Println("Host:   ", r.Host)
	fmt.Println("Method: ", r.Method)
	n, err := w.Write([]byte("Hello world"))
	fmt.Println(n, err)
}

func printNameHandler(w http.ResponseWriter, r *http.Request) {
	n, err := w.Write([]byte(r.URL.Path))
	fmt.Println(n, err)
}

func main() {
	server := easyserver.EasyHttpServe{}
	err := server.Init()
	if err == nil {
		server.AddRouter("/name/", printNameHandler)
		server.AddRouter("/", helloWorldHandler)
		server.Serve()
		defer server.Close()
	}

}
