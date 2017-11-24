package main

import (
	"sycki/log"
	"sycki/server"
)

func main() {
	log.Info("start mknode server...")
	server.StartServer()
}
