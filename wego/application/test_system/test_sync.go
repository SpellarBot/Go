package main

import (
	"fmt"
	"time"
	"sync"
)


func ChanTest(){
	var C chan int
	C = make(chan int,4)
	go func(){
		for{
			time.Sleep(6*time.Second)
			C <- 1
		}
	}()
	for{
		select{
		case <-C:
			fmt.Println("get a data")
		case <-time.After(5*time.Second):
				fmt.Println("wati 5 s")
		}
	}
}


func WaitGroupTest(){
	var W sync.WaitGroup
	W.Add(15)
	for i:=0;i<5;i++{
		go func(){
			fmt.Println("finish a fun")
			W.Done()
		}()
	}
	W.Wait()
}



func main(){
	//go ChanTest()
	//go WaitGroupTest()
	//time.Sleep(20*time.Second)
	var W sync.WaitGroup
	for i:=0;i<5;i++{
		go func(){
			fmt.Println("begin a fun")
			W.Add(1)
		}()
	}
	for i:=0;i<5;i++{
		go func(){
			fmt.Println("finish a fun")
			W.Done()
		}()
	}
	W.Wait()
}