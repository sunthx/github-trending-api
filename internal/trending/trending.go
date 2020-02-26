package trending

import (
	"encoding/json"
	"fmt"
	. "gtrending/internal"
	. "gtrending/internal/trending/model"
	. "gtrending/internal/trending/service"
	"net/http"
	"net/url"
	"strings"
)

func TrendRequestHandle(writer http.ResponseWriter, request *http.Request) {
	githubRequestUrl := getGithubRequestPathAndQueryString(GithubUrl, request)
	trendingCacheKey := "cache:trending:repo:request:" + githubRequestUrl
	var cacheValue []Repository
	cache := GetValueFromCache(trendingCacheKey,&cacheValue)
	if cache != nil {
		fmt.Println("From Cache")
		OK(writer,cache)
		return
	}

	trending, _ := GetTrending(githubRequestUrl)

	if trending != nil && len(trending) > 0 {
		jsonResult,_ := json.Marshal(trending)
		SetValueToCache(trendingCacheKey,jsonResult)
		fmt.Println("From Request")
		OK(writer,trending)
	}

	BadRequest(writer)
}

func DeveloperRequestHandle(writer http.ResponseWriter, r *http.Request) {
	trendingCacheKey := "cache:trending:developer"
	var cacheValue []Developer
	cache := GetValueFromCache(trendingCacheKey,&cacheValue)
	if cache != nil {
		fmt.Println("From [cache:trending:developer]")
		OK(writer,cache)
		return
	}

	githubRequestUrl := getGithubRequestPathAndQueryString(GithubUrl, r)
	developerTrend, _ := GetDeveloperTrending(githubRequestUrl)

	if developerTrend != nil && len(developerTrend) > 0 {
		jsonResult,_ := json.Marshal(developerTrend)
		SetValueToCache(trendingCacheKey,jsonResult)
		fmt.Println("From Request")
		OK(writer,developerTrend)
	}

	BadRequest(writer)
}

func getGithubRequestPathAndQueryString(baseUrl string, request *http.Request) string {
	request.ParseForm()

	requestPath, _ := url.Parse(baseUrl)
	requestPath.Path += request.URL.Path
	sinceValue := request.Form.Get("since")
	spokenValue := request.Form.Get("spoken_language_code")

	params := url.Values{}
	if sinceValue != "" {
		var since = Since(strings.ToLower(sinceValue))
		if err := since.IsValid(); err == nil {
			params.Add("since", sinceValue)
		}
	}

	if spokenValue != "" {
		var spoken = Spoken(strings.ToLower(spokenValue))
		if err := spoken.IsValid(); err == nil {
			params.Add("spoken_language_code", spokenValue)
		}
	}

	requestPath.RawQuery = params.Encode()
	res := requestPath.String()
	return res
}
