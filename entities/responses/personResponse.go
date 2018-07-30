package responses

import . "stock/entities"

type PersonResponse struct {
	Count int	`json:"count"`
	Items []*Person `json:"items"`
}


type PersonDropdownResponse struct {
	Id int	`json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}