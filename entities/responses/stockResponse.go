package responses

import . "stock/entities"

type StockItem struct {
	Id int				`json:"id"`
	Product *Product       `json:"productId"`
	Qty int             `json:"qty"`
	DealerId int        `json:"dealerId"`
	DealerName string	`json:"dealerName"`
	CreationDate int    `json:"creationDate"`
	UpdateDate int      `json:"updateDate"`
}

type StockResponse struct {
	Count int	`json:"count"`
	Items []*StockItem	`json:"items"`
}

