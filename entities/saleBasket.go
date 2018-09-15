package entities

type SaleBasketItem struct {
	Id int 	`json:"id"`
	Qty int			`json:"qty"`
	Discount float64	`json:"discount"`
}

type SaleBasket struct {
	Id int `json:"id"`
	CreationDate int `json:"creationDate"`
	Items []*SaleBasketItem `json:"items"`
	ItemsStr string   `json:"itemStr"`
	UserId int 		`json:"userId"`
	UserName string  `json:"userName"`
	TotalDiscount float64 `json:"totalDiscount"`
	TotalPrice	float64		`json:"totalPrice"`
}

type SaleBasketItems struct {
	Items	[]*SaleBasketItem	`json:"items"`
}