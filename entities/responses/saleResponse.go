package responses

import . "stock/entities"

type SaleResponse struct {
	Count int `json:"count"`
	Items []*Sale `json:"items"`
}
