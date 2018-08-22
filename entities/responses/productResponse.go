package responses


type ProductResponse struct{
	Count int 				`json:"count"`
	Items []*ProductItem		`json:"items"`
}

type ProductDropdownResponse struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Price float64 `json:"price"`
}

type ProductItem struct {
	Id            int     `json:"id"`
	Barcode       string  `json:"barcode"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Category      string  `json:"category"`
	PurchasePrice float64 `json:"purchasePrice"`
	SalePrice     float64 `json:"salePrice"`
	RegisterDate  int  `json:"registerDate"`
	UserId 		  int 	`json:"userId"`
	UserName 	string 	`json:"userName"`
}