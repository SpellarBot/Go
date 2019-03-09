package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	Lock sync.RWMutex
)
type I struct {
	i int
}

func (i *I) Read(){
	Lock.RLock()
	fmt.Println("begin to read")
	fmt.Println(i.i)
	fmt.Println("end to read")
	Lock.RUnlock()
}

func (i *I)Write(){
	fmt.Println("begin to write")
	i.i++
	time.Sleep(2*time.Second)
	fmt.Println(i.i)
	fmt.Println("end to write")
}

func main(){
	i := I{}
	go func(){
		fmt.Println("+++")
		i.Write()
		fmt.Println("+++")
	}()

	i.Write()
	i.Read()

	time.Sleep(10*time.Second)
}