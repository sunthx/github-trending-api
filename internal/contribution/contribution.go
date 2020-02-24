package contribution

import (
	"gtrending/internal"
	"gtrending/internal/contribution/service"
	"net/http"
)

func ContributionRequestHandle(writer http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		return
	}

	user := request.Form.Get("user")
	contributions, _ := service.GetContributions(user)
	internal.OK(writer,contributions)
}
