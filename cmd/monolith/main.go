package main

import (
	"github.com/quorauk/cholibot/pkg/discord"
	"github.com/quorauk/cholibot/pkg/webserver"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var version string

func main() {
	log.Info("Version: " + version)
	go webserver.Start()
	go discord.Start()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	log.Info("Bot is now running.  Press CTRL-C to exit.")
	<-sc
}
