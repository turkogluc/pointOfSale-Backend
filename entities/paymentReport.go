package entities

type PaymentReport struct {
	TotalExpenses		float64 	`json:"totalExpenses"`
	Expenses			[]float64 	`json:"expenses"`
	TotalPayments		float64		`json:"totalPayments"`
	Payments			[]float64 	`json:"payments"`
	OverduePayments		int		`json:"overduePayments"`
	TotalReceivings		float64		`json:"totalReceivings"`
	Receivings			[]float64		`json:"receivings"`
	OverdueReceivings	int		`json:"overdueReceivings"`
	Timestamps []string  `json:"timestamps"`
	ItemsAsObject		[]*PaymentList	`json:"itemsAsObject"`

}

type PaymentList struct {
	Person		string 	`json:"person"`
	Amount  	float64	`json:"amount"`
	Timestamp	int		`json:"timestamp"`
	Status		string `json:"status"`
	Detail		string	`json:"detail"`
	Type 		string 	`json:"type"`
}