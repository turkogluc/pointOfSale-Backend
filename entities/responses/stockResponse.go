package responses

import . "stock/entities"

type StockItem struct {
	Id int				`json:"id"`
	Product *Product    `json:"product"`
	Qty int             `json:"qty"`
	DealerId int        `json:"dealerId"`
	DealerName string	`json:"dealerName"`
	CreationDate int    `json:"creationDate"`
	UpdateDate int      `json:"updateDate"`
	UserId int			`json:"userId"`
	UserName string 	`json:"userName"`
}

type StockResponse struct {
	Count int	`json:"count"`
	Items []*StockItem	`json:"items"`
}

type CurrentStockReportItem struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Category string `json:"category"`
	Qty int	`json:"qty"`
	PurchasePrice float64 `json:"purchasePrice"`
	SalePrice float64 `json:"salePrice"`
	GrossValue float64 `json:"grossValue"`
	NetValue float64  `json:"netValue"`
	TotalProfit	float64 `json:"totalProfit"`
}

type CurrentStockReportResponse struct {
	Count int							`json:"count"`
	Items []*CurrentStockReportItem		`json:"items"`
	Total CurrentStockReportItem		`json:"total"`
}
