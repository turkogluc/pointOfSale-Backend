package interactors

import (

	 "stock/common/projectArch/interfaces"
)

var (
	//mysql
	ProductRepo                interfaces.ProductGateway
	StockRepo                  interfaces.StockGateway
	PersonRepo                 interfaces.PersonGateway
	ReceivingRepo              interfaces.ReceivingGateway
	PaymentRepo                interfaces.PaymentGateway
	ExpenseRepo                interfaces.ExpenseGateway
	UserRepo                   interfaces.UserGateway
	SaleBasketRepo             interfaces.SaleBasketGateway
	SaleDetailRepo             interfaces.SaleDetailGateway
	SaleSummaryReportDailyRepo interfaces.SaleSummaryReportDailyGateway

)


