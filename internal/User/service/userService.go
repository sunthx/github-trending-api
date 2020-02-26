package service

import (
	"gtrending/internal"
	"gtrending/internal/User/model"
	"io/ioutil"
	"net/http"
)

func GetUser(userName string) (model.User,error) {
	requestUrl := internal.GithubUrl + "/" + userName
	response,err := http.Get(requestUrl)
	if err != nil {
		return model.User{},err
	}

	contentBytes, err := ioutil.ReadAll(response.Body)
	if err != nil || contentBytes == nil {
		return model.User{}, err
	}

	user := resolveUser(string(contentBytes))
	user.Website = requestUrl
	return user,err
}

func resolveUser(content string) model.User {
	avatarExp := `(?<=<img[\s\S]+?avatar-before-user-status" src=")\S+?(?=")`
	nameExp := `(?<=<span[\s\S]+?itemprop="name">)\S+?(?=<)`
	nickNameExp :=`(?<=<span[\s\S]+?itemprop="additionalName">)\S+?(?=<)`

	res := model.User{}
	res.Avatar = internal.TrimSpace(internal.FindFirstOrDefaultMatchUseRegex2(content,avatarExp))
	res.Name= internal.TrimSpace(internal.FindFirstOrDefaultMatchUseRegex2(content,nameExp))
	res.NickName= internal.TrimSpace(internal.FindFirstOrDefaultMatchUseRegex2(content,nickNameExp))

	return res
}