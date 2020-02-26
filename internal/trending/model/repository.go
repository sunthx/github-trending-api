package model

type Repository struct {
	Name        		string 		`json:"name"`
	Description 		string 		`json:"description"`
	Url         		string 		`json:"url"`
	Star        		string 		`json:"star"`
	StarToday 			string 		`json:"star_today"`
	Fork        		string 		`json:"fork"`
	Lang        		string 		`json:"lang"`
	AuthorAvatar 		string		`json:"avatar"`
}
