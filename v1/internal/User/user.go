package User

import (
	"encoding/json"
	"gtrending/internal"
	"gtrending/internal/User/model"
	"gtrending/internal/User/service"
	"net/http"
)

func DetailRequestHandle(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	userName := request.Form.Get("name")

	if userName == "" {
		internal.BadRequest(writer)
		return
	}

	userCacheKey := "cache:user:" + userName
	var cacheValue model.User
	cache := internal.GetValueFromCache(userCacheKey,&cacheValue)
	if cache != nil {
		internal.OK(writer,cache)
		return
	}

	user,err := service.GetUser(userName)
	if err != nil {
		internal.BadRequest(writer)
		return
	}

	jsonResult,_ := json.Marshal(user)
	internal.SetValueToCache(userCacheKey,jsonResult)
	internal.OK(writer,user)
}