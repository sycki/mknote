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
	level int
	g     *log.Logger
)

func init() {
	level = Config.LogLevel
	logFile := Config.LogFile
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
	out, err1 := os.OpenFile(logFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)

	// postpone action to main method execut finish of the close handler
	// if not call close hook, the file handler will be gc recycle
	//	defer out.Close()

	if err1 != nil {
		panic(err1)
	}

	log.Println("log to:", logFile)
	g = log.New(out, "", log.LstdFlags)
}

func GetLogger() *log.Logger {
	return g
}

func format(pre string, megs ...interface{}) string {
	if megs == nil {
		return ""
	}
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Base(file)
	return fmt.Sprintf("%v %v:%v %v\n", pre, file, line, strings.Trim(fmt.Sprint(megs), "[]"))
}

func Debug(megs ...interface{}) {
	if level <= 0 {
		g.Print(format(debug, megs))
	}
}

func Info(megs ...interface{}) {
	if level <= 1 {
		g.Print(format(info, megs))
	}
}

func Warn(megs ...interface{}) {
	if level <= 2 {
		g.Print(format(warn, megs))
	}
}

func Error(megs ...interface{}) {
	if level <= 3 {
		g.Print(format(errors, megs))
	}
}

func Fatal(megs ...interface{}) {
	if level <= 4 {
		g.Print(format(fatal, megs))
		os.Exit(13)
	}
}
