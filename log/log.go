package log

import (
	"fmt"
	"log"
	"mknote/config"
	"os"
	"path/filepath"
	"runtime"
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
	FATAL = "FATAL"
)

var (
	g   *log.Logger
	out *os.File
)

func init() {
	logFile := config.Get("log.file")
	logPath := filepath.Dir(logFile)

	// create all parent directory if not exists
	_, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(logPath, 0666)
		if err != nil {
			log.Fatal(err)
		}
	}

	// open log file handle and create log file if not exists
	out, err = os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)

	// postpone action to main method execut finish of the close handler
	// if not call close hook, the file handler will be gc recycle
	//	defer out.Close()

	if err != nil {
		panic(err)
	}

	log.Println("log to:", logFile)
	g = log.New(out, "", log.LstdFlags)
}

func format(pre string, info interface{}) string {
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Base(file)
	return fmt.Sprintf("%v %v:%v %v", pre, file, line, info)
}

func Debug(info interface{}) {
	g.Println(format(DEBUG, info))
}

func Info(info interface{}) {
	g.Println(format(INFO, info))
}

func Warn(info interface{}) {
	g.Println(format(WARN, info))
}

func Error(info interface{}) {
	g.Println(format(ERROR, info))
}

func Fatal(info interface{}) {
	g.Println(format(FATAL, info))
	os.Exit(13)
}

//func Close() {
//	Info("call logger shutting down hook")
//	if f, ok := out.(*os.File); ok {
//		out.Sync()
//		out.Close()
//	} else if f, ok := out.(io.Closer); ok {
//		f.Close()
//	}
//}
