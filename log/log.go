package log

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sycki/config"
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
	FATAL = "FATAL"
)

var g *log.Logger

func init() {
	logFile := config.Get("LOG_FILE")
	logPath := filepath.Dir(logFile)

	// create parent directory if not exists
	_, e := os.Stat(logPath)
	if os.IsNotExist(e) {
		e := os.MkdirAll(logPath, 0666)
		if e != nil {
			log.Fatal(e)
		}
	}

	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	log.Println("log to:", logFile)
	g = log.New(f, "", log.LstdFlags)
}

func format(pre string, info interface{}) string {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%v [%v:%v], msg: [%v]", pre, file, line, info)
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
