package responses

import ."stock/entities"

type PaymentResponse struct {
	Count int					`json:"count"`
	Items []*PaymentsItem		`json:"items"`
}
