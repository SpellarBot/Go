package utils

import (
	"encoding/json"
	"fmt"
)
type Json struct{
	Content string
	key []string
}

// json to map
func (J *Json) Parse()(map[string]interface{}){
	var f interface{}
	json.Unmarshal([]byte(J.Content), &f)
	m := f.(map[string]interface{})
	for key,_ := range m{
		J.key = append(J.key,key)
	}
	return m
}

//map to json
func (J *Json) Set(M map[string]interface{}){
	D,e := json.Marshal(M)
	if e == nil {
		J.Content = string(D)
	}
}


func JsonTest(){
	S := `{"nAmE":"Bob","Food":"Pickle", "inner":"changed"}`
	J := Json{
		Content:S,
	}
	m := J.Parse()
	for i,j:= range m{
		fmt.Println(i,j)
	}

	J.Set(m)
	fmt.Println(J.Content)

}