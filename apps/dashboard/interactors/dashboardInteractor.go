package interactors

import(
	jwt_lib "stock/lib/jwt-go"
	. "stock/entities"
	"stock/common/projectArch/interactors"
	. "stock/common/logger"
	"stock/entities/responses"
	"time"
	"strings"
	"strconv"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
)

var errorMap map[string]map[int]string
type DashboardInteractor struct{}


func init(){
	errorMap = make(map[string]map[int]string)

	errorMap["tr"] = make(map[int]string)
	errorMap["en"] = make(map[int]string)


	errorMap["tr"][10001] = "Kullanici Email adresi bulunamadi"
	errorMap["tr"][10002] = "Kullanici bulunamadi"
	errorMap["tr"][10003] = "Yanlis sifre girildi"
	errorMap["tr"][10004] = "Token olusturulamadi"
	errorMap["tr"][10005] = "Kullanici adi bulunamadi"


	errorMap["en"][10100] = "Fuel entry couldn't be found"
	errorMap["en"][10101] = "An error occured while creating the vehicle"
	errorMap["en"][10102] = "The vehicle to be updated couldn't be found"
	errorMap["en"][10103] = "The vehicle to be deleted couldn't be found"

}

func (DashboardInteractor) Login(email, password string, secret string) (*User, string, *ErrorType) {

	u,err := interactors.UserRepo.SelectUserByEmail(email)
	if err != nil {
		LogError(err)
		return nil,"",GetError(0)
	}

	validUser := comparePasswords(u.Token, []byte(password))
	if !validUser{

		return  nil, "", GetError(0)
	}

	// Create the token
	token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))

	// Set some claims
	token.Claims["exp"] = time.Now().Add(time.Hour * 24 ).Unix()
	token.Claims["userId"] = u.Id
	token.Claims["email"] = u.Email

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil{
		LogError(err)
		return  nil, "", GetError(0)
	}

	return u, tokenString, nil
}


func (DashboardInteractor) FillProductTable() *ErrorType{

	xlsx, err := excelize.OpenFile("./list.xlsx")
	if err != nil {
		fmt.Println(err)
		return GetError(0)
	}
	rows := xlsx.GetRows("Sheet1")
	timeNow := int(time.Now().Unix())
	for _, row := range rows {
		p := &Product{
			Barcode:row[0],
			Name:row[1],
			RegisterDate:timeNow,
		}

		err := interactors.ProductRepo.InsertProduct(p)
		if err != nil {
			LogError(err)
			return GetError(0)
		}
	}

	return nil

}

func (DashboardInteractor) CreateProduct(p *Product) *ErrorType{

	p.RegisterDate = int(time.Now().Unix())
	err := interactors.ProductRepo.InsertProduct(p)
	if err != nil{
		LogError(err)
		//return GetError(0)
	}
	return nil

}

func (DashboardInteractor) UpdateProduct(p *Product) *ErrorType{

	p.RegisterDate = int(time.Now().Unix())
	err := interactors.ProductRepo.UpdateProductById(p,p.Id)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) GetProductById(id int) (*Product,*ErrorType){

	p,err := interactors.ProductRepo.SelectProductById(id)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

func (DashboardInteractor) GetProducts(barcode,name,description,category,orderBy,orderAs string,pageNumber, pageSize int) (*responses.ProductResponse,  *ErrorType){

	p,err := interactors.ProductRepo.SelectProducts(barcode,name,description,category,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

func (DashboardInteractor) DeleteProducts(ids []int) *ErrorType{

	err := interactors.ProductRepo.DeleteProducts(ids)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil
}

// ########################################################################

func (DashboardInteractor) CreateStock(p *Stock) *ErrorType{

	p.CreationDate = int(time.Now().Unix())
	p.UpdateDate = int(time.Now().Unix())
	err := interactors.StockRepo.InsertStock(p)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) UpdateStock(p *Stock) *ErrorType{

	p.UpdateDate = int(time.Now().Unix())
	err := interactors.StockRepo.UpdateStockById(p,p.Id)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) GetStockById(id int) (*Stock,*ErrorType){

	p,err := interactors.StockRepo.SelectStockById(id)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

func (DashboardInteractor) GetStocks(barcode,name,description,category,orderBy,orderAs string,pageNumber, pageSize,dealerId int) (*responses.StockResponse,  *ErrorType){

	p,err := interactors.StockRepo.SelectStocks(barcode,name,description,category,orderBy,orderAs,pageNumber, pageSize,dealerId)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil
}

func (DashboardInteractor) DeleteStocks(ids []int) *ErrorType{

	err := interactors.StockRepo.DeleteStocks(ids)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil
}

// ###################################################################

func (DashboardInteractor) CreatePerson(p *Person) *ErrorType{

	p.CreationDate = int(time.Now().Unix())
	err := interactors.PersonRepo.InsertPerson(p)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) UpdatePerson(p *Person) *ErrorType{

	p.CreationDate = int(time.Now().Unix())
	err := interactors.PersonRepo.UpdatePersonById(p,p.Id)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) GetPersonById(id int) (*Person,*ErrorType){

	p,err := interactors.PersonRepo.SelectPersonById(id)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

func (DashboardInteractor) GetPeople(name,pType,orderBy,orderAs string,pageNumber, pageSize int) (*responses.PersonResponse,  *ErrorType){

	p,err := interactors.PersonRepo.SelectPeople(name,pType,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

func (DashboardInteractor) DeletePersons(ids []int) *ErrorType{

	err := interactors.PersonRepo.DeletePersons(ids)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil
}

// #################################################################

func (DashboardInteractor) CreateReceiving(p *Receiving) *ErrorType{

	p.CreationDate = int(time.Now().Unix())
	p.UpdateDate = int(time.Now().Unix())
	err := interactors.ReceivingRepo.InsertReceiving(p)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) UpdateReceiving(p *Receiving) *ErrorType{

	p.UpdateDate = int(time.Now().Unix())
	err := interactors.ReceivingRepo.UpdateReceivingById(p,p.Id)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) GetReceivingById(id int) (*Receiving,*ErrorType){

	p,err := interactors.ReceivingRepo.SelectReceivingById(id)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

func (DashboardInteractor) GetReceivings(person,status,orderBy,orderAs string,pageNumber, pageSize int) (*responses.ReceivingResponse,  *ErrorType){

	p,err := interactors.ReceivingRepo.SelectReceivings(person,status,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}

	// insert product details to productList
	for k,item := range p.Items{
		var prodInt []int
		prodStr := strings.Split(item.ProductIds,",")
		for _,v := range prodStr{
			vInt,_ := strconv.Atoi(v)
			prodInt = append(prodInt,vInt)
		}
		for _,v := range prodInt{
			product,err := interactors.ProductRepo.SelectProductById(v)
			if err != nil{
				LogError(err)
				return nil,GetError(0)
			}

			p.Items[k].ProductList = append(p.Items[k].ProductList,product)
		}

	}

	return p,nil

}

func (DashboardInteractor) SetReceivingStatus(status string,id int) *ErrorType{

	err := interactors.ReceivingRepo.SetStatus(status,id)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) DeleteReceivings(ids []int) *ErrorType{

	err := interactors.ReceivingRepo.DeleteReceivings(ids)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil
}

// ###################################################################

func (DashboardInteractor) CreatePayment(p *Payment) *ErrorType{

	p.CreationDate = int(time.Now().Unix())
	p.UpdateDate = int(time.Now().Unix())
	err := interactors.PaymentRepo.InsertPayment(p)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) UpdatePayment(p *Payment) *ErrorType{

	p.UpdateDate = int(time.Now().Unix())
	err := interactors.PaymentRepo.UpdatePaymentById(p,p.Id)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) GetPaymentById(id int) (*Payment,*ErrorType){

	p,err := interactors.PaymentRepo.SelectPaymentById(id)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

func (DashboardInteractor) GetPayments(person,status,orderBy,orderAs string,pageNumber, pageSize int) (*responses.PaymentResponse,  *ErrorType){

	p,err := interactors.PaymentRepo.SelectPayments(person,status,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

func (DashboardInteractor) DeletePayments(ids []int) *ErrorType{

	err := interactors.PaymentRepo.DeletePayments(ids)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil
}

func (DashboardInteractor) SetPaymentStatus(status string,id int) *ErrorType{

	err := interactors.PaymentRepo.SetPaymentStatus(status,id)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

// ###########################################################

func (DashboardInteractor) CreateExpense(p *Expense) *ErrorType{

	p.CreateDate = int(time.Now().Unix())
	p.UpdateDate = int(time.Now().Unix())
	err := interactors.ExpenseRepo.InsertExpense(p)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) UpdateExpense(p *Expense) *ErrorType{

	p.UpdateDate = int(time.Now().Unix())
	err := interactors.ExpenseRepo.UpdateExpenseById(p,p.Id)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) GetExpenseById(id int) (*Expense,*ErrorType){

	p,err := interactors.ExpenseRepo.SelectExpenseById(id)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

func (DashboardInteractor) GetExpenses(name,description,orderBy,orderAs string,pageNumber, pageSize int) (*responses.ExpenseResponse,  *ErrorType){

	p,err := interactors.ExpenseRepo.SelectExpenses(name,description,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

func (DashboardInteractor) DeleteExpenses(ids []int) *ErrorType{

	err := interactors.ExpenseRepo.DeleteExpenses(ids)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil
}

// ###################################################################

func (DashboardInteractor) CreateUser(p *User) *ErrorType{

	p.RegisterDate = int(time.Now().Unix())

	p.Token = hashAndSalt([]byte(p.Password))

	//LogDebug(string(p.Token))
	err := interactors.UserRepo.InsertUser(p)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) UpdateUser(p *User) *ErrorType{

	p.Token = hashAndSalt([]byte(p.Password))

	err := interactors.UserRepo.UpdateUserById(p,p.Id)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) GetUserById(id int) (*User,*ErrorType){

	p,err := interactors.UserRepo.SelectUserById(id)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

func (DashboardInteractor) GetUsers(name,email,orderBy,orderAs string,pageNumber, pageSize int) (*responses.UserResponse,  *ErrorType){

	p,err := interactors.UserRepo.SelectUsers(name,email,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

func (DashboardInteractor) DeleteUsers(ids []int) *ErrorType{

	err := interactors.UserRepo.DeleteUsers(ids)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil
}


// ###########################################################

func (DashboardInteractor) CreateSale(p *Sale) *ErrorType{

	p.CreationDate = int(time.Now().Unix())
	err := interactors.SaleRepo.InsertSale(p)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) UpdateSale(p *Sale) *ErrorType{

	p.CreationDate = int(time.Now().Unix())
	err := interactors.SaleRepo.UpdateSaleById(p,p.Id)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) GetSaleById(id int) (*Sale,*ErrorType){

	p,err := interactors.SaleRepo.SelectSaleById(id)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

func (DashboardInteractor) GetSales(tInterval,orderBy,orderAs string,pageNumber, pageSize int) (*responses.SaleResponse,  *ErrorType){

	strInter := strings.Split(tInterval,",")
	intInter := []int{}
	for _,str := range strInter{
		i,_ := strconv.Atoi(str)
		intInter = append(intInter,i)
	}

	p,err := interactors.SaleRepo.SelectSales(intInter,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

func (DashboardInteractor) DeleteSales(ids []int) *ErrorType{

	err := interactors.SaleRepo.DeleteSales(ids)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil
}

// ###################################################################

func GetError(code int) (*ErrorType){
	return &ErrorType{Code: code, Message: errorMap["tr"][code]}
}