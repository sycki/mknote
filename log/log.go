package log

import (
	"fmt"
	"log"
	"os"
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
	logFile := config.GetOr("LOG_FILE", "/var/log/mknode/mknode.log")
	f, err := os.OpenFile(logFile, os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		panic(err)
	}
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
	os.Exit(3)
}
