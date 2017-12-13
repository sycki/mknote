package ctx

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	debug  = "DEBUG"
	info   = "INFO"
	warn   = "WARN"
	errors = "ERROR"
	fatal  = "FATAL"
)

var (
	g   *log.Logger
	out *os.File
)

func init() {
	logFile := Get("log.file")
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

func GetLogger() *log.Logger {
	return g
}

func format(pre string, info ...interface{}) string {
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Base(file)
	return fmt.Sprintf("%v %v:%v %v", pre, file, line, strings.Trim(fmt.Sprint(info), "[]"))
}

func Debug(megs ...interface{}) {
	if megs == nil {
		return
	}
	g.Println(format(debug, megs))
}

func Info(megs ...interface{}) {
	if megs == nil {
		return
	}
	g.Println(format(info, megs))
}

func Warn(megs ...interface{}) {
	if megs == nil {
		return
	}
	g.Println(format(warn, megs))
}

func Error(megs ...interface{}) {
	if megs == nil {
		return
	}
	g.Println(format(errors, megs))
}

func Fatal(megs ...interface{}) {
	g.Println(format(fatal, megs))
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
