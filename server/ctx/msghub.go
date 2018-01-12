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
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	l               = &sync.Mutex{}
	stoperList      = []chan string{}
	eventListenList = []chan string{}
)

func init() {
	go monitorSignal()
}

func RegistryStoper(msg chan string) {
	l.Lock()
	stoperList = append(stoperList, msg)
	l.Unlock()
}

func StopAll(msg string) {
	for _, c := range stoperList {
		c <- msg
	}
}

func monitorSignal() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	signal.Ignore(syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	close(sigs)
	Warn("recived signal:", sig)
	StopAll("stop")
}

func Registry(msg chan string) {
	l.Lock()
	eventListenList = append(eventListenList, msg)
	l.Unlock()
}

func Notify(msg string) {
	for _, c := range eventListenList {
		c <- msg
	}
}
