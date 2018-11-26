package interfaces

import (
	. "stock/entities"
	"stock/entities/responses"
)

// mysql

type ProductGateway interface {
	SelectProductById(id int)(*Product,error)
	SelectProductCategories()([]string,error)
	SelectProducts(barcode,name,description,category,orderBy,orderAs string,pageNumber, pageSize int) (*responses.ProductResponse,  error)
	InsertProduct(p *Product)(error)
	UpdateProductById(p *Product, IdToUpdate int)(error)
	DeleteProductById(Id int)(error)
	DeleteProducts(ids []int)(error)
	Close()
}

type StockGateway interface {
	SelectStockById(id int)(*Stock,error)
	SelectStockByProductId(productId int)(*Stock,error)
	SelectStocks(timeInterval []int,barcode,name,description,category,orderBy,orderAs string,pageNumber, pageSize,dealerId int,creatorId int,isFavorite bool) (*responses.StockResponse,  error)
	SelectCurrentStockReport(name,category,orderBy,orderAs string,pageNumber, pageSize int) (*responses.CurrentStockReportResponse,  error)
	InsertStock(p *Stock)(error)
	UpdateStockById(p *Stock, IdToUpdate int)(error)
	DecrementProductFromStock(productId,count int)(error)
	SetFavoriteByProductId(productId int,fav bool)(error)
	DeleteStockById(Id int)(error)
	DeleteStocks(ids []int)(error)
	Close()
}

type PersonGateway interface {
	SelectPersonById(id int)(*Person,error)
	SelectPeople(name,pType,orderBy,orderAs string,pageNumber, pageSize int) (*responses.PersonResponse,  error)
	InsertPerson(p *Person)(error)
	UpdatePersonById(p *Person, IdToUpdate int)(error)
	DeletePersonById(Id int)(error)
	DeletePersons(ids []int)(error)
	Close()
}

type ReceivingGateway interface {
	SelectReceivingById(id int)(*Receiving,error)
	SelectReceivings(timeInterval []int,person,status,orderBy,orderAs string,pageNumber, pageSize,creator int) (*responses.ReceivingResponse,  error)
	InsertReceiving(p *Receiving)(error)
	UpdateReceivingById(p *Receiving, IdToUpdate int)(error)
	DeleteReceivingById(Id int)(error)
	DeleteReceivings(ids []int)(error)
	SetStatus(status string, IdToUpdate int)(error)
	Close()
}

type PaymentGateway interface {
	SelectPaymentById(id int)(*Payment,error)
	SelectPayments(timeInterval []int,person,status,orderBy,orderAs string,pageNumber, pageSize,creator int) (*responses.PaymentResponse,  error)
	InsertPayment(p *Payment)(error)
	UpdatePaymentById(p *Payment, IdToUpdate int)(error)
	SetPaymentStatus(status string, IdToUpdate int)(error)
	DeletePaymentById(Id int)(error)
	DeletePayments(ids []int)(error)
	Close()
}

type ExpenseGateway interface {
	SelectExpenseById(id int)(*Expense,error)
	SelectExpenses(timeInterval []int,name,description,orderBy,orderAs string,pageNumber, pageSize,creator int) (*responses.ExpenseResponse,  error)
	InsertExpense(p *Expense)(error)
	UpdateExpenseById(p *Expense, IdToUpdate int)(error)
	DeleteExpenseById(Id int)(error)
	DeleteExpenses(ids []int)(error)
	Close()
}

type UserGateway interface {
	SelectUserById(id int)(*User,error)
	SelectUserByEmail(email string)(*User,error)
	SelectUsers(name,email,orderBy,orderAs string,pageNumber, pageSize int) (*responses.UserResponse,  error)
	InsertUser(p *User)(error)
	UpdateUserById(p *User, IdToUpdate int)(error)
	DeleteUserById(Id int)(error)
	DeleteUsers(ids []int)(error)
	Close()
}

type SaleBasketGateway interface {
	SelectSaleBasketById(id int)(*SaleBasket,error)
	SelectSaleBaskets(timeInterval []int,userId int,orderBy,orderAs string,pageNumber, pageSize int) (*responses.SaleBasketResponse,  error)
	RetrieveNotProcessedRecords()(*SaleSummaryObject,error)
	InsertSaleBasket(p *SaleBasket)(error)
	UpdateSaleBasketById(p *SaleBasket, IdToUpdate int)(error)
	SetSaleBasketIsProcessedStatus(IdToUpdate int,status bool)(error)
	DeleteSaleBasketById(Id int)(error)
	DeleteSaleBaskets(ids []int)(error)
	Close()
}

type SaleDetailGateway interface {
	SelectSaleDetailById(id int)(*SaleDetail,error)
	InsertSaleDetail(p *SaleDetail)(error)
	UpdateSaleDetailById(p *SaleDetail, IdToUpdate int)(error)
	DeleteSaleDetailById(Id int)(error)
	DeleteSaleDetails(ids []int)(error)
	SelectSaleDetails(timeInterval []int,productName string,category string, userId int)(*ProductReport,error)
	//RetrieveNotProcessedRecords()(*SaleSummaryObject,error)
	//ChangeIsProcessedStatus(status bool, basketIdToUpdate int)(error)
}

type SaleSummaryReportGateway interface {
	SelectSaleSummaryReportById(id int)(*SaleSummaryReport,error)
	InsertSaleSummaryReport(p *SaleSummaryObjectItem)(error)
	UpdateSaleSummaryReportById(p *SaleSummaryObjectItem, IdToUpdate int)(error)
	DeleteSaleSummaryReportById(Id int)(error)
	DeleteSaleSummaryReport(ids []int)(error)
	SelectSaleSummaryReportItems(timeInterval []int) (*SaleSummaryReport,  error)
	SelectSaleSummaryReportByDate(id int)(*SaleSummaryObjectItem,error)
	Close()
}