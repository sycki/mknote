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

package logger

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
	level = 1
	out   = log.New(os.Stdout, "", log.LstdFlags)
)

func SetLevel(l int) {
	level = l
}

func GetLogger() *log.Logger {
	return out
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
		out.Print(format(debug, megs))
	}
}

func Info(megs ...interface{}) {
	if level <= 1 {
		out.Print(format(info, megs))
	}
}

func Warn(megs ...interface{}) {
	if level <= 2 {
		out.Print(format(warn, megs))
	}
}

func Error(megs ...interface{}) {
	if level <= 3 {
		out.Print(format(errors, megs))
	}
}

func Fatal(megs ...interface{}) {
	if level <= 4 {
		out.Print(format(fatal, megs))
		os.Exit(13)
	}
}
