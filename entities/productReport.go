package entities

// product | numOfSales | grossProfit | netProfit | totalDiscount | numOfReturn |
// markup-kar marjı (satış fiyatı - alış fiyatı) / alış
//

type ProductReportItem struct {
	ProductId	int		`json:"productId"`
	ProductName	string	`json:"productName"`
	Qty		int			`json:"qty"`
	GrossProfit	float64	`json:"grossProfit"`
	NetProfit	float64	`json:"netProfit"`
	Discount	float64	`json:"discount"`
	Markup		float64	`json:"markup"`
	ProfitPercentage float64	`json:"profitPercentage"`
	NumberOfReturn	int	 `json:"numberOfReturn"`
}

type ProductReport struct {
	Count int	`json:"count"`
	Items	[]*ProductReportItem	`json:"items"`
	TotalQty	int		`json:"totalQty"`
	TotalGrossProfit	float64	`json:"totalGrossProfit"`
	TotalNetProfit	float64	`json:"totalNetProfit"`
	TotalDiscount	float64	`json:"totalDiscount"`
	TotalNumberOfReturn	int	 `json:"totalNumberOfReturn"`
}