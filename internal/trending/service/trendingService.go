package service

import (
	"fmt"
	. "gtrending/internal"
	"gtrending/internal/User/model"
	. "gtrending/internal/trending/model"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)


// 获取项目排行榜
func GetTrending(requestUrl string) ([]Repository, error) {
	//no use http/2
	//http.DefaultTransport.(*http.Transport).TLSNextProto = make(map[string]func(string, *tls.Conn) http.RoundTripper)
	response,err := http.Get(requestUrl)
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
func GetDeveloperTrending(requestUrl string) ([]Developer,error) {
	res,error := http.Get(requestUrl)
	if error != nil{
		return nil,error
	}

	contentBytes, err := ioutil.ReadAll(res.Body)
	if err != nil || contentBytes == nil {
		return nil, err
	}

	return resolveDeveloperTrending(string(contentBytes)),nil
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
	name := TrimSpace(getRepositoryName(content))
	lang := TrimSpace(getRepositoryLang(content))
	desc := TrimSpace(getRepositoryDescription(content))
	star := TrimSpace(getRepositoryStar(content))
	starToday := TrimSpace(getRepositoryTodayStar(content))
	fork := TrimSpace(getRepositoryFork(content))
	url := GithubUrl + name

	if lang == "" {
		lang = "Mars"
	}

	return Repository{Name: name, Description: desc, Lang: lang, Star: star,StarToday:starToday,Fork: fork, Url: url}
}

func getRepositoryName(content string) string {
	repositoryItemNameExp := `(?<=<h1[\s\S]+<a href=")[\S]+(?=")`
	name := FindFirstOrDefaultMatchUseRegex2(content,repositoryItemNameExp)
	return name
}

func getRepositoryDescription(content string) string {
	repositoryItemDescriptionExp := `(?<=<p class="col-9 text-gray my-1 pr-4">)[\s\S]+?(?=<\/p>)`
	desc := FindFirstOrDefaultMatchUseRegex2(content,repositoryItemDescriptionExp)
	return desc
}

func getRepositoryLang(content string) string {
	repositoryItemLangExp := `(?<=<span itemprop="programmingLanguage">)[\s\S]+?(?=<\/span>)`
	return FindFirstOrDefaultMatchUseRegex2(content, repositoryItemLangExp)
}

func getRepositoryStar(content string) string {
	repositoryItemStarTagExp := `(?<=<a class="[\s]?muted-link d-inline-block mr-3"[\s\S]+stargazers">[\s\S]+g>)[\s\S]*?(?=<\/a>)`
	starValue := FindFirstOrDefaultMatchUseRegex2(content, repositoryItemStarTagExp)
	return starValue
}

func getRepositoryFork(content string) string {
	repositoryItemForkTagExp := `(?<=<a class="[\s]?muted-link d-inline-block mr-3"[\s\S]+network/[\S]*">[\s\S]+g>)[\s\S]*?(?=<\/a>)`
	forkValue := FindFirstOrDefaultMatchUseRegex2(content, repositoryItemForkTagExp)
	return forkValue
}

func getRepositoryTodayStar(content string) string {
	repositoryItemStarTagExp := `(?<=<span class="[\s]?d-inline-block float-sm-right">[\s\S]+g>)[\s\S]*?(?=stars)`
	starValue := FindFirstOrDefaultMatchUseRegex2(content, repositoryItemStarTagExp)
	return starValue
}

func resolveDeveloperTrending(content string) []Developer {
	rep := `<article class="Box-row d-flex"[\s\S]+?>[\s\S\n]+?<\/article>`
	regexp := regexp.MustCompile(rep)
	match := regexp.FindAll([]byte(content), -1)
	if match == nil {
		return nil
	}

	result := make([]Developer, 0)
	for i := 0; i < len(match); i++ {
		result = append(result, resolveDeveloperTrendDataItem(string(match[i])))
	}

	return result
}

func resolveDeveloperTrendDataItem(content string) Developer {
	developIndexExp := `(?<=<a class="text-gray f6 text-center"[\s\S]+?>)[\s\S]+?(?=<\/a>)`
	avatarExp := `(?<=<img[\s\S]+src=")[\S]+(?=")`
	userNameExp := `(?<=<h1 class="h3 lh-condensed">[\s\S]+>)[\s\S]+?(?=<\/a>[\s\S]+<\/h1>)`
	nickNameExp := `(?<=<p class="f4 text-normal mb-1">[\s\S]+>)[\s\S]+?(?=<\/a>[\s\S]+<\/p>)`

	user := model.User{
		Name:     TrimSpace(FindFirstOrDefaultMatchUseRegex2(content,userNameExp)),
		NickName: TrimSpace(FindFirstOrDefaultMatchUseRegex2(content,nickNameExp)),
		Avatar:   TrimSpace(FindFirstOrDefaultMatchUseRegex2(content,avatarExp)),
	}

	user.Website = GithubUrl + "/" + user.NickName

	repositoryNameExp := `(?<=<article>[\s\S]+h1[\s\S]+\/svg>)[\s\S]+?(?=<)`
	repositoryUrlExp := `(?<=<article>[\s\S]+href=")[\s\S]+?(?=")`
	repositoryDescriptionExp := `(?<=<div class="f6 text-gray mt-1">)[\s\S]+?(?=<)`

	repo := Repository{
		Name:        TrimSpace(FindFirstOrDefaultMatchUseRegex2(content,repositoryNameExp)),
		Description: TrimSpace(FindFirstOrDefaultMatchUseRegex2(content,repositoryDescriptionExp)),
		Url:         GithubUrl + TrimSpace(FindFirstOrDefaultMatchUseRegex2(content,repositoryUrlExp)),
	}

	index,_ := strconv.Atoi(TrimSpace(FindFirstOrDefaultMatchUseRegex2(content,developIndexExp)))
	return Developer{
		Index:             index,
		User:              user,
		PopularRepository: repo,
	}
}
