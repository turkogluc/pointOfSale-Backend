package entities

type Expense struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Price float64 `json:"price"`
	CreateDate int `json:"createDate"`
	UpdateDate int `json:"updateDate"`
	UserId int		`json:"userId"`
}
