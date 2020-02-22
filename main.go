package main

import (
	"encoding/json"
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var githubUrl = "https://github.com"

func main() {
	listenPort := ":8080"
	contributionRouteString := "/contributions"
	trendingRouteString := "/trending/"
	developerTrendingRouteString := "/trending/developers/"

	fmt.Println("github_trending_api_server running...\n" +
		listenPort + "\n" +
		contributionRouteString + "?user=[username]\n" +
		trendingRouteString + "\n" +
		developerTrendingRouteString)

	http.HandleFunc(contributionRouteString, contributionAPIHandle)
	http.HandleFunc(trendingRouteString, trendingAPIHandle)
	http.HandleFunc(developerTrendingRouteString,trendingDeveloperAPIHandle)

	http.HandleFunc("/",defaultHttpRequestHandle)
	http.ListenAndServe(listenPort, nil)
}

func defaultHttpRequestHandle(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Hello world.")
}

func trendingDeveloperAPIHandle(writer http.ResponseWriter,r *http.Request){
	developerTrend, _ := getDeveloperTrending(r)
	ok(writer,developerTrend)
}

func contributionAPIHandle(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		return
	}

	user := request.Form.Get("user")
	contributions, _ := getContributions(user)
	ok(writer,contributions)
}

func trendingAPIHandle(writer http.ResponseWriter, request *http.Request) {
	trending, _ := getTrending(request)
	ok(writer,trending)
}

func ok(response http.ResponseWriter,data interface{}){
	result := ApiResponse{
		Code: 200,
		Data: data,
		Date: time.Now(),
	}

	jsonValue, _ := json.Marshal(result)
	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(response, string(jsonValue))
}

// 获取项目排行榜
func getTrending(request *http.Request) ([]Repository, error) {
	//no use http/2
	//http.DefaultTransport.(*http.Transport).TLSNextProto = make(map[string]func(string, *tls.Conn) http.RoundTripper)
	githubRequestUrl := getGithubRequestPathAndQueryString(githubUrl,request)
	response,err := http.Get(githubRequestUrl)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}

	contentBytes, err := ioutil.ReadAll(response.Body)
	if err != nil || contentBytes == nil {
		return nil, err
	}

	repositories := resolveRepositories(string(contentBytes))
	return repositories, nil
}

// 获取开发者排行榜
func getDeveloperTrending(request *http.Request) ([]DeveloperTrendDataItem,error) {
	githubRequestUrl := getGithubRequestPathAndQueryString(githubUrl,request)

		res,error := http.Get(githubRequestUrl)
		if error != nil{
			return nil,error
		}

		contentBytes, err := ioutil.ReadAll(res.Body)
		if err != nil || contentBytes == nil {
			return nil, err
		}


	return resolveDeveloperTrending(string(contentBytes)),nil
}

// 获取指定用户活跃榜
func getContributions(userName string) ([]Contribution, error) {
	if userName == "" {
		return nil, errors.New("required:username")
	}

	currentTime := time.Now()
	requestUrl := githubUrl+"/users/" + userName + "/contributions?to=" + currentTime.Format("2006-01-02")
	res, err := http.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	contentBytes, err := ioutil.ReadAll(res.Body)
	if err != nil || contentBytes == nil {
		return nil, err
	}

	return resolveContributions(string(contentBytes)), nil
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
	starToday := stringFormat(getRepositoryTodayStar(content))
	fork := stringFormat(getRepositoryFork(content))
	url := githubUrl + name

	return Repository{Name: name, Description: desc, Lang: lang, Star: star,StarToday:starToday,Fork: fork, Url: url}
}

func getRepositoryName(content string) string {
	repositoryItemNameExp := `(?<=<h1[\s\S]+<a href=")[\S]+(?=")`
	name := findFirstOrDefaultMatch(content,repositoryItemNameExp)
	return name
}

func getRepositoryDescription(content string) string {
	repositoryItemDescriptionExp := `(?<=<p class="col-9 text-gray my-1 pr-4">)[\s\S]+?(?=<\/p>)`
	desc := findFirstOrDefaultMatch(content,repositoryItemDescriptionExp)
	return desc
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

func getRepositoryTodayStar(content string) string {
	repositoryItemStarTagExp := `(?<=<span class="[\s]?d-inline-block float-sm-right">[\s\S]+g>)[\s\S]*?(?=stars)`
	starValue := findFirstOrDefaultMatch(content, repositoryItemStarTagExp)
	return starValue
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

func resolveContributions(content string) []Contribution {
	rectTags := resolveRectTags(content)
	if len(rectTags) == 0 {
		return nil
	}

	res := make([]Contribution,0)
	for i := 0; i < len(rectTags); i++ {
		contributionData := createContributionByRectTag(rectTags[i])
		res = append(res,contributionData)
	}

	for left, right := 0, len(res)-1; left < right; left, right = left+1, right-1 {
		res[left], res[right] = res[right], res[left]
	}

	return res
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

func createContributionByRectTag(tag string) Contribution{
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
	date,_ := time.Parse("2006-01-02",dateMatchResult.String())

	contributionData := Contribution{
		OfficialColor: colorMatchResult.String(),
		Date:          date,
		Total:         dataCount,
	}

	contributionData.Weekday = int(date.Weekday())
	contributionData.Month = date.Month().String()
	contributionData.Year = date.Year()

	switch contributionData.OfficialColor {
	case "#ebedf0":
		contributionData.Level = 0
	case "#c6e48b":
		contributionData.Level = 1
	case "#7bc96f":
		contributionData.Level = 2
	case "#239a3b":
		contributionData.Level = 3
	case "#196127":
		contributionData.Level = 4
	}

	return contributionData
}

func resolveDeveloperTrending(content string) []DeveloperTrendDataItem {
	rep := `<article class="Box-row d-flex"[\s\S]+?>[\s\S\n]+?<\/article>`
	regexp := regexp.MustCompile(rep)
	match := regexp.FindAll([]byte(content), -1)
	if match == nil {
		return nil
	}

	result := make([]DeveloperTrendDataItem, 0)
	for i := 0; i < len(match); i++ {
		result = append(result, resolveDeveloperTrendDataItem(string(match[i])))
	}

	return result
}

func resolveDeveloperTrendDataItem(content string) DeveloperTrendDataItem {
	developIndexExp := `(?<=<a class="text-gray f6 text-center"[\s\S]+?>)[\s\S]+?(?=<\/a>)`
	avatarExp := `(?<=<img[\s\S]+src=")[\S]+(?=")`
	userNameExp := `(?<=<h1 class="h3 lh-condensed">[\s\S]+>)[\s\S]+?(?=<\/a>[\s\S]+<\/h1>)`
	nickNameExp := `(?<=<p class="f4 text-normal mb-1">[\s\S]+>)[\s\S]+?(?=<\/a>[\s\S]+<\/p>)`

	user := User{
		Name:     stringFormat(findFirstOrDefaultMatch(content,userNameExp)),
		NickName: stringFormat(findFirstOrDefaultMatch(content,nickNameExp)),
		Avatar:   stringFormat(findFirstOrDefaultMatch(content,avatarExp)),
	}

	user.Website = githubUrl + "/" + user.NickName

	repositoryNameExp := `(?<=<article>[\s\S]+h1[\s\S]+\/svg>)[\s\S]+?(?=<)`
	repositoryUrlExp := `(?<=<article>[\s\S]+href=")[\s\S]+?(?=")`
	repositoryDescriptionExp := `(?<=<div class="f6 text-gray mt-1">)[\s\S]+?(?=<)`

	repo := Repository{
		Name:        stringFormat(findFirstOrDefaultMatch(content,repositoryNameExp)),
		Description: stringFormat(findFirstOrDefaultMatch(content,repositoryDescriptionExp)),
		Url:         githubUrl + stringFormat(findFirstOrDefaultMatch(content,repositoryUrlExp)),
	}

	index,_ := strconv.Atoi(stringFormat(findFirstOrDefaultMatch(content,developIndexExp)))
	return DeveloperTrendDataItem{
		Index:             index,
		User:              user,
		PopularRepository: repo,
	}
}

func getGithubRequestPathAndQueryString(baseUrl string,request *http.Request) string {
	request.ParseForm()

	requestPath,_ := url.Parse(baseUrl)
	requestPath.Path += request.URL.Path
	sinceValue := request.Form.Get("since")
	spokenValue := request.Form.Get("spoken_language_code")

	params := url.Values{}
	if sinceValue != "" {
		var since = Since(strings.ToLower(sinceValue))
		if err:= since.IsValid(); err == nil {
			params.Add("since",sinceValue)
		}
	}

	if spokenValue != "" {
		var spoken = Spoken(strings.ToLower(spokenValue))
		if err:= spoken.IsValid(); err == nil {
			params.Add("spoken_language_code",spokenValue)
		}
	}

	requestPath.RawQuery = params.Encode()
	res := requestPath.String()
	fmt.Println("request path :" + res)
	return res
}

func stringFormat(content string) string {
	content = strings.Replace(content, "\n", "", -1)
	content = strings.TrimSpace(content)
	return content
}

type ApiResponse struct {
	Code 	int 			`json:"code"`
	Data	interface{}		`json:"data"`
	Date	time.Time		`json:"time"`
}

type Repository struct {
	Name        	string `json:"name"`
	Description 	string `json:"description"`
	Url         	string `json:"url"`
	Star        	string `json:"star"`
	StarToday 		string `json:"star_today"`
	Fork        	string `json:"fork"`
	Lang        	string `json:"lang"`
}


type Contribution struct {
	Level				int			`json:"level"`
	OfficialColor 		string		`json:"color"`
	Date				time.Time	`json:"time"`
	Year				int			`json:"year"`
	Month				string		`json:"month"`
	Weekday				int			`json:"weekday"`
	Total				int			`json:"total"`
}

type User struct {
	Name		string `json:"name"`
	NickName	string `json:"nick_name"`
	Avatar		string `json:"avatar"`
	Website		string `json:"website"`
}

type DeveloperTrendDataItem struct {
	Index				int		`json:"index"`
	User				User		`json:"user"`
	PopularRepository	Repository	`json:"popular_repository"`
}

type Since string
const (
	Daily 		Since = "daily"
	Weekly		Since = "weekly"
	Monthly		Since = "monthly"
)

func (lt Since) IsValid() error {
	switch lt {
	case Daily,Weekly,Monthly:
		return nil
	}
	return errors.New("Invalid since type")
}

type Spoken string
const(
	Chinese		Spoken = "zh"
	English		Spoken = "en"
)

func(sp Spoken) IsValid() error {
	switch sp {
	case Chinese,English:
		return nil
	}
	return errors.New("invalid spoken type")
}

