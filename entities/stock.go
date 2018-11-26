package entities

type Stock struct {
	Id int				`json:"id"`
	ProductId int       `json:"productId"`
	Qty int             `json:"qty"`
	DealerId int        `json:"dealerId"`
	CreationDate int    `json:"creationDate"`
	UpdateDate int      `json:"updateDate"`
	UserId int			`json:"userId"`
	IsFavorite bool		`json:"isFavorite"`
}
