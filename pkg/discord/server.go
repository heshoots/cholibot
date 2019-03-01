package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/heshoots/cholibot/pkg/models"
	"github.com/jinzhu/configor"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

var authToken string
var discordInstance *discordgo.Session

func GetRolesForGuild(guild string) ([]*discordgo.Role, error) {
	return discordInstance.GuildRoles(guild)
}

var Config = struct {
	DiscordBot struct {
		AuthToken string `required:"true"`
	}
}{}

func init() {
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

func joinedServer(s *discordgo.Session, g *discordgo.GuildCreate) {
	var newGuild = models.DiscordGuild{}
	db := models.GetDB()
	db.FirstOrCreate(&newGuild, &models.DiscordGuild{GuildID: g.ID})
}

func leftServer(s *discordgo.Session, g *discordgo.GuildDelete) {
	if !g.Unavailable {
		log.Println("Removing guild (have been removed from the server)")
		var guild = models.DiscordGuild{GuildID: g.ID}
		db := models.GetDB()
		db.First(&guild)
		db.Delete(&guild)
	}
}

func addRole(s *discordgo.Session, g *discordgo.GuildRoleCreate) {
	log.Println("Adding new guild role")
	db := models.GetDB()
	var guild = models.DiscordGuild{GuildID: g.GuildID}
	db.First(&guild)
	var role = &models.DiscordRole{
		Name:    g.Role.Name,
		RoleID:  g.Role.ID,
		GuildID: g.GuildID,
	}
	db.Create(&role)
}

func GetGuild(guildID string) (*discordgo.Guild, error) {
	return discordInstance.Guild(guildID)
}

func hasPattern(pattern string, m *discordgo.MessageCreate) bool {
	regex := regexp.MustCompile(`^!showroles$`)
	return regex.MatchString(m.Content)
}

func ShowRoles(s *discordgo.Session, m *discordgo.MessageCreate) {
	if hasPattern(`!showroles$`, m) {
		db := models.GetDB()
		var roles []models.DiscordRole
		db.Where("guild_id = ?", m.GuildID).Find(&roles)
		out := "```Available roles\n-------------\n"
		guildRoles, err := s.GuildRoles(m.GuildID)
		if err != nil {
			log.Println(err)
			return
		}
		for _, grole := range guildRoles {
			for _, role := range roles {
				if grole.ID == role.RoleID {
					out += "!iam " + grole.Name + "\n"
				}
			}
		}
		out += "```"
		_, err = s.ChannelMessageSend(m.ChannelID, out)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func Start() {
	var err error
	discordInstance, err = discordgo.New("Bot " + authToken)
	if err != nil {
		panic(err)
		return
	}

	discordInstance.AddHandler(joinedServer)
	discordInstance.AddHandler(leftServer)
	discordInstance.AddHandler(addRole)
	discordInstance.AddHandler(ShowRoles)
	discordInstance.Open()
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discordInstance.Close()
}
