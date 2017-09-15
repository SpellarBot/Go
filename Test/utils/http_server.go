package utils

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"io/ioutil"
	"encoding/json"
)




func Controller(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析参数, 默认是不会解析的
	M := make(map[string]string)
	fmt.Fprintf(w, "Hello I am your new server!\n") //输出到客户端的信息
	if r.Method == "GET"{
		for key,value := range r.Form{
			M[key] = strings.Join(value,"")
			fmt.Fprintf(w,fmt.Sprintf("Get key: %s\n",M[key]))
		}
	}
	if r.Method == "POST"{
		S,_ := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		var M map[string]interface{}
		json.Unmarshal([]byte(S),&M)
		for key,value := range M{
			fmt.Fprintf(w,fmt.Sprintf("Post %s: %s\n",key,value))
		}
		}
	}



func Http_Server_Test() {
	http.HandleFunc("/", Controller)       //设置访问的路由
	err := http.ListenAndServe(":1521", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

