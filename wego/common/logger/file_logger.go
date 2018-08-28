package logger

// 文件日志记录器
// 先Init，再AddLogger
type FileLogger map[string]*FileLogWriter

func NewFileLogger() FileLogger {
	return make(FileLogger)
}

func (l FileLogger) AddLogger(moudle string, path string, filename string, bufferlength int, console bool) error {
	fileLogWriter, err := NewDefaultFileLogWriter(path, filename, bufferlength, console)
	l[moudle] = fileLogWriter
	return err
}

func (l FileLogger) AddHourlyLogger(moudle string, path string, filename string, bufferlength int, console bool, maxbackup int) error {
	fileLogWriter, err := NewHourlyFileLogWriter(path, filename, bufferlength, console, maxbackup)
	l[moudle] = fileLogWriter
	return err
}

func (l FileLogger) AddDailyLogger(moudle string, path string, filename string, bufferlength int, console bool, maxbackup int) error {
	fileLogWriter, err := NewDailytFileLogWriter(path, filename, bufferlength, console, maxbackup)
	l[moudle] = fileLogWriter
	return err
}

func (l FileLogger) AddSizeLogger(moudle string, path string, filename string, maxsize int64, bufferlength int, console bool, maxbackup int) error {
	fileLogWriter, err := NewSizeFileLogWriter(path, filename, maxsize, bufferlength, console, maxbackup)
	l[moudle] = fileLogWriter
	return err
}

func (l FileLogger) AddLineLogger(moudle string, path string, filename string, maxline int64, bufferlength int, console bool, maxbackup int) error {
	fileLogWriter, err := NewLineFileLogWriter(path, filename, maxline, bufferlength, console, maxbackup)
	l[moudle] = fileLogWriter
	return err
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
