package main

import (
	"mknote/server"
	"mknote/server/ctx"
	"os"
)

func main() {
	var isTls = true

	if _, err := os.Stat(ctx.Get("server.tls.cert.file")); os.IsNotExist(err) {
		isTls = false
	}

	if _, err := os.Stat(ctx.Get("server.tls.key.file")); os.IsNotExist(err) {
		isTls = false
	}

	if isTls {
		server.StartTLS()
	} else {
		server.Start()
	}

}
