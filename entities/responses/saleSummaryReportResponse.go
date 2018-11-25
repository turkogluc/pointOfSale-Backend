package responses

import . "stock/entities"

type SaleSummaryReportResponse struct {
	Count int                  `json:"count"`
	Items []*SaleSummaryReport `json:"items"`
}
