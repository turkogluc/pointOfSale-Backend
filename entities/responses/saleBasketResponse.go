package responses

import . "stock/entities"

type SaleBasketResponse struct {
	Count int `json:"count"`
	Items []*SaleBasket `json:"items"`
}
