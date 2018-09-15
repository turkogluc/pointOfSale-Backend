package entities

type SaleItem struct {
	barcode string 	`json:"barcode"`
	qty int			`json:"qty"`
}

type Sale struct {
	Id int `json:"id"`
	CreationDate int `json:"creationDate"`
	Items []*SaleItem `json:"items"`
	ItemsStr string   `json:"itemStr"`
	UserId int 		`json:"userId"`
	UserName string  `json:"userName"`
}