package controllers

import (
	"time"
	"fmt"
)
var(
	Messages chan string
	DeleteStatus bool
	TempConsumerNum int64
	T1 int64
	T0 int64
	MaxConsumerNum int64
)

func GetError(e error){
	logClient.Error("[CTR][Kafka Error]"+e.Error())
}

func GetMessage(Message []byte){
	Messages <- string(Message)
}

func GetLog(){
	kafkaClient.RunConsumerFromNewest(GetMessage, GetError)
}

func AddLog(){
	var t1,t2 int64
	for {
		switch *status {
		case "Work":
			Message := <-Messages
			ctrModel.AddLog(Message)
			TempConsumerNum ++
		case "Stop":
			t1 = time.Now().UnixNano()
			ctrModel.Save(*modelFile)
			*status = "Pause"
			t2 = time.Now().UnixNano()
			logClient.Info(fmt.Sprintf("Stop and save cost %f ms",float64(t2-t1)/1e6))
		case "Restore":
			t1 = time.Now().UnixNano()
			*status = "Work"
			t2 = time.Now().UnixNano()
			logClient.Info(fmt.Sprintf("Restore cost %f ms",float64(t2-t1)/1e6))
		}
		if DeleteStatus == true {
			ctrModel.Delete()
			DeleteStatus = false
			logClient.Info("Clear Videos Whithout Consume")
		}
		if TempConsumerNum == MaxConsumerNum{
			t1 = time.Now().UnixNano()
			ctrModel.Save(*modelFile)
			t2 = time.Now().UnixNano()
			logClient.Info(fmt.Sprintf("Save Success cost %f ms",float64(t2-t1)/1e6))
			TempConsumerNum = 0
			T1 := time.Now().UnixNano()
			T := T1- T0
			T0 = T1
			logClient.Info(fmt.Sprintf("Consume Speed is %d/s",int(float64(MaxConsumerNum)/(float64(T)/1e9))))
		}
	}
}

func Delete(){
	for {
		t1 := time.NewTimer(6*time.Hour)
		<-t1.C
		DeleteStatus = true

	}
}


// 根据状态循环读取kafka/暂停/初始化，并定期清理
func Consumer(){
	Messages = make(chan string, 20000)
	DeleteStatus = false
	TempConsumerNum = 0
	MaxConsumerNum = 10000
	T1 = 0
	T0 = time.Now().UnixNano()
	time.Sleep(10*time.Second)
	go GetLog()
	go AddLog()
	go Delete()
}