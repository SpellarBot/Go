package logger

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Level int

const (
	FINEST   = Level(0)
	FINE     = Level(1)
	DEBUG    = Level(2)
	TRACE    = Level(3)
	INFO     = Level(4)
	WARNING  = Level(5)
	ERROR    = Level(6)
	CRITICAL = Level(7)
)

var (
	levels = map[Level]string{
		FINEST:   "FINEST",
		FINE:     "FINE",
		DEBUG:    "DEBUG",
		TRACE:    "TRACE",
		INFO:     "INFO",
		WARNING:  "WARNING",
		ERROR:    "ERROR",
		CRITICAL: "CRITICAL",
	}
)

type Locate struct {
	FileName string
	Func     string
	Line     int
}

type LogRecord struct {
	LogLevel Level     // The log level
	Location Locate    // The Log's location
	Created  time.Time // The time at which the log message was created (nanoseconds)
	Message  string    // The log message
}

// 单个日志记录器
// 日志Roate规则优先级 hourly > daily > size = lines
type FileLogWriter struct {
	rec  chan *LogRecord
	read chan bool

	// The opened file
	Path       string
	FilePrefix string
	file       *os.File

	// print in console
	Console bool

	// buffer length
	BufferLength int

	// Keep old logfiles (.001, .002, etc)
	Maxbackup int

	// Rotate at linecount
	Maxline        int64
	maxlineCurline int64

	// Rotate at size
	Maxsize        int64
	maxsizeCursize int64

	// Current file nums
	currentFileNums int64

	// Rotate daily/hourly
	Daily   bool
	Hourly  bool
	logTime time.Time

	// Define log level and rotate type
	MinLevel   Level
	rotateType string

	// Lock
	rotateMutex sync.RWMutex
}

// a logwriter without rotate
func NewDefaultFileLogWriter(path string, filename string, bufferlength int, console bool) *FileLogWriter {
	var fileLogWriter FileLogWriter
	fileLogWriter = FileLogWriter{
		Path:         path,
		FilePrefix:   filename,
		Console:      console,
		Hourly:       false,
		Daily:        false,
		Maxsize:      0,
		Maxline:      0,
		BufferLength: bufferlength,
		Maxbackup:    168,
		MinLevel:     FINEST,
	}
	return &fileLogWriter
}

// a hourly logwriter
func NewHourlyFileLogWriter(path string, filename string, bufferlength int, console bool, maxbackup int) *FileLogWriter {
	fileLogWriter := NewDefaultFileLogWriter(path, filename, bufferlength, console)
	fileLogWriter.SetHourly()
	fileLogWriter.SetMaxbackup(maxbackup)
	return fileLogWriter
}

// a daily logwriter
func NewDailytFileLogWriter(path string, filename string, bufferlength int, console bool, maxbackup int) *FileLogWriter {
	fileLogWriter := NewDefaultFileLogWriter(path, filename, bufferlength, console)
	fileLogWriter.SetDaily()
	fileLogWriter.SetMaxbackup(maxbackup)
	return fileLogWriter
}

// a maxsize logwriter
func NewSizeFileLogWriter(path string, filename string, maxsize int64, bufferlength int, console bool, maxbackup int) *FileLogWriter {
	fileLogWriter := NewDefaultFileLogWriter(path, filename, bufferlength, console)
	fileLogWriter.SetMaxSize(maxsize)
	fileLogWriter.SetMaxbackup(maxbackup)
	return fileLogWriter
}

// a maxline logwriter
func NewLineFileLogWriter(path string, filename string, maxline int64, bufferlength int, console bool, maxbackup int) *FileLogWriter {
	fileLogWriter := NewDefaultFileLogWriter(path, filename, bufferlength, console)
	fileLogWriter.SetMaxLine(maxline)
	fileLogWriter.SetMaxbackup(maxbackup)
	return fileLogWriter
}

func (w *FileLogWriter) Init() {
	switch {
	case w.Hourly:
		w.rotateType = "Hourly"
	case w.Daily:
		w.rotateType = "Daily"
	case w.Maxsize > 0:
		w.rotateType = "Maxsize"
	case w.Maxline > 0:
		w.rotateType = "Maxline"
	default:
		w.rotateType = "None"
	}
	w.logTime = getFileCreateTime(w.Path, w.FilePrefix)
	w.file = initFile(w.Path, w.FilePrefix)
	w.maxlineCurline = getFileLine(w.Path, w.FilePrefix)
	w.maxsizeCursize = getFileSize(w.Path, w.FilePrefix)
	w.read = make(chan bool)
	w.rec = make(chan *LogRecord, w.BufferLength)
	w.rotateMutex = sync.RWMutex{}
	go w.write()
	if w.rotateType != "None" {
		if w.Maxbackup <= 0 {
			w.Maxbackup = 100000000000
		}
		go w.delete()
	}

}

func (w *FileLogWriter) GetLogFun(level Level) func(string) {
	return func(msg string) {
		pc, file, line, ok := runtime.Caller(1)
		if ok == true {
			f := runtime.FuncForPC(pc)
			_, currentfile := filepath.Split(file)
			fun := f.Name()
			rec := LogRecord{
				Created:  time.Now(),
				LogLevel: level,
				Message:  msg,
				Location: Locate{
					FileName: currentfile,
					Line:     line,
					Func:     fun,
				},
			}
			w.rec <- &rec
		}
	}
}

func (w *FileLogWriter) Info(msg string) {
	w.addLog(INFO, msg)
}
func (w *FileLogWriter) Finest(msg string) {
	w.addLog(FINEST, msg)
}
func (w *FileLogWriter) Fine(msg string) {
	w.addLog(FINE, msg)
}
func (w *FileLogWriter) Debug(msg string) {
	w.addLog(DEBUG, msg)
}
func (w *FileLogWriter) Trace(msg string) {
	w.addLog(TRACE, msg)
}
func (w *FileLogWriter) Warning(msg string) {
	w.addLog(WARNING, msg)
}
func (w *FileLogWriter) Critical(msg string) {
	w.addLog(CRITICAL, msg)
}
func (w *FileLogWriter) Error(msg string) {
	w.addLog(ERROR, msg)
}

func (w *FileLogWriter) addLog(level Level, msg string) {
	pc, file, line, ok := runtime.Caller(2)
	if ok == true {
		f := runtime.FuncForPC(pc)
		_, currentfile := filepath.Split(file)
		fun := f.Name()
		rec := LogRecord{
			Created:  time.Now(),
			LogLevel: level,
			Message:  msg,
			Location: Locate{
				FileName: currentfile,
				Line:     line,
				Func:     fun,
			},
		}
		w.rec <- &rec
	}
}

// ToDo: delete log file in time
func (w *FileLogWriter) delete() {
	if w.rotateType == "None" {
		return
	}
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		<-ticker.C
		fileList, err := ioutil.ReadDir(w.Path)
		removeList := make([]os.FileInfo, 0, len(fileList))
		if err == nil {

			for _, file := range fileList {
				if strings.Contains(file.Name(), w.FilePrefix) {
					removeList = append(removeList, file)
				}
			}
			if len(removeList) <= w.Maxbackup+1 {
				continue
			}
			removeNum := len(removeList) - w.Maxbackup - 1
			for i := 0; i < len(removeList); i++ {
				for j := i + 1; j < len(removeList); j++ {

					if removeList[j-1].ModTime().Nanosecond() >= removeList[j].ModTime().Nanosecond() {
						tempInfo := removeList[j]
						removeList[j] = removeList[j-1]
						removeList[j-1] = tempInfo
					}
				}
			}
			for i := 0; i < removeNum; i++ {
				err = os.Remove(filepath.Join(w.Path, removeList[i].Name()))
				if err == nil {
					w.Info("Remove OutDate LogFile Succ: " + removeList[i].Name())
				} else {
					w.Info("Remove OutDate LogFile Fail: " + err.Error())
				}
			}
		}

	}

}

func (w *FileLogWriter) Close() {
	w.Flush()
	// close(w.rec)
}

func (w *FileLogWriter) write() {
	for {
		w.rotate()
		record := <-w.rec
		log := LogRecord2String(record)
		if record.LogLevel >= w.MinLevel {
			fmt.Fprintln(w.file, log)
			atomic.AddInt64(&w.maxlineCurline, 1)
			atomic.AddInt64(&w.maxsizeCursize, int64(len([]byte(log))))
		}
		if w.Console {
			fmt.Println(log)
		}
	}
}

func (w *FileLogWriter) Flush() {
	w.file.Sync()
}

func (w *FileLogWriter) rotate() {
	switch w.rotateType {
	case "Hourly":
		w.rotateHourly()
	case "Daily":
		w.rotateDaily()
	case "Maxsize":
		w.rotateMaxsize()
	case "Maxline":
		w.rotateMaxline()
	default:
		return
	}
}

func (w *FileLogWriter) rotateMaxsize() {
	if w.maxsizeCursize = getFileSize(w.Path, w.FilePrefix); w.maxsizeCursize >= w.Maxsize {
		oldfname := filepath.Join(w.Path, w.FilePrefix)
		newfname := filepath.Join(w.Path, fmt.Sprintf("%s.%s.%s", w.FilePrefix,
			w.logTime.Format("20060102150405"), strconv.Itoa(int(w.currentFileNums))))
		w.changefile(oldfname, newfname)
	}
}

func (w *FileLogWriter) rotateMaxline() {
	if w.maxlineCurline >= w.Maxline {
		oldfname := filepath.Join(w.Path, w.FilePrefix)
		newfname := filepath.Join(w.Path, fmt.Sprintf("%s.%s.%s", w.FilePrefix,
			w.logTime.Format("20060102150405"), strconv.Itoa(int(w.currentFileNums))))
		w.changefile(oldfname, newfname)
	}
}

func (w *FileLogWriter) rotateDaily() {
	currentDay := time.Now().Day()
	if w.logTime.Day() != currentDay {
		oldfname := filepath.Join(w.Path, w.FilePrefix)
		newfname := filepath.Join(w.Path, w.FilePrefix+w.logTime.Format(".20060102"))
		w.changefile(oldfname, newfname)
	}
}

func (w *FileLogWriter) rotateHourly() {
	currentHour := time.Now().Hour()
	currentDay := time.Now().Day()
	if w.logTime.Hour() != currentHour || w.logTime.Day() != currentDay {
		oldfname := filepath.Join(w.Path, w.FilePrefix)
		newfname := filepath.Join(w.Path, w.FilePrefix+w.logTime.Format(".2006010215"))
		w.changefile(oldfname, newfname)
	}
}

func (w *FileLogWriter) changefile(oldfname string, newfname string) {
	w.rotateMutex.Lock()
	defer w.rotateMutex.Unlock()
	w.file.Close()
	os.Rename(oldfname, newfname)
	w.file = initFile(w.Path, w.FilePrefix)
	atomic.StoreInt64(&w.maxlineCurline, 0)
	atomic.StoreInt64(&w.maxsizeCursize, 0)
	atomic.AddInt64(&w.currentFileNums, 1)
	w.logTime = time.Now()

}

func (w *FileLogWriter) SetMaxbackup(maxbackup int) {
	w.Maxbackup = maxbackup
}

func (w *FileLogWriter) SetConsole(console bool) {
	w.Console = console
}
func (w *FileLogWriter) SetDaily() {
	w.Daily = true
	w.Hourly = false
}
func (w *FileLogWriter) SetHourly() {
	w.Hourly = true
	w.Daily = false
}
func (w *FileLogWriter) SetMaxSize(maxsize int64) {
	if maxsize > 0 {
		w.Daily = false
		w.Hourly = false
		w.Maxline = 0
		w.Maxsize = maxsize
	} else {
		panic("maxsize must > 0")
	}
}
func (w *FileLogWriter) SetMaxLine(maxline int64) {
	if maxline > 0 {
		w.Daily = false
		w.Hourly = false
		w.Maxsize = 0
		w.Maxline = maxline
	} else {
		panic("maxline must > 0")
	}
}

func (w *FileLogWriter) SetBufferLength(bufferlength int) {
	w.BufferLength = bufferlength
}

func (w *FileLogWriter) SetLogLevel(level Level) {
	w.MinLevel = level
}

// LogRecord 2 string
func LogRecord2String(logRecord *LogRecord) string {
	sTime := logRecord.Created.Format("2006/01/02 15:04:05 MST")
	sLocation := fmt.Sprintf("%s:%s:%d", logRecord.Location.FileName, logRecord.Location.Func, logRecord.Location.Line)
	sLevel := logRecord.LogLevel
	return fmt.Sprintf("[%s] [%s] (%s) : %s", sTime, levels[sLevel], sLocation, logRecord.Message)
}

func initFile(path string, fname string) *os.File {
	filename := filepath.Join(path, fname)
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0)
	if err != nil {
		panic("Open LogFile Failed")
	}
	return file
}

func getFileSize(path string, fname string) int64 {
	filename := filepath.Join(path, fname)
	fileinfo, _ := os.Stat(filename)
	return fileinfo.Size()
}

func getFileCreateTime(path string, fname string) time.Time {
	filename := filepath.Join(path, fname)
	fileinfo, err := os.Stat(filename)
	if err == nil {
		return fileinfo.ModTime()
	}
	return time.Now()
}

func getFileLine(path string, fname string) int64 {
	filename := filepath.Join(path, fname)
	file, err := os.Open(filename)
	defer file.Close()
	if err == nil {
		k := 0
		reader := bufio.NewReader(file)
		for {
			_, _, err = reader.ReadLine()
			if err != io.EOF {
				k++
			} else {
				break
			}
		}
		return int64(k)
	}
	return 0
}
