package trending

import (
	"fmt"
	. "gtrending/internal"
	. "gtrending/internal/trending/model"
	. "gtrending/internal/trending/service"
	"net/http"
	"net/url"
	"strings"
)

func TrendRequestHandle(writer http.ResponseWriter, request *http.Request) {
	githubRequestUrl := getGithubRequestPathAndQueryString(GithubUrl,request)
	trending, _ := GetTrending(githubRequestUrl)
	OK(writer,trending)
}

func DeveloperRequestHandle(writer http.ResponseWriter,r *http.Request){
	githubRequestUrl := getGithubRequestPathAndQueryString(GithubUrl,r)
	developerTrend, _ := GetDeveloperTrending(githubRequestUrl)
	OK(writer,developerTrend)
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