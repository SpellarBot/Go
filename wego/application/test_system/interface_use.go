package utils

import "fmt"

type Man struct {
	Name string
	Age int
}
func (M Man) Think(){
	fmt.Println("I can think")
}
type Woman struct {
	Name string
	Age int
}
func (W Woman) Think(){
	fmt.Println("I can think")
}

type People interface {
	Think()
}


func InterfaceTest(){
	var P People
	M := Man{}
	P = M
	P.Think()


	m := make(map[string]string)
	m["name"] = "zhang"
	m["age"] = "12"
	var I interface{}
	I = m
	fmt.Println(I)
}
