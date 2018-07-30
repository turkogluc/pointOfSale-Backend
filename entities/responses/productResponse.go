package responses

import . "stock/entities"

type ProductResponse struct{
	Count int 				`json:"count"`
	Items []*Product		`json:"items"`
}

type ProductDropdownResponse struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Price float64 `json:"price"`
}
