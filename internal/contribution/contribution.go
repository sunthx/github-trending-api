package contribution

import (
	"encoding/json"
	"gtrending/internal"
	. "gtrending/internal/contribution/model"
	"gtrending/internal/contribution/service"
	"net/http"
)

func UserContributionRequestHandle(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		internal.BadRequest(writer)
		return
	}

	user := request.Form.Get("user")
	if user == "" {
		internal.BadRequest(writer)
		return
	}

	userCacheKey := "cache:contribution:" + user
	var cacheValue []Contribution
	cache := internal.GetValueFromCache(userCacheKey,&cacheValue)
	if cache != nil {
		internal.OK(writer,cache)
		return
	}

	contributions, _ := service.GetContributions(user)
	if contributions != nil && len(contributions) > 0{
		jsonResult,_ := json.Marshal(user)
		internal.SetValueToCache(userCacheKey,jsonResult)
		internal.OK(writer,contributions)
		return
	}

	internal.BadRequest(writer)
}
