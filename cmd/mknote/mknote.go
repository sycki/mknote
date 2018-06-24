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

package main

import (
	"github.com/sycki/mknote/server"
	"github.com/sycki/mknote/logger"
	"os"
	"context"
	"time"
	"os/signal"
	"syscall"
	"github.com/sycki/mknote/server/persistent"
	"flag"
	"github.com/sycki/mknote/cmd/mknote/options"
	"github.com/sycki/mknote/server/view"
)

var version = "mknote-v2.3.0"

func main() {
	cmd := flag.CommandLine
	config := options.GetDefaultConfig()
	config.AddFlags(cmd)
	flag.Parse()

	if config.Version {
		println(version)
		os.Exit(0)
	}

	var isTls = true
	if _, err := os.Stat(config.TlsCertFile); os.IsNotExist(err) {
		isTls = false
	}
	if _, err := os.Stat(config.TlsKeyFile); os.IsNotExist(err) {
		isTls = false
	}

	// create context
	c, cancel := context.WithTimeout(context.Background(), time.Second)

	// init logger
	logger.Init(config)

	// init view
	view.Init(config)

	// start fs monitor
	persistent.Start(config, c)

	// start http server
	if isTls {
		server.StartTLS(config, c)
	} else {
		server.Start(config, c)
	}

	// start listen stop signal from system
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	cancel()
	close(sig)

	logger.Info("mknote gracefully stopped")
}
