package main

import "fmt"

type node struct {
	value int
	next *node
}


func main(){
	N := node{value:0}
	fmt.Println(N.next)
}