package internal

import "time"

type ApiResponse struct {
	Code 	int 			`json:"code"`
	Data	interface{}		`json:"data"`
	Date	time.Time		`json:"time"`
}
