package entities

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	Address string `json:"address"`
	Email string `json:"email"`
	Password string `json:"password"`
	Token []byte `json:"token"`
	RegisterDate int `json:"registerDate"`
}
