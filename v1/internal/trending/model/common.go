package model

import "errors"

type Since string
const (
	Daily 		Since = "daily"
	Weekly		Since = "weekly"
	Monthly		Since = "monthly"
)

func (lt Since) IsValid() error {
	switch lt {
	case Daily,Weekly,Monthly:
		return nil
	}
	return errors.New("Invalid since type")
}

type Spoken string
const(
	Chinese		Spoken = "zh"
	English		Spoken = "en"
)

func(sp Spoken) IsValid() error {
	switch sp {
	case Chinese,English:
		return nil
	}
	return errors.New("invalid spoken type")
}
