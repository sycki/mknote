package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	levelText = map[int]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
		FATAL: "FATAL",
	}
	level = INFO
	out   = log.New(os.Stdout, "", log.LstdFlags)
)

func GetLevel(l int) {level = l}
func SetLevel(l int) {level = l}
func GetLogger() *log.Logger {return out}
func SetLogger(o *log.Logger) {out = o}

func output(str string) {
	out.Print(str)
}

func format(level int, megs ...interface{}) string {
	if megs == nil {
		return ""
	}
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Base(file)
	return fmt.Sprintf("%s %s:%d %s", levelText[level], file, line, fmt.Sprintln(megs...))
}

func Debug(megs ...interface{}) {
	if level <= DEBUG {
		output(format(DEBUG, megs...))
	}
}

func Debugf(f string, megs ...interface{}) {
	if level <= DEBUG {
		output(format(DEBUG, fmt.Sprintf(f, megs...)))
	}
}

func Info(megs ...interface{}) {
	if level <= INFO {
		output(format(INFO, megs...))
	}
}

func Infof(f string, megs ...interface{}) {
	if level <= INFO {
		output(format(INFO, fmt.Sprintf(f, megs...)))
	}
}

func Warn(megs ...interface{}) {
	if level <= WARN {
		output(format(WARN, megs...))
	}
}

func Warnf(f string, megs ...interface{}) {
	if level <= WARN {
		output(format(WARN, fmt.Sprintf(f, megs...)))
	}
}

func Error(megs ...interface{}) {
	if level <= ERROR {
		output(format(ERROR, megs...))
	}
}

func Errorf(f string, megs ...interface{}) {
	if level <= ERROR {
		output(format(ERROR, fmt.Sprintf(f, megs...)))
	}
}

func Fatal(megs ...interface{}) {
	if level <= FATAL {
		output(format(FATAL, megs...))
		os.Exit(FATAL)
	}
}

func Fatalf(f string, megs ...interface{}) {
	if level <= FATAL {
		output(format(FATAL, fmt.Sprintf(f, megs...)))
		os.Exit(FATAL)
	}
}
