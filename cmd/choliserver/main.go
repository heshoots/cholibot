package main

import (
	"github.com/quorauk/cholibot/pkg/webserver"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

const version = "v0.0.1"

func main() {
	log.Info("Version: " + version)
	go webserver.Start()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	log.Info("Server is now running.  Press CTRL-C to exit.")
	<-sc
}
