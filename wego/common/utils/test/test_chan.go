package main

import (
	"fmt"
	"time"
)


var (
	A int
)

func produce(p chan<- int){
	for i:=0;i<100;i++{
		p <- i
		fmt.Println("produce",i)
	}
}

func consume(p <-chan int){
	for i:=0;i<10;i++{
		a:= <-p
		A += a
		fmt.Println("consume",a)
	}
}

func print(){
	for{
		fmt.Println("print",A)
	}
}

type people struct{
	name string
	age int
}


func main(){
	P := make(chan int,5)
	A = 0
	go produce(P)
	go consume(P)
	go print()
	time.Sleep(5*time.Second)
}

