package responses

import ."stock/entities"

type ReceivingResponse struct {
	Count int					`json:"count"`
	Items []*ReceivingsItem		`json:"items"`
}
