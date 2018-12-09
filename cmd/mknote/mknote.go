package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"github.com/sycki/mknote/cmd/mknote/options"
	"github.com/sycki/mknote/controller"
	"github.com/sycki/mknote/logger"
	"github.com/sycki/mknote/server"
	"github.com/sycki/mknote/storage"
)

var version string

func main() {
	cmd := flag.CommandLine
	config := options.NewDefaultConfig()
	config.AddFlags(cmd)
	flag.Parse()

	if config.Version {
		println(version)
		os.Exit(0)
	}

	logger.SetLevel(config.LogLevel)

	// create context for all threads
	errCh := make(chan error, 1)
	// create storage layer manager and start file cache manager
	sm := storage.NewManager(config)
	sm.Start(errCh)
	defer sm.Stop()

	// create controller manager of page and rest api
	ctr := controller.NewManager(config, sm)
	defer ctr.Close()

	// start http server
	s := server.NewServer(config, ctr)
	s.Start(errCh)
	defer s.Stop()
	logger.Info("mknote is started")

	// start listen system signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		logger.Error("occur error when start mknote:", err)
	case signal := <-sig:
		logger.Warn("receive a signal:", signal)
	}

	close(errCh)
	close(sig)

	logger.Info("mknote stopped gracefully")
}
