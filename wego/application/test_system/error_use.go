package utils

import (
	"errors"
	"fmt"
)

func String_Error(a interface{}) error{
	_,err := a.(string)
	if !err {
		return errors.New("Your should input a string")
	}
	return nil
}
func Error_Test(){
	a := 10
	err := String_Error(a)
	if err!=nil{
		fmt.Println(err)
	}
}