package entities


type SaleDetail struct {
	Id int `json:"id"`
	CreationDate int `json:"creationDate"`
	BasketId	int		`json:"basketId"`
	ProductId	int		`json:"productId"`
	Qty			int		`json:"qty"`
	Discount	float64		`json:"discount"`
	UserId int 		`json:"userId"`
}