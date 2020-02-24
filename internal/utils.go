package internal

import (
	"encoding/json"
	"fmt"
	"github.com/dlclark/regexp2"
	"net/http"
	"time"
)

func FindFirstOrDefaultMatchUseRegex2(content string, exp string) string {
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

func OK(response http.ResponseWriter,data interface{}){
	result := ApiResponse{
		Code: 200,
		Data: data,
		Date: time.Now(),
	}

	jsonValue, _ := json.Marshal(result)
	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(response, string(jsonValue))
}