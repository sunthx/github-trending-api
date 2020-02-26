package User

import (
	"fmt"
	"gtrending/internal"
	"gtrending/internal/User/service"
	"net/http"
)

func DetailRequestHandle(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	userName := request.Form.Get("name")

	fmt.Println("userName" + userName)
	if userName == "" {
		internal.BadRequest(writer)
		return
	}

	user,err := service.GetUser(userName)
	if err != nil {
		internal.BadRequest(writer)
		return
	}

	internal.OK(writer,user)
}