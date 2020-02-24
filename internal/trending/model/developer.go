package model

type Developer struct {
	Index             int        `json:"index"`
	User              User       `json:"user"`
	PopularRepository Repository `json:"popular_repository"`
}
