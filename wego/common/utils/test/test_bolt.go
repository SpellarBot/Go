package main


import (
	"os"
	"strconv"
	"fmt"
	"bufio"
)


type People struct{
	name string
	age int
	number []int
}

func main(){
	P := People{name:"zhang",age:18}
	nu := []int{4,5,6}
	P.number = nu

	file,err := os.Open("test.dat")
	fmt.Println(err)
	if err == nil{
		os.Remove("test.dat")
	}
	file,err = os.Create("test.dat")
	defer file.Close()

	os.Create("test.dat")
	Writer := bufio.NewWriter(file)

	Writer.WriteString(P.name)
	Writer.WriteString(strconv.Itoa(P.age))

	for _,j := range P.number {
		Writer.WriteString(strconv.Itoa(j))
	}
	Writer.Flush()
	fmt.Println("end")
}