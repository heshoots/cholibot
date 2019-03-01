package webserver

import (
	"encoding/json"
	"github.com/heshoots/cholibot/pkg/models"
	"github.com/jinzhu/configor"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var discordRedirect string

var Config = struct {
	WebServer struct {
		Port         string `default:"3000"`
		ClientID     string `required:"true"`
		ClientSecret string `required:"true"`
		RedirectURL  string `required:"true"`
		Scope        string `required:"true"`
	}
}{}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type GuildRole struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Color       int64  `json:"color"`
	Hoist       bool   `json:"hoist"`
	Position    int64  `json:"position"`
	Permissions int64  `json:"permissions"`
	Managed     bool   `json:"managed"`
	Mentionable bool   `json:"mentionable"`
}

type GuildRoles []GuildRole

type UserResponse struct {
	ID string `json:"id"`
}

type UserGuildResponse []UserGuild

type UserGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Owner       bool   `json:"owner"`
	Permissions int64  `json:"permissions"`
	BotJoined   bool
}

func init() {
	environment, envSet := os.LookupEnv("ENV")
	if !envSet {
		environment = "development"
	}
	err := configor.Load(&Config, "./config/"+environment+".yaml")
	if err != nil {
		log.Fatal(err)
	}
	var escapedRedirect string = url.QueryEscape(Config.WebServer.RedirectURL)
	var scope = url.QueryEscape(Config.WebServer.Scope)
	discordRedirect = "https://discordapp.com/api/oauth2/authorize?client_id=" + Config.WebServer.ClientID + "&redirect_uri=" + escapedRedirect + "&response_type=code&scope=" + scope

}

func login(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, discordRedirect, 307)
}

func createData(api_key string) url.Values {
	data := url.Values{}
	data.Set("client_id", Config.WebServer.ClientID)
	data.Set("client_secret", Config.WebServer.ClientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", api_key)
	data.Set("redirect_uri", Config.WebServer.RedirectURL)
	data.Set("scope", "identify email connections")
	return data
}

func getUrl() string {
	apiUrl := "https://discordapp.com"
	resource := "/api/oauth2/token"
	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()
	return urlStr
}

func genericDataResponse(url, api_key string, data interface{}) error {
	client := &http.Client{}
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	r.Header.Add("Authorization", "Bearer "+api_key)
	resp, _ := client.Do(r)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	return nil
}

func readTokenResponse(resp *http.Response) TokenResponse {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	var jsondata TokenResponse // TopTracks
	err = json.Unmarshal(body, &jsondata)
	if err != nil {
		panic(err.Error())
	}
	return jsondata
}

func callback(w http.ResponseWriter, r *http.Request) {
	api_key := r.URL.Query()["code"]
	if api_key != nil && len(api_key) == 0 {
		log.Fatal("No api_key returned")
		return
	}
	data := createData(api_key[0])
	urlStr := getUrl()

	client := &http.Client{}
	r, _ = http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(r)
	jsondata := readTokenResponse(resp)
	cookie := http.Cookie{
		Name:  "session",
		Value: jsondata.AccessToken,
		//Domain:   "localhost:3000",
		Path:     "/",
		MaxAge:   60 * 60,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", 301)
}

func getGuildRoles(guildID string, api_key string) (m GuildRoles) {
	genericDataResponse("https://discordapp.com/api/guilds/"+guildID+"/roles", api_key, &m)
	return m
}

func getUser(api_key string) (m UserResponse) {
	genericDataResponse("https://discordapp.com/api/users/@me", api_key, &m)
	return m
}

func getUsersGuilds(api_key string) (m UserGuildResponse) {
	genericDataResponse("https://discordapp.com/api/users/@me/guilds", api_key, &m)
	return m
}

func getGuild(guildID string, api_key string) (m UserGuild) {
	genericDataResponse("https://discordapp.com/api/guilds/guildID", api_key, &m)
	return m
}

type HomeData struct {
	KnownGuilds []models.DiscordGuild
	UserGuilds  UserGuildResponse
}

func Start() {
	router := getRouter()
	log.Fatal(http.ListenAndServe(":"+Config.WebServer.Port, router))
}
