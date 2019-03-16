package models

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/configor"
	log "github.com/sirupsen/logrus"
	"os"
)

var client *ModelClient

var Config = struct {
	Redis struct {
		Address  string `required:"true"`
		Password string `required:"false"`
		DB       int    `required:"false"`
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
}

func GetClient() *ModelClient {
	Environment()
	if client == nil {
		client = &ModelClient{redis.NewClient(&redis.Options{
			Addr:     Config.Redis.Address,
			Password: Config.Redis.Password,
			DB:       Config.Redis.DB,
		})}
		return client
	} else {
		return client
	}
}

func (m *ModelClient) Pong() {
	pong, err := m.Ping().Result()
	fmt.Println(pong, err)
}

func (m *ModelClient) GetRoles(guild string) ([]string, error) {
	pong, err := m.HKeys("guild:" + guild).Result()
	if err != nil {
		return nil, err
	}
	return pong, nil
}

func (m *ModelClient) GetRole(guild string, role string) (string, error) {
	return m.HGet("guild:"+guild, role).Result()
}

func (m *ModelClient) SetRole(guild string, role string, id string) error {
	res := m.HSet("guild:"+guild, role, id)
	return res.Err()
}

func (m *ModelClient) RemoveRole(guild string, role string) error {
	res := m.HDel("guild:"+guild, role)
	return res.Err()
}

func (m *ModelClient) GetCustomCommands(guild string) ([]string, error) {
	return m.HKeys("commands:" + guild).Result()
}

func (m *ModelClient) SetCustomCommand(guild string, command string, response string) error {
	res := m.HSet("commands:"+guild, command, response)
	return res.Err()
}

func (m *ModelClient) GetCustomCommand(guild string, command string) (string, error) {
	return m.HGet("commands:"+guild, command).Result()
}

func (m *ModelClient) HasChallonge(guild string) (bool, error) {
	val, err := m.Exists("challonge:apikey:" + guild).Result()
	if err != nil {
		return false, err
	}
	return val != 0, nil
}

func (m *ModelClient) GetChallonge(guild string) (string, string, error) {
	apikey, err := m.Get("challonge:apikey:" + guild).Result()
	if err != nil {
		return "", "", err
	}
	subdomain, err := m.Get("challonge:subdomain:" + guild).Result()
	if err != nil {
		return "", "", err
	}
	return apikey, subdomain, nil
}

type ModelSource interface {
	GetRoles(guild string) ([]string, error)
	GetRole(guild string, role string) (string, error)
	SetRole(guild string, role string) error
	RemoveRole(guild string, role string) error
	GetCustomCommands(guild string) ([]string, error)
	SetCustomCommand(guild, command, response string) error
	GetCustomCommand(guild, command string) (string, error)
	HasChallonge(guild string) (bool, error)
	GetChallonge(guild string) (string, string, error)
	Pong()
}

type ModelClient struct {
	*redis.Client
}
