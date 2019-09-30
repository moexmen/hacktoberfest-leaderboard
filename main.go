package main

import (
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env"
	"html/template"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

type config struct {
	GithubToken     string `env:"GHTOKEN"`
	Authors         string `env:"AUTHORS"`
	RefreshInterval int    `env:"REFRESH_INTERVAL" envDefault:"1800"`
}

type AuthorData struct {
	Author    string
	PrCount   int
	AvatarURL string
}

type LeaderboardData struct {
	AuthorData      []AuthorData
	RefreshInterval int
	Year            int
	UpdatedTime     string
}

type avatarResult struct {
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
}

type prCountResult struct {
	PrCount int `json:"total_count"`
}

func leaderboard(writer http.ResponseWriter, request *http.Request) {
	t := template.Must(template.ParseFiles("leaderboard.html"))
	authorData := getAuthorData()
	leaderboardData := LeaderboardData{AuthorData: authorData, RefreshInterval: cfg.RefreshInterval, Year: calcYear(), UpdatedTime: time.Now().Format("2 Jan 2006 3:04 PM")}
	t.Execute(writer, leaderboardData)
}

func leaderboardJSON(writer http.ResponseWriter, request *http.Request) {
	jsonString, _ := json.Marshal(getAuthorData())
	fmt.Fprintf(writer, "%s", jsonString)
}

// return slice of AuthorData structs ordered by PR count descending
func getAuthorData() []AuthorData {
	authors := strings.Split(cfg.Authors, ":")
	authorData := make([]AuthorData, len(authors))
	fmt.Printf("Authors: %v\n", authors)
	for i, author := range authors {
		avatarData := getAvatar(author)
		currentAuthor := AuthorData{Author: avatarData.Name, PrCount: getPrCount(author), AvatarURL: avatarData.AvatarURL}
		authorData[i] = currentAuthor
		fmt.Printf("Author: %s, PR count: %d\n", currentAuthor.Author, currentAuthor.PrCount)
	}
	sort.Slice(authorData, func(i, j int) bool { return authorData[i].PrCount > authorData[j].PrCount })
	return authorData
}

func getAvatar(author string) avatarResult {
	response, err := makeAuthorizedRequest("https://api.github.com/users/%s", author)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		fmt.Println("Failed to fetch avatar. %s\n", err)
		return avatarResult{}
	} else {
		ghData, _ := ioutil.ReadAll(response.Body)
		result := avatarResult{}
		json.Unmarshal([]byte(ghData), &result)
		return result
	}
}

func getPrCount(author string) (prCount int) {
	year := calcYear()
	response, err := makeAuthorizedRequest("https://api.github.com/search/issues?q=created:%d-09-30T00:00:00-12:00..%d-10-31T23:59:59-12:00+type:pr+is:public+author:%s", year, year, author)
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

func makeAuthorizedRequest(urlFormat string, arguments ...interface{}) (*http.Response, error) {
	url := fmt.Sprintf(urlFormat, arguments...)
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	if cfg.GithubToken != "" {
		request.Header.Set("Authorization", "token "+cfg.GithubToken)
	}
	return client.Do(request)
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
	http.HandleFunc("/leaderboard.json", leaderboardJSON)
	http.HandleFunc("/leaderboard", leaderboard)
	http.ListenAndServe(":4000", nil)
}
