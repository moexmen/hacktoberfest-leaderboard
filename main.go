package main

import (
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type config struct {
	GithubToken string `env:"GHTOKEN"`
	Authors     string `env:"AUTHORS"`
}

type UserData struct {
	PrCount   int
	AvatarURL string
}

type avatarResult struct {
	AvatarURL string `json:"avatar_url"`
}

type prCountResult struct {
	PrCount int `json:"total_count"`
}

func leaderboard(writer http.ResponseWriter, request *http.Request) {
	t := template.Must(template.ParseFiles("leaderboard.html"))
	userData := leaderboard_data()
	t.Execute(writer, userData)
}

func leaderboard_json(writer http.ResponseWriter, request *http.Request) {
	jsonString, _ := json.Marshal(leaderboard_data())
	fmt.Fprintf(writer, "%s", jsonString)
}

// return map of usernames to UserData struct containing PR count and avatar URL.
func leaderboard_data() map[string]UserData {
	userData := make(map[string]UserData)
	authors := strings.Split(cfg.Authors, ":")
	fmt.Printf("Authors: %v\n", authors)
	for _, author := range authors {
		userData[author] = UserData{PrCount: getPrCount(author), AvatarURL: getAvatar(author)}
		fmt.Printf("Author: %s, PR count: %d\n", author, userData[author].PrCount)
	}
	return userData
}

func getAvatar(author string) string {
	url := fmt.Sprintf("https://api.github.com/users/%s", author)
	response, err := http.Get(url)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		fmt.Println("Failed to fetch avatar. %s\n", err)
		return ""
	} else {
		ghData, _ := ioutil.ReadAll(response.Body)
		result := avatarResult{}
		json.Unmarshal([]byte(ghData), &result)
		return result.AvatarURL
	}
}

func getPrCount(author string) (prCount int) {
	year := calcYear()
	url := fmt.Sprintf("https://api.github.com/search/issues?q=created:%d-09-30T00:00:00-12:00..%d-10-31T23:59:59-12:00+type:pr+is:public+author:%s", year, year, author)
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	if cfg.GithubToken != "" {
		request.Header.Set("Authorization", "token "+cfg.GithubToken)
	}
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		fmt.Println("Failed to fetch PR count. %s\n", err)
		return -1
	} else {
		ghData, _ := ioutil.ReadAll(response.Body)
		result := prCountResult{}
		json.Unmarshal([]byte(ghData), &result)
		return result.PrCount
	}
}

func calcYear() int {
	currentTime := time.Now()
	dateTimeString := fmt.Sprintf("30 Sep %d 0:00 -1200", currentTime.Year()-2000)
	hacktoberfestStart, _ := time.Parse(time.RFC822Z, dateTimeString)
	if currentTime.Before(hacktoberfestStart) {
		return currentTime.Year() - 1
	} else {
		return currentTime.Year()
	}
}

// global config
var cfg = config{}

func main() {
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	fs := http.FileServer(http.Dir("assets"))

	http.Handle("/", fs)
	http.HandleFunc("/leaderboard.json", leaderboard_json)
	http.HandleFunc("/leaderboard", leaderboard)
	http.ListenAndServe(":4000", nil)
}
