package model

import "time"

type Contribution struct {
	Level				int			`json:"level"`
	OfficialColor 		string		`json:"color"`
	Date				time.Time	`json:"time"`
	Year				int			`json:"year"`
	Month				string		`json:"month"`
	Weekday				int			`json:"weekday"`
	Total				int			`json:"total"`
}
