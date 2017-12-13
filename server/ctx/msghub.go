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
