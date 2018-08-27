package main

import (
	"time"
	"wego/common/logger"
)

var (
	fileLogs logger.FileLogger
)

func Log(writer func(string)) {
	for i := 0; i < 10000; i++ {
		writer("I am a Info")
	}
}

func main() {
	fileLogs = logger.NewFileLogger()
	defer fileLogs.Close()
	fileLogs.AddDailyLogger("Daily", "./", "DailyLog")
	// fileLogs.AddLineLogger("Line", "./", "LineLog", 100000)
	// fileLogs.AddSizeLogger("Size", "./", "SizeLog", 100000)
	// fileLogs.AddLogger("Default", "./", "DefaultLog")

	fileLog := fileLogs.GetWriter("Daily")
	logFun := fileLogs.GetInfoLogFun("Daily")

	logFun("I am logfun")
	fileLog.Info("I am logger")
	time.Sleep(10 * time.Second)
}
