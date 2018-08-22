package entities

type Person struct {
	Id int				`json:"id"`
	Name string         `json:"name"`
	Phone string        `json:"phone"`
	Email string        `json:"email"`
	Address string      `json:"address"`
	Type string         `json:"type"`
	CreationDate int    `json:"creationDate"`
	UserId int			`json:"userId"`
}