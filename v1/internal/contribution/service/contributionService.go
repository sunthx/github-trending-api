package service

import (
	"errors"
	"github.com/dlclark/regexp2"
	"gtrending/internal"
	. "gtrending/internal/contribution/model"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

// 获取指定用户活跃榜
func GetContributions(userName string) ([]Contribution, error) {
	if userName == "" {
		return nil, errors.New("required:username")
	}

	requestUrl := internal.GithubUrl +"/users/" + userName + "/contributions"
	res, err := http.Get(requestUrl)
	if res.StatusCode != http.StatusOK || err != nil {
		return nil, err
	}

	contentBytes, err := ioutil.ReadAll(res.Body)
	if err != nil || contentBytes == nil {
		return nil, err
	}

	return resolveContributions(string(contentBytes)), nil
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
	exp := `<rect.*?\>`
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

