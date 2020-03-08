package models

import (
	log "github.com/sirupsen/logrus"
	"github.com/jinzhu/configor"
	"github.com/bwmarrin/discordgo"
	"github.com/quorauk/dmux"
	"encoding/base64"
	"encoding/json"
	"errors"
	"crypto/aes"
	"crypto/cipher"
)

type DiscordDB struct {
	Session *discordgo.Session
	Secret string
}


var discordClient *DiscordDB

var DiscordDBConfig = struct {
	DiscordBot struct {
		Secret string `required:"true"`
	}
}{}

func GetDiscordClient(session dmux.Session) *DiscordDB {
	if discordClient == nil {
		configor.Load(&DiscordDBConfig, "config/development.yaml")
		log.Info(DiscordDBConfig.DiscordBot.Secret)
		discordClient = &DiscordDB{Session: session.RawSession(), Secret: DiscordDBConfig.DiscordBot.Secret}
	}
	return discordClient
}

func getOrCreateDataChannel(s *discordgo.Session) *discordgo.Channel {
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

func getLastChannelMessage(s *discordgo.Session, c *discordgo.Channel) (string, error) {
	lastmessage, _ := s.ChannelMessages(c.ID, int(1), "", "", "")
	if len(lastmessage) > 0 {
		return lastmessage[0].Content, nil
	}
	return "", errors.New("no data set")
}

type DiscordDBState struct {
	Roles map[string]string `json:"roles"`
	Commands map[string]string `json:"commands"`
	Challonge string `json:"challonge"`
}

func (d *DiscordDB) getGCM() cipher.AEAD {
	key := []byte(d.Secret)
	c, err := aes.NewCipher(key)
	if (err != nil) {
		log.Fatal(err)
	}
	gcm, _ := cipher.NewGCM(c)
	return gcm
}

func (d *DiscordDB) encrypt(data []byte) []byte {
	gcm := d.getGCM()
	nonce := make([]byte, gcm.NonceSize())
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return []byte(base64.StdEncoding.EncodeToString(ciphertext))
}

func (d *DiscordDB) decrypt(b64text []byte) ([]byte, error) {
	gcm := d.getGCM()

	ciphertext, _ := base64.StdEncoding.DecodeString(string(b64text))
	nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }

    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    return gcm.Open(nil, nonce, ciphertext, nil)
}

func decodeDataMessage(message string) *DiscordDBState {
	var state *DiscordDBState

	json.Unmarshal([]byte(message), &state)
	return state
}

func defaultDataMessage() *DiscordDBState {
	return &DiscordDBState{
		Roles: make(map[string]string),
		Commands: make(map[string]string),
	}
}

func (d *DiscordDB) setData(guild string, data *DiscordDBState) error {
	channel := getOrCreateDataChannel(d.Session)
	jsonData, _ := json.Marshal(data)
	encryptedData := d.encrypt(jsonData)
	d.Session.ChannelMessageSend(channel.ID, string(encryptedData))
	return nil
}

func (d *DiscordDB) GetGuildState(guild string) *DiscordDBState {
	channel := getOrCreateDataChannel(d.Session)
	lastMessage, err := getLastChannelMessage(d.Session, channel)
	if err != nil {
		defaultMessage := defaultDataMessage()
		d.setData(guild, defaultMessage)
		return defaultMessage
	}
	plaintext, err := d.decrypt([]byte(lastMessage))
	log.Info(plaintext)
	log.Info(err)
	return decodeDataMessage(string(plaintext))
}

func (d *DiscordDB) GetRoles(guild string) ([]string, error) {
	state := d.GetGuildState(guild)
	log.Info(state)
	keys := make([]string, 0, len(state.Roles))
	for key := range state.Roles {
		keys = append(keys, key)
	}

	return keys, nil
}

func (d *DiscordDB) GetRole(guild string, role string) (string, error) {
	state := d.GetGuildState(guild)

	return state.Roles[role], nil
}

func (d *DiscordDB) SetRole(guild string, role string, id string) error {
	state := d.GetGuildState(guild)
	state.Roles[role] = id
	d.setData(guild, state)
	return nil
}

func (d *DiscordDB) RemoveRole(guild string, toremove string) error {
	state := d.GetGuildState(guild)
	for index, role := range state.Roles {
		if role == toremove {
			delete(state.Roles, index)
			d.setData(guild, state)
		}
	}
	return nil
}

func (d *DiscordDB) SetCustomCommand(guild, command, response string) error {
	state := d.GetGuildState(guild)

	state.Commands[command] = response
	d.setData(guild, state)
	return nil
}

func (d *DiscordDB) GetCustomCommand(guild, command string) (string, error) {
	state := d.GetGuildState(guild)

	return state.Commands[command], nil
}

func (d *DiscordDB) GetCustomCommands(guild string) ([]string, error) {
	state := d.GetGuildState(guild)
	keys := make([]string, 0, len(state.Commands))
	for key := range state.Commands {
		keys = append(keys, key)
	}

	return keys, nil
}