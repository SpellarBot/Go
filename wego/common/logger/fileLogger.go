package logger

// 文件日志记录器
// 先Init，再AddLogger
type FileLogger map[string]*FileLogWriter

func NewFileLogger() FileLogger {
	return make(FileLogger)
}

func (l FileLogger) AddLogger(moudle string, path string, filename string) {
	fileLogWriter := NewDefaultFileLogWriter(path, filename)
	fileLogWriter.Init()
	l[moudle] = fileLogWriter
}

func (l FileLogger) AddHourlyLogger(moudle string, path string, filename string) {
	fileLogWriter := NewHourlyFileLogWriter(path, filename)
	fileLogWriter.Init()
	l[moudle] = fileLogWriter
}

func (l FileLogger) AddDailyLogger(moudle string, path string, filename string) {
	fileLogWriter := NewDailytFileLogWriter(path, filename)
	fileLogWriter.Init()
	l[moudle] = fileLogWriter
}

func (l FileLogger) AddSizeLogger(moudle string, path string, filename string, maxsize int64) {
	fileLogWriter := NewSizeFileLogWriter(path, filename, maxsize)
	fileLogWriter.Init()
	l[moudle] = fileLogWriter
}

func (l FileLogger) AddLineLogger(moudle string, path string, filename string, maxline int64) {
	fileLogWriter := NewLineFileLogWriter(path, filename, maxline)
	fileLogWriter.Init()
	l[moudle] = fileLogWriter
}

func (l FileLogger) Close() {
	for name, filt := range l {
		filt.Close()
		delete(l, name)
	}
}

func (l FileLogger) GetLogger(moudle string) *FileLogWriter {
	return l[moudle]
}

func (l FileLogger) GetInfoWriter(moudle string) func(string) {
	return l[moudle].GetLogger()
}
