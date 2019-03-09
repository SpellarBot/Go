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
	// logWriter = logger.NewDefaultFileLogWriter("./", "logTest", 1000, true)
	logWriter = logger.NewLineFileLogWriter("./", "logTest", 10000, 1000, true, 3)
	logWriter.Init()
	go setinfo()
	go seterror()
	go setcritical()
	time.Sleep(10 * time.Second)
}
