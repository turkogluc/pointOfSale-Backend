package responses

import . "stock/entities"

type ExpenseResponse struct{
	Count int 				`json:"count"`
	Items []*Expense		`json:"items"`
}

//type ExpenseDropdownResponse struct {
//	Id int `json:"id"`
//	Name string `json:"name"`
//	Price float64 `json:"price"`
//}
