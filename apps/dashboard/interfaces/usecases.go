package interfaces

import (
	. "stock/entities"
	"stock/entities/responses"
)

type DashboardUseCases interface {
	CreateProduct(p *Product) *ErrorType
	UpdateProduct(p *Product) *ErrorType
	GetProductById(id int) (*Product,*ErrorType)
	GetProducts(barcode,name,description,category,orderBy,orderAs string,pageNumber, pageSize int) (*responses.ProductResponse,  *ErrorType)
	DeleteProducts(ids []int) *ErrorType
	RetrieveCategories()([]string,*ErrorType)

	CreateStock(p *Stock) *ErrorType
	UpdateStock(p *Stock) *ErrorType
	GetStockById(id int) (*Stock,*ErrorType)
	GetStocks(timeInterval string,barcode,name,description,category,orderBy,orderAs string,pageNumber, pageSize,dealerId,userId int,isFavorite bool) (*responses.StockResponse,  *ErrorType)
	DeleteStocks(ids []int) *ErrorType
	SetFavoriteProduct(productId int, isFavorite bool) *ErrorType

	CreatePerson(p *Person) *ErrorType
	UpdatePerson(p *Person) *ErrorType
	GetPersonById(id int) (*Person,*ErrorType)
	GetPeople(name,pType,orderBy,orderAs string,pageNumber, pageSize int) (*responses.PersonResponse,  *ErrorType)
	DeletePersons(ids []int) *ErrorType

	CreateReceiving(p *Receiving) *ErrorType
	UpdateReceiving(p *Receiving) *ErrorType
	GetReceivingById(id int) (*Receiving,*ErrorType)
	GetReceivings(timeInterval string,person,status,orderBy,orderAs string,pageNumber, pageSize,creator int) (*responses.ReceivingResponse,  *ErrorType)
	DeleteReceivings(ids []int) *ErrorType
	SetReceivingStatus(status string,id int) *ErrorType

	CreatePayment(p *Payment) *ErrorType
	UpdatePayment(p *Payment) *ErrorType
	GetPaymentById(id int) (*Payment,*ErrorType)
	GetPayments(timeInterval string,person,status,orderBy,orderAs string,pageNumber, pageSize,creator int) (*responses.PaymentResponse,  *ErrorType)
	DeletePayments(ids []int) *ErrorType
	SetPaymentStatus(status string,id int) *ErrorType

	CreateExpense(p *Expense) *ErrorType
	UpdateExpense(p *Expense) *ErrorType
	GetExpenseById(id int) (*Expense,*ErrorType)
	GetExpenses(timeInterval string,name,description,orderBy,orderAs string,pageNumber, pageSize,creator int) (*responses.ExpenseResponse,  *ErrorType)
	DeleteExpenses(ids []int) *ErrorType

	CreateUser(p *User) *ErrorType
	UpdateUser(p *User) *ErrorType
	GetUserById(id int) (*User,*ErrorType)
	GetUsers(name,email,orderBy,orderAs string,pageNumber, pageSize int) (*responses.UserResponse,  *ErrorType)
	DeleteUsers(ids []int) *ErrorType

	CreateSaleBasket(p *SaleBasket) *ErrorType
	UpdateSaleBasket(p *SaleBasket) *ErrorType
	GetSaleBasketById(id int) (*SaleBasket,*ErrorType)
	GetSaleBaskets(timeInterval string,userId int,orderBy,orderAs string,pageNumber, pageSize int) (*responses.SaleBasketResponse,  *ErrorType)
	DeleteSaleBaskets(ids []int) *ErrorType

	// # Reports #

	GetSaleSummaryReport(tInterval string) (*SaleSummaryReport,  *ErrorType)
	GetCurrentStockReport(name,category,orderBy,orderAs string,pageNumber, pageSize int) (*responses.CurrentStockReportResponse,  *ErrorType)
	GetActivityLog(tInterval string,userId int)(*ActivityLogs,*ErrorType)
	GetPaymentReport(tInterval string) (*PaymentReport,  *ErrorType)
	GetProductReport(tInterval string,productName string,category string,userId int) (*ProductReport,  *ErrorType)


	// # Reports As Excel#

	GetSaleSummaryReportAsExcel(tInterval string) (string,  *ErrorType)
	GetCurrentStockReportAsExcel(name,category,orderBy,orderAs string,pageNumber, pageSize int) (string,  *ErrorType)
	GetPaymentReportAsExcel(tInterval string) (string,  *ErrorType)
	GetProductReportAsExcel(tInterval string,productName string,category string,userId int)(string,*ErrorType)

	// # Util #

	Login(email, password string, secret string) (*User, string, *ErrorType)
	FillProductTable(userId int) *ErrorType
}
