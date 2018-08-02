package interfaces

import (
	. "stock/entities"
	"stock/entities/responses"
)

// mysql

type ProductGateway interface {
	SelectProductById(id int)(*Product,error)
	SelectProducts(barcode,name,description,category,orderBy,orderAs string,pageNumber, pageSize int) (*responses.ProductResponse,  error)
	InsertProduct(p *Product)(error)
	UpdateProductById(p *Product, IdToUpdate int)(error)
	DeleteProductById(Id int)(error)
	DeleteProducts(ids []int)(error)
	Close()
}

type StockGateway interface {
	SelectStockById(id int)(*Stock,error)
	SelectStocks(barcode,name,description,category,orderBy,orderAs string,pageNumber, pageSize,dealerId int) (*responses.StockResponse,  error)
	InsertStock(p *Stock)(error)
	UpdateStockById(p *Stock, IdToUpdate int)(error)
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
	SelectReceivings(person,status,orderBy,orderAs string,pageNumber, pageSize int) (*responses.ReceivingResponse,  error)
	InsertReceiving(p *Receiving)(error)
	UpdateReceivingById(p *Receiving, IdToUpdate int)(error)
	DeleteReceivingById(Id int)(error)
	DeleteReceivings(ids []int)(error)
	Close()
}

type PaymentGateway interface {
	SelectPaymentById(id int)(*Payment,error)
	SelectPayments(person,status,orderBy,orderAs string,pageNumber, pageSize int) (*responses.PaymentResponse,  error)
	InsertPayment(p *Payment)(error)
	UpdatePaymentById(p *Payment, IdToUpdate int)(error)
	DeletePaymentById(Id int)(error)
	DeletePayments(ids []int)(error)
	Close()
}

type ExpenseGateway interface {
	SelectExpenseById(id int)(*Expense,error)
	SelectExpenses(name,description,orderBy,orderAs string,pageNumber, pageSize int) (*responses.ExpenseResponse,  error)
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

type SaleGateway interface {
	SelectSaleById(id int)(*Sale,error)
	SelectSales(timeInterval []int,orderBy,orderAs string,pageNumber, pageSize int) (*responses.SaleResponse,  error)
	InsertSale(p *Sale)(error)
	UpdateSaleById(p *Sale, IdToUpdate int)(error)
	DeleteSaleById(Id int)(error)
	DeleteSales(ids []int)(error)
	Close()
}