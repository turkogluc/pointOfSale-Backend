package entities

type LoginParams struct {
	Email 		string		`json:"email"`
	Password 	string		`json:"password"`
}

type LoginResponse struct{
	Token 	string   		`json:"token"`
	User 	*User        	`json:"user"`
}