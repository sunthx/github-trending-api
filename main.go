package main

import (
	"fmt"
	"gtrending/internal/contribution"
	"gtrending/internal/trending"
	"net/http"
)


func main() {
	listenPort := ":8080"
	contributionRouteString := "/contributions"
	trendingRouteString := "/trending/"
	developerTrendingRouteString := "/trending/developers/"

	fmt.Println("github_trending_api_server running...\n" +
		listenPort + "\n" +
		contributionRouteString + "?user=[username]\n" +
		trendingRouteString + "\n" +
		developerTrendingRouteString)

	http.HandleFunc(contributionRouteString, contribution.ContributionRequestHandle)
	http.HandleFunc(trendingRouteString, trending.TrendRequestHandle)
	http.HandleFunc(developerTrendingRouteString,trending.DeveloperRequestHandle)
	http.ListenAndServe(listenPort, nil)
}
















