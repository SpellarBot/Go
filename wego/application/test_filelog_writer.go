package main

import (
	"time"
	"wego/common/logger"
)

var (
	logWriter *logger.FileLogWriter
)

func setinfo() {
	for {
		logWriter.Info("I am a Info")
	}
}
func seterror() {
	for {
		logWriter.Error("I am a error")
	}
}
func setcritical() {
	for {
		logWriter.Critical("I am a error")
	}
}

func main() {
	logWriter = logger.NewDefaultFileLogWriter("./", "logTest")
	logWriter.Daily = false
	logWriter.Hourly = false
	logWriter.Maxbackup = 3
	logWriter.Maxline = 200000
	logWriter.Maxsize = 0
	logWriter.Console = false
	logWriter.Init()
	go setinfo()
	go seterror()
	go setcritical()
	time.Sleep(5 * time.Second)
}
