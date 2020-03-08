package discord

import (
	"github.com/quorauk/dmux"
	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/configor"
	"github.com/quorauk/cholibot/pkg/models"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

var authToken string
var discordInstance dmux.Session

var Config = struct {
	DiscordBot struct {
		AuthToken   string `required:"true"`
		TestAccount string `required:"false"`
		TestGuild   string `required:"false"`
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

func GetOrCreateDataChannel(s *discordgo.Session) *discordgo.Channel {
	channels, _ := s.GuildChannels("386488116426440704")
	for _, channel := range channels {
		if (channel.Name == "choli-data") {
			return channel
		}
	}
	channel, _ := s.GuildChannelCreate(
		"386488116426440704",
		"choli-data",
		discordgo.ChannelTypeGuildText,
	)
	return channel
}

func GetLastChannelMessage(s *discordgo.Session, c *discordgo.Channel) string {
	lastmessage, err := s.ChannelMessages(c.ID, int(1), "", "", "")
	log.Info(err)
	log.Info(lastmessage)
	if len(lastmessage) > 0 {
		return lastmessage[0].Content
	}
	return ""
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

	discordDb := models.GetDiscordClient(discordInstance)
	roles, err := discordDb.GetRoles("386488116426440704")
	log.Info(roles)
	log.Info(err)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discordInstance.Close()
}
