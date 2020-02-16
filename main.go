package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/pkg/errors"
)

func main() {
	fmt.Println("github_trending_api_server running...\n:8080\n/contributions?user=[username]\n/trending")
	
	http.HandleFunc("/contributions", contributionAPIHandle)
	http.HandleFunc("/trending", trendingAPIHandle)

	http.HandleFunc("/",defaultHttpRequestHandle)
	http.ListenAndServe(":8080", nil)
}

func defaultHttpRequestHandle(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Hello world.")
}

func contributionAPIHandle(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		return
	}

	user := request.Form.Get("user")
	contributions, _ := getContributions(user)
	result, _ := json.Marshal(contributions)

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(writer, string(result))
}

func trendingAPIHandle(writer http.ResponseWriter, request *http.Request) {
	trending, _ := getTrending()
	result, _ := json.Marshal(trending)

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(writer, string(result))
}

type Trending struct {
	Repositories []Repository
}

type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Star        string `json:"star"`
	Fork        string `json:"fork"`
	Lang        string `json:"lang"`
}

func getTrending() (Trending, error) {
	//no use http/2
	//http.DefaultTransport.(*http.Transport).TLSNextProto = make(map[string]func(string, *tls.Conn) http.RoundTripper)
	requestUrl := "https://github.com/trending"
	response,err := http.Get(requestUrl)
	if err != nil {
		fmt.Println(err)
		return Trending{},err
	}

	contentBytes, err := ioutil.ReadAll(response.Body)
	if err != nil || contentBytes == nil {
		return Trending{}, err
	}

	htmlContent := string(contentBytes)
	repositories := resolveRepositories(htmlContent)
	return Trending{Repositories:repositories}, nil
}

func resolveRepositories(content string) []Repository {
	rep := `<article class="Box-row">[\s\S]*?<\/article>`
	regexp := regexp.MustCompile(rep)
	match := regexp.FindAll([]byte(content), -1)
	if match == nil {
		return nil
	}

	result := make([]Repository, 0)
	for i := 0; i < len(match); i++ {
		result = append(result, resolveRepositoryTag(string(match[i])))
	}

	return result
}

func resolveRepositoryTag(content string) Repository {
	name := stringFormat(getRepositoryName(content))
	lang := stringFormat(getRepositoryLang(content))
	desc := stringFormat(getRepositoryDescription(content))
	star := stringFormat(getRepositoryStar(content))
	fork := stringFormat(getRepositoryFork(content))
	url := "https://github.com" + name

	return Repository{Name: name, Description: desc, Lang: lang, Star: star, Fork: fork, Url: url}
}

func stringFormat(content string) string {
	content = strings.Replace(content, "\n", "", -1)
	content = strings.TrimSpace(content)
	return content
}

func getRepositoryName(content string) string {
	repositoryItemNameExp := `(?<=<h1[\s\S]+<a href=")[\S]+(?=")`
	name := findFirstOrDefaultMatch(content,repositoryItemNameExp)
	return name
}

func getRepositoryDescription(content string) string {
	repositoryItemDescriptionExp := `(?<=<p class="col-9 text-gray my-1 pr-4">)[\s\S]+?(?=<\/p>)`
	desc := findFirstOrDefaultMatch(content,repositoryItemDescriptionExp)
	return desc;
}

func getRepositoryLang(content string) string {
	repositoryItemLangExp := `(?<=<span itemprop="programmingLanguage">)[\s\S]+?(?=<\/span>)`
	return findFirstOrDefaultMatch(content, repositoryItemLangExp)
}

func getRepositoryStar(content string) string {
	repositoryItemStarTagExp := `(?<=<a class="[\s]?muted-link d-inline-block mr-3"[\s\S]+stargazers">[\s\S]+g>)[\s\S]*?(?=<\/a>)`
	starValue := findFirstOrDefaultMatch(content, repositoryItemStarTagExp)
	return starValue
}

func getRepositoryFork(content string) string {
	repositoryItemForkTagExp := `(?<=<a class="[\s]?muted-link d-inline-block mr-3"[\s\S]+network/[\S]*">[\s\S]+g>)[\s\S]*?(?=<\/a>)`
	forkValue := findFirstOrDefaultMatch(content, repositoryItemForkTagExp)
	return forkValue
}

func findFirstOrDefaultMatch(content string, exp string) string {
	regexp2 := regexp2.MustCompile(exp, 0)
	match, err := regexp2.FindStringMatch(content)
	if err != nil || match == nil {
		return ""
	}

	groups := match.Groups()
	if len(groups) > 1 {
		return groups[1].Capture.String()
	}

	return match.String()
}

type Contribution struct {
	DataCount int    `json:"count"`
	Date      string `json:"date"`
	Color     string `json:"color"`
}

func getContributions(userName string) ([]Contribution, error) {
	if userName == "" {
		return nil, errors.New("required:username")
	}

	currentTime := time.Now()
	requestUrl := "https://github.com/users/" + userName + "/contributions?to=" + currentTime.Format("2006-01-02")
	res, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	contentBytes, err := ioutil.ReadAll(res.Body)
	if err != nil || contentBytes == nil {
		return nil, err
	}

	contentString := string(contentBytes)
	return resolveContributions(contentString), nil
}

func resolveContributions(content string) []Contribution {
	rectTags := resolveRectTags(content)
	if len(rectTags) == 0 {
		return nil
	}

	contributions := make([]Contribution, 0)
	for i := 0; i < len(rectTags); i++ {
		contribution := createContributionByRectTag(rectTags[i])
		contributions = append(contributions, contribution)
	}

	return contributions
}

func resolveRectTags(content string) []string {
	exp := `<rect.*?\/>`
	regexp := regexp.MustCompile(exp)
	match := regexp.FindAll([]byte(content), -1)
	if match == nil {
		return nil
	}

	result := make([]string, 0)
	for i := 0; i < len(match); i++ {
		result = append(result, string(match[i]))
	}

	return result
}

func createContributionByRectTag(tag string) Contribution {
	exp := `(?<=<rect.*data-count=").*(?="\s*data-date.*\/>)|(?<=<rect.*\s*data-date=").*(?="\s?.*\/>)|(?<=<rect.*fill=").*(?="\s*data-count.*\/>)`
	regexp2 := regexp2.MustCompile(exp, 0)
	colorMatchResult, err := regexp2.FindStringMatch(tag)
	if err != nil || colorMatchResult == nil {
		return Contribution{}
	}

	dataCountMatch, _ := regexp2.FindNextMatch(colorMatchResult)
	if dataCountMatch == nil {
		return Contribution{}
	}

	dateMatchResult, _ := regexp2.FindNextMatch(dataCountMatch)
	if dateMatchResult == nil {
		return Contribution{}
	}

	dataCount, _ := strconv.Atoi(dataCountMatch.String())
	date := dateMatchResult.String()
	color := colorMatchResult.String()

	return Contribution{Color: color, DataCount: dataCount, Date: date}
}
