package responses

type UserResponse struct {
	Count int `json:"count"`
	Items []*UserItem `json:"items"`
}

type UserItem struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	Address string `json:"address"`
	Email string `json:"email"`
	Token []byte `json:"token"`
	RegisterDate int `json:"registerDate"`
}

type UserDropdownResponse struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}