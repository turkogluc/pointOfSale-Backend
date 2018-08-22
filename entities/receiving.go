package entities

type Receiving struct {
	Id int					`json:"id"`
	PersonId int            `json:"personId"`
	Amount float64          `json:"amount"`
	CreationDate int        `json:"creationDate"`
	UpdateDate int          `json:"updateDate"`
	ExpectedDate int        `json:"expectedDate"`
	ProductIds string		`json:"productIds"`
	Status string           `json:"status"`
	UserId int				`json:"userId"`
}

type ReceivingsItem struct {
	Id int                   `json:"id"`
	PersonId int             `json:"personId"`
	Amount float64           `json:"amount"`
	CreationDate int         `json:"creationDate"`
	UpdateDate int           `json:"updateDate"`
	ExpectedDate int         `json:"expectedDate"`
	Status string            `json:"status"`
	PersonName string        `json:"personName"`
	PersonPhone string       `json:"personPhone"`
	ProductIds string		 `json:"productIds"`
	ProductList []*Product	 `json:"productList"`
	UserId int				`json:"userId"`
	UserName string			`json:"userName"`
}
