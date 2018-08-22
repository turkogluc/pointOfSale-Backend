package responses

type PersonResponse struct {
	Count int	`json:"count"`
	Items []*PersonItem `json:"items"`
}


type PersonDropdownResponse struct {
	Id int	`json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type PersonItem struct {
	Id int				`json:"id"`
	Name string         `json:"name"`
	Phone string        `json:"phone"`
	Email string        `json:"email"`
	Address string      `json:"address"`
	Type string         `json:"type"`
	CreationDate int    `json:"creationDate"`
	UserId int			`json:"userId"`
	UserName string		`json:"userName"`
}