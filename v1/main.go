package main

import (
	"fmt"
	"gtrending/internal/User"
	"gtrending/internal/contribution"
	"gtrending/internal/trending"
	"net/http"
)


func main() {
	listenPort := ":8080"
	contributionRouteString := "/contributions"
	trendingRouteString := "/trending/"
	developerTrendingRouteString := "/trending/developers/"
	userRoutString := "/user"

	fmt.Println("github_trending_api_server running...\n" +
		listenPort + "\n" +
		contributionRouteString + "?user=[username]\n" +
		trendingRouteString + "\n" +
		developerTrendingRouteString + "\n" +
		userRoutString + "?name=[username]")

	http.HandleFunc(contributionRouteString, contribution.UserContributionRequestHandle)
	http.HandleFunc(trendingRouteString, trending.TrendRequestHandle)
	http.HandleFunc(developerTrendingRouteString,trending.DeveloperRequestHandle)
	http.HandleFunc(userRoutString,User.DetailRequestHandle)

	http.ListenAndServe(listenPort, nil)
}
















