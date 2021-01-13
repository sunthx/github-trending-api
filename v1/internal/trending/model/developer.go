package model

import "gtrending/internal/User/model"

type Developer struct {
	Index             int        `json:"index"`
	User              model.User `json:"user"`
	PopularRepository Repository `json:"popular_repository"`
}
