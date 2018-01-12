/*
Copyright 2017 sycki.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
