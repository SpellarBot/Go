package main

import "fmt"
import "time"

func main() {
	t1 := time.NewTimer(time.Second * 5)
	fmt.Println(time.Now())
	fmt.Println("====")
	t := <-t1.C
	fmt.Println(t)
	fmt.Println("====")
	fmt.Println(time.Now())
}