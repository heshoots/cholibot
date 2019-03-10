package discord

import (
	"fmt"
	"github.com/heshoots/dmux"
	"github.com/jinzhu/configor"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var authToken string
var discordInstance dmux.Session

var Config = struct {
	DiscordBot struct {
		AuthToken string `required:"true"`
	}
}{}

func Environment() {
	environment, envSet := os.LookupEnv("ENV")
	if !envSet {
		environment = "development"
	}
	err := configor.Load(&Config, "./config/"+environment+".yaml")
	if err != nil {
		log.Fatal(err)
	}
	authToken = Config.DiscordBot.AuthToken
}

func Start() {
	Environment()

	discordInstance, err := dmux.Router(authToken)
	if err != nil {
		panic(err)
		return
	}
	for _, handler := range Handlers() {
		discordInstance.AddHandler(handler)
	}
	discordInstance.Open()
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discordInstance.Close()
}
