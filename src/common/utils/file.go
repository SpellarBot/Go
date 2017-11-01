package utils

import (
	"os"
	"bufio"
	"fmt"
)

type File struct {
	Filename string
}
//read
func (F File )Read2Array (M []string )(int){
	f1,err := os.Open(F.Filename)
	if err != nil{
		return -1
	} else {
		reader := bufio.NewReader(f1)
		k := 0
		for{
			S,er := reader.ReadString('\n')
			if er != nil{
				break
			} else{
				M[k] = S
				k++
			}
		}
		f1.Close()
		return 0
	}
}
//create and write
func (F File)Create_Write(A string){
	f1,_ := os.Create(F.Filename)
	writer := bufio.NewWriter(f1)
	for i:=1;i<10;i++{
		writer.WriteString(fmt.Sprintf("%s\r\n",A))
		//fmt.Fprintln(writer,A)
	}
	writer.Flush()
	f1.Close()
	return
}

//write in the next
func (F File)Append_Write(A string){
	f1,_ := os.OpenFile(F.Filename,os.O_APPEND|os.O_WRONLY,0644)
	writer := bufio.NewWriter(f1)
	for i:=1;i<10;i++{
		writer.WriteString(fmt.Sprintf("%s\r\n",A))
		//fmt.Fprintln(writer,A)
	}
	writer.Flush()
	f1.Close()
	return
}
