package main

import (
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env"
	"io/ioutil"
	"net/http"
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
	userData := make(map[string]UserData)
	authors := []string{"teo-shaowei", "fonglh", "cflee", "naomilwx", "jchiam", "jeremyyap", "wongherlung", "xeluna", "weelillad"}
	for _, author := range authors {
		userData[author] = UserData{PrCount: getPrCount(author), AvatarURL: getAvatar(author)}
	}
	jsonString, _ := json.Marshal(userData)
	fmt.Fprintf(writer, "%s", jsonString)
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

func getPrCount(author string) int {
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
		fmt.Println("The HTTP request failed with error %s\n", err)
		return 0
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

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	fmt.Printf("%+v\n", cfg)

	http.HandleFunc("/leaderboard", leaderboard)
	http.ListenAndServe(":4000", nil)
}
