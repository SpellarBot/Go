package main

import (
	"vidmate.com/ligo/ligo_ctr/common/services"
	"fmt"
)

func main(){
	C := services.CTRModel{}
	C.Init(0,0,0,"IN")
	fmt.Println(C.Country)
}