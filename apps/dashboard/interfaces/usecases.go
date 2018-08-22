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

	CreateStock(p *Stock) *ErrorType
	UpdateStock(p *Stock) *ErrorType
	GetStockById(id int) (*Stock,*ErrorType)
	GetStocks(barcode,name,description,category,orderBy,orderAs string,pageNumber, pageSize,dealerId int) (*responses.StockResponse,  *ErrorType)
	DeleteStocks(ids []int) *ErrorType

	CreatePerson(p *Person) *ErrorType
	UpdatePerson(p *Person) *ErrorType
	GetPersonById(id int) (*Person,*ErrorType)
	GetPeople(name,pType,orderBy,orderAs string,pageNumber, pageSize int) (*responses.PersonResponse,  *ErrorType)
	DeletePersons(ids []int) *ErrorType

	CreateReceiving(p *Receiving) *ErrorType
	UpdateReceiving(p *Receiving) *ErrorType
	GetReceivingById(id int) (*Receiving,*ErrorType)
	GetReceivings(person,status,orderBy,orderAs string,pageNumber, pageSize int) (*responses.ReceivingResponse,  *ErrorType)
	DeleteReceivings(ids []int) *ErrorType
	SetReceivingStatus(status string,id int) *ErrorType

	CreatePayment(p *Payment) *ErrorType
	UpdatePayment(p *Payment) *ErrorType
	GetPaymentById(id int) (*Payment,*ErrorType)
	GetPayments(person,status,orderBy,orderAs string,pageNumber, pageSize int) (*responses.PaymentResponse,  *ErrorType)
	DeletePayments(ids []int) *ErrorType
	SetPaymentStatus(status string,id int) *ErrorType

	CreateExpense(p *Expense) *ErrorType
	UpdateExpense(p *Expense) *ErrorType
	GetExpenseById(id int) (*Expense,*ErrorType)
	GetExpenses(name,description,orderBy,orderAs string,pageNumber, pageSize int) (*responses.ExpenseResponse,  *ErrorType)
	DeleteExpenses(ids []int) *ErrorType

	CreateUser(p *User) *ErrorType
	UpdateUser(p *User) *ErrorType
	GetUserById(id int) (*User,*ErrorType)
	GetUsers(name,email,orderBy,orderAs string,pageNumber, pageSize int) (*responses.UserResponse,  *ErrorType)
	DeleteUsers(ids []int) *ErrorType

	CreateSale(p *Sale) *ErrorType
	UpdateSale(p *Sale) *ErrorType
	GetSaleById(id int) (*Sale,*ErrorType)
	GetSales(timeInterval,orderBy,orderAs string,pageNumber, pageSize int) (*responses.SaleResponse,  *ErrorType)
	DeleteSales(ids []int) *ErrorType

	Login(email, password string, secret string) (*User, string, *ErrorType)
	FillProductTable() *ErrorType
}
