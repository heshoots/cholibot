package challonge

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/heshoots/cholibot/pkg/models"
	"github.com/heshoots/dmux"
	"net/http"
)

func GetCredentials(m dmux.MessageContext) (string, string, error) {
	db := models.GetClient()
	if exists, err := db.HasChallonge(m.Guild().ID()); !exists || err != nil {
		return "", "", errors.New("Guild does not have challonge credentials")
	}
	return db.GetChallonge(m.Guild().ID())
}

func CreateTournament(apikey string, subdomain string, name string, game string) (string, error) {
	client := &http.Client{}
	tournamentvalues := map[string]string{"name": name, "url": name, "subdomain": subdomain, "game_name": game, "tournament_type": "double elimination"}
	values := map[string]map[string]string{"tournament": tournamentvalues}
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest("POST", "https://api.challonge.com/v1/tournaments.json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("api_key", apikey)
	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return "", errors.New(resp.Status + "challonge create failed " + buf.String())
	}
	return "http://" + subdomain + ".challonge.com/" + name, nil
}

/*
func ChallongeHandler(apikey string, subdomain string, postto []string, errorto []string) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		command := discord.GetCommand(m)
		split := strings.SplitAfterN(command, " ", 2)
		if len(split) != 2 {
			for _, channel := range errorto {
				s.ChannelMessageSend(channel, "not enough input, command: !challonge url game_name")
			}
			return
		}
		name := strings.Trim(split[0], " ")
		game := strings.Trim(split[1], " ")
		url, err := CreateTournament(apikey, subdomain, name, game)
		if err != nil {
			log.Panic(err)
			for _, channel := range errorto {
				s.ChannelMessageSend(channel, "couldn't create tournament: "+err.Error())
			}
			return
		}
		for _, channel := range postto {
			s.ChannelMessageSend(channel, url)
		}
	}
}*/
