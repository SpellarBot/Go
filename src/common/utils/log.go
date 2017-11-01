// Log Client
// @Author: Golion
// @Date: 2017.5

package utils

import (
	"fmt"
	"time"
	"os"
	"sync"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/theckman/go-flock"
)

const (
	defaultLogPrefix               = "logclient_default"
	defaultTimeZone                = "Asia/Shanghai"
	defaultFlushBufferLimit        = 1000
	defaultFlushBufferMilliseconds = 300
	maxRecoverTryCnt               = 3
	autoRotateTrySeconds           = 1
)

type LogClient struct {
	LogDir            string // 日志存储目录
	LogPrefix         string // 日志前缀
	LogLevel          string // NONE/ERROR/INFO/WARNING/DEBUG/ALL
	MaxDays           int    // 设为0，则只存一个文件；大于0，则每天会把xxx.log移到xxx.YYYY-MM-DD.log，且自动删掉过期日志
	TimeZone          string // 时区。默认是北京时间。
	FlushLimit        int    // 日志Flush条数限制
	FlushMilliSeconds int    // 日志Flush时间限制

	f                 *os.File
	locale            *time.Location
	bufferChan        chan log
	bufferChanSize    int32
	bufferChanLimit   int32
	bufferChanCnt     int32
	bufferIndex       int
	buffer1           []string
	buffer2           []string
	flushMutex        sync.RWMutex
	switchMutex       sync.RWMutex
	rotateMutex       sync.RWMutex
	current           string
	flushTicker       *time.Ticker
	rotateTicker      *time.Ticker
	fatal             bool
	flushing          bool
	rotating          bool
	recoverTryCnt     int
	fileLock          *flock.Flock
}

func NewLogClient(logDir string, logPrefix string, logMaxDays int) *LogClient {
	logClient := LogClient{
		LogDir:    logDir,
		LogPrefix: logPrefix,
		MaxDays:   logMaxDays,
	}
	logClient.Init()
	return &logClient
}

type log struct {
	Level  string
	Msg    string
	Caller string
	Line   int
}

func (l *LogClient) Init() {
	var err error

	// 检查目录
	if l.LogDir == "" {
		// 检查GOPATH
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			fmt.Printf("[LogClient][Init] Fatal Error! GOPATH Is Empty!\n")
			l.fatal = true
			return
		}
		l.LogDir = gopath + "/logs"
	}
	if !IsDirExist(l.LogDir) {
		err = os.Mkdir(l.LogDir, 0666)
		if err != nil {
			fmt.Printf("[LogClient][Init] Fatal Error! Dir Not Exist And Mkdir Failed! dir=[%v] error=[%v]\n", l.LogDir, err.Error())
			l.fatal = true
			return
		} else {
			fmt.Printf("[LogClient][Init] Mkdir! dir=[%v]\n", l.LogDir)
		}
	}

	// 检查配置
	if l.LogPrefix == "" {
		l.LogPrefix = defaultLogPrefix
	}
	if l.TimeZone == "" {
		l.TimeZone = defaultTimeZone
	}
	if l.FlushLimit <= 0 {
		l.FlushLimit = defaultFlushBufferLimit
	}
	if l.FlushMilliSeconds <= 0 {
		l.FlushMilliSeconds = defaultFlushBufferMilliseconds
	}

	// 设置时区
	l.locale, err = time.LoadLocation(l.TimeZone)
	if err != nil {
		fmt.Printf("[LogClient][Init] Fatal Error! error=[%v]\n", err.Error())
		l.fatal = true
	} else {
		fmt.Printf("[LogClient][Init] Use Default Timezone! timeZone=[%v]\n", defaultTimeZone)
		l.locale, _ = time.LoadLocation(defaultTimeZone)
	}

	// 时间
	l.current = time.Now().In(l.locale).Format("2006-01-02")

	// 缓冲区
	l.bufferChanCnt = 0
	l.bufferChanSize = defaultFlushBufferLimit * 100
	l.bufferChanLimit = defaultFlushBufferLimit * 90
	l.bufferChan = make(chan log, l.bufferChanSize)
	l.buffer1 = []string{}
	l.buffer2 = []string{}
	go l.runBufferConsumer()

	l.initFile()

	// 开启定时flush缓冲区到文件的协程
	if l.flushTicker == nil {
		go l.autoFlush()
	}

	// 开启定时搬运日志文件的协程
	if (l.rotateTicker == nil) && (l.MaxDays > 0) {
		go l.autoRotate()
	}
}

func (l *LogClient) initFile() {
	var err error
	fileName := l.LogDir + "/" + l.LogPrefix + ".log"
	if !IsFileExist(fileName) {
		os.Create(fileName)
	}
	l.f, err = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("[LogClient][Init] Fatal Error! error=[%v]\n", err.Error())
		l.fatal = true
	} else {
		fmt.Printf("[LogClient] Start! fileName=[%v]\n", fileName)
		l.Infof("Start! fileName=[%v]", fileName)
	}
}

func (l *LogClient) Debugf(format string, a ...interface{}) {
	l.log("Debug", fmt.Sprintf(format, a...))
}

func (l *LogClient) Debug(msg string) {
	l.log("Debug", msg)
}

func (l *LogClient) Warning(msg string) {
	l.log("Warning", msg)
}

func (l *LogClient) Warningf(format string, a ...interface{}) {
	l.log("Warning", fmt.Sprintf(format, a...))
}

func (l *LogClient) Info(msg string) {
	l.log("Info", msg)
}

func (l *LogClient) Infof(format string, a ...interface{}) {
	l.log("Info", fmt.Sprintf(format, a...))
}

func (l *LogClient) Errorf(format string, a ...interface{}) {
	l.log("Error", fmt.Sprintf(format, a...))
}

func (l *LogClient) Error(msg string) {
	l.log("Error", msg)
}

func (l *LogClient) getLevelCode(level string) int {
	switch strings.ToUpper(level) {
	case "NONE":
		return 5
	case "ERROR":
		return 4
	case "INFO":
		return 3
	case "WARNING":
		return 2
	case "DEBUG":
		return 1
	case "ALL":
		return 0
	default:
		return 0
	}
}

func (l *LogClient) log(level string, msg string) {
	if l.fatal {
		return
	}
	if l.LogLevel != "" && l.getLevelCode(l.LogLevel) > l.getLevelCode(level) {
		return
	}
	if l.bufferChanCnt >= l.bufferChanLimit {
		return
	}
	if l.f == nil {
		l.Init()
	}
	atomic.AddInt32(&l.bufferChanCnt, 1)
	if pc, _, line, ok := runtime.Caller(2); ok {
		l.bufferChan <- log{Level: level, Msg: msg, Caller: runtime.FuncForPC(pc).Name(), Line: line}
	} else {
		l.bufferChan <- log{Level: level, Msg: msg}
	}
}

func (l *LogClient) getCurrentBuffer() *[]string {
	l.switchMutex.RLock()
	defer l.switchMutex.RUnlock()
	if l.bufferIndex == 1 {
		return &l.buffer1
	} else {
		return &l.buffer2
	}
}

func (l *LogClient) getAnotherBuffer() *[]string {
	l.switchMutex.RLock()
	defer l.switchMutex.RUnlock()
	if l.bufferIndex == 1 {
		return &l.buffer2
	} else {
		return &l.buffer1
	}
}

func (l *LogClient) switchBufferIndex() {
	l.switchMutex.Lock()
	defer l.switchMutex.Unlock()
	if l.bufferIndex == 1 {
		l.bufferIndex = 0
	} else {
		l.bufferIndex = 1
	}
}

func (l *LogClient) runBufferConsumer() {
	for log := range l.bufferChan {
		atomic.AddInt32(&l.bufferChanCnt, -1)
		l.printLog(log.Level, log.Msg, log.Caller, log.Line)
	}
}

func (l *LogClient) printLog(level string, msg string, caller string, line int) {
	currentYMD := time.Now().In(l.locale).Format("2006-01-02")
	if (!l.rotating) && (l.MaxDays > 0) && (currentYMD != l.current) {
		l.flush()
		go l.Rotate(currentYMD)
	}
	currentHMS := time.Now().In(l.locale).Format("15:04:05")
	logContent :=
		"[" + level + "]" +
		"[" + currentYMD + " " + currentHMS + "]" +
		"[" + caller + ":" + fmt.Sprintf("%d", line) + "]" +
		" " + msg + "\n";
	l.flushMutex.RLock()
	currentBuffer := l.getCurrentBuffer()
	*currentBuffer = append(*currentBuffer, logContent)
	l.flushMutex.RUnlock()
	if (!l.flushing) && (!l.rotating) && (len(*currentBuffer) > l.FlushLimit) {
		l.flush()
	}
}

// 定时检查flush
func (l *LogClient) autoFlush() {
	l.flushTicker = time.NewTicker(time.Duration(l.FlushMilliSeconds) * time.Millisecond)
	for _ = range l.flushTicker.C {
		if (!l.flushing) && (!l.rotating) && (len(*l.getCurrentBuffer()) > 0) {
			l.flush()
		}
	}
}

// 退出时触发
// 正确用法：在main()里面defer LogClient.CleanUp()
func (l *LogClient) CleanUp() {
	l.flush()
	l.fatal = true
	l.recoverTryCnt = maxRecoverTryCnt + 1
	l.f.Close()
	l.f = nil
}

// 把缓冲区的日志一次写入文件
func (l *LogClient) flush() {
	l.flushMutex.Lock()
	l.rotateMutex.RLock()
	if len(*l.getCurrentBuffer()) > 0 {
		l.flushing = true
		l.switchBufferIndex()
		anotherBuffer := l.getAnotherBuffer()
		_, err := l.f.WriteString(strings.Join(*anotherBuffer, ""))
		if err != nil {
			fmt.Printf("[LogClient][flush] Fatal Error! Flush Failed! error=[%v]\n", err.Error())
			l.fatal = true
		}
		*anotherBuffer = []string{}
		l.flushing = false
	}
	l.rotateMutex.RUnlock()
	l.flushMutex.Unlock()
}

// 定时检查rotate
func (l *LogClient) autoRotate() {
	l.rotateTicker = time.NewTicker(time.Duration(autoRotateTrySeconds) * time.Second)
	for _ = range l.rotateTicker.C {
		l.checkRecover()
		currentYMD := time.Now().In(l.locale).Format("2006-01-02")
		if (!l.rotating) && (currentYMD != l.current) {
			l.Rotate(currentYMD)
		}
	}
}

// 自动恢复
func (l *LogClient) checkRecover() {
	if l.fatal {
		l.recoverTryCnt++
		if l.recoverTryCnt > maxRecoverTryCnt {
			return
		}
		blankFileName := l.LogDir + "/" + l.LogPrefix + ".log"
		if !IsFileExist(blankFileName) {
			os.Create(blankFileName)
			fmt.Printf("[LogClient][checkRecover] create %v\n", blankFileName)
		}
		var err error
		l.f, err = os.OpenFile(blankFileName, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("[LogClient][checkRecover][OpenFile] Fatal Error! error=[%v]\n", err.Error())
		} else {
			l.fatal = false
			l.recoverTryCnt = 0
		}
	}
}

// 把xxx.log移到xxx.YYYY-MM-DD.log，且自动删掉过期日志
func (l *LogClient) Rotate(current string) {
	l.rotateMutex.Lock()
	defer l.rotateMutex.Unlock()
	if l.current == current {
		return
	}
	l.rotating = true
	fmt.Printf("[LogClient][Rotate] Start Rotate current=[%v]\n", current)

	oriFileName := l.LogDir + "/" + l.LogPrefix + ".log"
	lockFileName  := l.LogDir + "/" + l.LogPrefix + ".lock"
	backFileName  := l.LogDir + "/" + l.LogPrefix + "." + l.current + ".log"

	l.f.Close()

	if !IsFileExist(backFileName) {
		// 使用文件锁，避免竞争
		l.fileLock = flock.NewFlock(lockFileName)
		if locked, _ := l.fileLock.TryLock(); locked {

			// 删除过期日志
			outdatedFileName := l.LogDir + "/" + l.LogPrefix + "." + time.Unix(time.Now().Unix() - int64(86400 * (l.MaxDays + 1)), 0).In(l.locale).Format("2006-01-02") + ".log"
			if IsFileExist(outdatedFileName) {
				fmt.Printf("[LogClient][Rotate] rm %v\n", outdatedFileName)
				Remove(outdatedFileName)
			}

			// 将xxx.log移动到xxx.YYYY-MM-DD.log
			Exec("mv", oriFileName, backFileName)
			fmt.Printf("[LogClient][Rotate] mv %v %v\n", oriFileName, backFileName)

			// 创建新文件
			Exec("cp", "/dev/null", oriFileName)
			fmt.Printf("[LogClient][Rotate] cp /dev/null %v\n", oriFileName)

			// 延时解除文件锁&删除锁文件
			go l.closeFileLock(lockFileName)

		} else {
			// 竞争失败，睡眠一段时间，等待竞争成功者完成操作
			time.Sleep(time.Duration(defaultFlushBufferMilliseconds) * time.Millisecond)
		}
	}

	l.initFile()

	l.current = current
	l.rotating = false
}

func (l *LogClient) closeFileLock(lockFileName string) {
	time.Sleep(time.Duration(autoRotateTrySeconds * 10) * time.Second)
	l.fileLock.Unlock()
	Remove(lockFileName)
}