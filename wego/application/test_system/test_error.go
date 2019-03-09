package main


import (
	"fmt"
	"errors"
)

func Error(Msg string) error{
	var e error = errors.New(Msg)
	return e
}

func main(){
	S := Error("It me")
	fmt.Println(S.Error())

	B := []byte("I am a boy134 &&dfrdfr")
	A := string(B)
	fmt.Println(B)
	fmt.Println(A)

}