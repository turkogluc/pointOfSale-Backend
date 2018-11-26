package entities

type Product struct {
	Id            int     `json:"id"`
	Barcode       string  `json:"barcode"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Category      string  `json:"category"`
	PurchasePrice float64 `json:"purchasePrice"`
	SalePrice     float64 `json:"salePrice"`
	RegisterDate  int  `json:"registerDate"`
	UserId 		  int 	`json:"userId"`
	ImagePath	string	`json:"imagePath"`
}