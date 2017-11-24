package main

import (
	"mknote/log"
	"mknote/server"
)

func main() {
	log.Info("start mknode server...")
	server.StartServer()
}
