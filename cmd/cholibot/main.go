package main

import (
	"github.com/heshoots/cholibot/pkg/discord"
	//"github.com/heshoots/cholibot/pkg/webserver"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//go webserver.Start()
	go discord.Start()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
