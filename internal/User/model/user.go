package model

type User struct {
	Name		string `json:"name"`
	NickName	string `json:"nick_name"`
	Avatar		string `json:"avatar"`
	Website		string `json:"website"`
}