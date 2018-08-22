package entities

type Payment struct {
	Id int					`json:"id"`
	PersonId int            `json:"personId"`
	Amount float64          `json:"amount"`
	CreationDate int        `json:"creationDate"`
	UpdateDate int          `json:"updateDate"`
	ExpectedDate int        `json:"expectedDate"`
	Status string           `json:"status"`
	Summary string 			`json:"summary"`
	UserId int			`json:"userId"`
}

type PaymentsItem struct {
	Id int                   `json:"id"`
	PersonId int             `json:"personId"`
	Amount float64           `json:"amount"`
	CreationDate int         `json:"creationDate"`
	UpdateDate int           `json:"updateDate"`
	ExpectedDate int         `json:"expectedDate"`
	Status string            `json:"status"`
	PersonName string        `json:"personName"`
	PersonPhone string       `json:"personPhone"`
	Summary string 			 `json:"summary"`
	UserId int			`json:"userId"`
	UserName string			`json:"userName"`
}
