package entities

type SaleSummaryReport struct {
	Id 				int 		`json:"id"`
	GrossProfits	[]float64		`json:"grossProfits"`
	GrossProfit		float64		`json:"grossProfit"`
	NetProfits		[]float64		`json:"netProfits"`
	NetProfit		float64		`json:"netProfit"`
	SaleCounts		[]int			`json:"saleCounts"`
	SaleCount		int			`json:"saleCount"`
	ItemCounts		[]int			`json:"itemCounts"`
	ItemCount		int			`json:"itemCount"`
	CustomerCounts	[]int			`json:"customerCounts"`
	CustomerCount	int			`json:"customerCount"`
	Discounts		[]float64		`json:"discounts"`
	Discount		float64		`json:"discount"`
	BasketValues		[]float64		`json:"basketValues"`
	BasketValue		float64		`json:"basketValue"`
	BasketSizes		[]float64		`json:"basketSizes"`
	BasketSize		float64		`json:"basketSize"`
	Timestamp		int			`json:"timestamp"`
	Timestamps		[]int			`json:"timestamps"`

	AsObject		*SaleSummaryObject `json:"asObject"`

}

type SaleSummaryObject struct{
	Count int 						`json:"count"`
	Items []*SaleSummaryObjectItem	`json:"items"`
}

type SaleSummaryObjectItem struct {
	Id				int			`json:"id"`
	GrossProfit		float64		`json:"grossProfit"`
	NetProfit		float64		`json:"netProfit"`
	SaleCount		int			`json:"saleCount"`
	ItemCount		int			`json:"itemCount"`
	CustomerCount	int			`json:"customerCount"`
	Discount		float64		`json:"discount"`
	BasketValue		float64		`json:"basketValue"`
	BasketSize		float64		`json:"basketSize"`
	Timestamp		int			`json:"timestamp"`
}