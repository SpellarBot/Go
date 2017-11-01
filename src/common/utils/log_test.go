// Log Client Test Cases
// @Author: Golion
// @Date: 2017.7

package utils

import (
	"testing"
	"time"
)

func TestLogClient_Rotate(t *testing.T) {
	logClient := NewLogClient("/tmp", "test_logclient_rotate", 3)
	go func() {
		ticker := time.NewTicker(time.Duration(100) * time.Microsecond)
		for _ = range ticker.C {
			logClient.Debug(NewRandLenChars(10))
		}
	}()
	go func() {
		ticker := time.NewTicker(time.Duration(20) * time.Second)
		for _ = range ticker.C {
			logClient.Rotate(NewRandLenChars(5))
		}
	}()
	select {}
}

func TestLogClient_LogLevel(t *testing.T) {
	logClient := NewLogClient("/tmp", "test_logclient_loglevel", 3)
	logClient.LogLevel = "INFO"
	logClient.Error("Error")
	logClient.Info("Info")
	logClient.Warning("Warning")
	logClient.Debug("Debug")
	select {}
}