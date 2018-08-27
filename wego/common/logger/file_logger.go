package logger

// 文件日志记录器
// 先Init，再AddLogger
type FileLogger map[string]*FileLogWriter

func NewFileLogger() FileLogger {
	return make(FileLogger)
}

func (l FileLogger) AddLogger(moudle string, path string, filename string, bufferlength int, console bool) {
	fileLogWriter := NewDefaultFileLogWriter(path, filename, bufferlength, console)
	fileLogWriter.Init()
	l[moudle] = fileLogWriter
}

func (l FileLogger) AddHourlyLogger(moudle string, path string, filename string, bufferlength int, console bool, maxbackup int) {
	fileLogWriter := NewHourlyFileLogWriter(path, filename, bufferlength, console, maxbackup)
	fileLogWriter.Init()
	l[moudle] = fileLogWriter
}

func (l FileLogger) AddDailyLogger(moudle string, path string, filename string, bufferlength int, console bool, maxbackup int) {
	fileLogWriter := NewDailytFileLogWriter(path, filename, bufferlength, console, maxbackup)
	fileLogWriter.Init()
	l[moudle] = fileLogWriter
}

func (l FileLogger) AddSizeLogger(moudle string, path string, filename string, maxsize int64, bufferlength int, console bool, maxbackup int) {
	fileLogWriter := NewSizeFileLogWriter(path, filename, maxsize, bufferlength, console, maxbackup)
	fileLogWriter.Init()
	l[moudle] = fileLogWriter
}

func (l FileLogger) AddLineLogger(moudle string, path string, filename string, maxline int64, bufferlength int, console bool, maxbackup int) {
	fileLogWriter := NewLineFileLogWriter(path, filename, maxline, bufferlength, console, maxbackup)
	fileLogWriter.Init()
	l[moudle] = fileLogWriter
}

func (l FileLogger) Close() {
	for name, filt := range l {
		filt.Close()
		delete(l, name)
	}
}

// 获取一个FileLogWriter对象
func (l FileLogger) GetWriter(moudle string) *FileLogWriter {
	return l[moudle]
}

// 获取FileLogWriter对象的函数
func (l FileLogger) GetLogFun(level Level, moudle string) func(string) {
	return l[moudle].GetLogFun(level)
}
