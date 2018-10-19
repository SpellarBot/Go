package main

import "fmt"
import "time"

func test() {
	panic("error exists")
}

func main() {
	fmt.Println("I am tets")
	go test()
	time.Sleep(3 * time.Second)
}
