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
	"sort"
	"encoding/json"
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


func (DashboardInteractor) FillProductTable(userId int) *ErrorType{

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
			UserId:userId,
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

func (DashboardInteractor) RetrieveCategories()([]string,*ErrorType){

	p,err := interactors.ProductRepo.SelectProductCategories()
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

// ########################################################################

func (DashboardInteractor) CreateStock(p *Stock) *ErrorType{

	product,err := interactors.StockRepo.SelectStockByProductId(p.ProductId)
	if product != nil {
		temp := &Stock{
			Id:			product.Id,
			ProductId:	p.ProductId,
			Qty: 		product.Qty+p.Qty, // add up the quantity of product if it is already exist
			UpdateDate: int(time.Now().Unix()),
			UserId:		p.UserId,
			DealerId:	p.DealerId,
			IsFavorite: p.IsFavorite,
		}

		err := interactors.StockRepo.UpdateStockById(temp,temp.Id)
		if err != nil {
			LogError(err)
			return GetError(0)
		}

	}else if product == nil {
		p.CreationDate = int(time.Now().Unix())
		p.UpdateDate = int(time.Now().Unix())
		err = interactors.StockRepo.InsertStock(p)
		if err != nil{
			LogError(err)
			return GetError(0)
		}
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

func (DashboardInteractor) GetStocks(tInterval string,barcode,name,description,category,orderBy,orderAs string,pageNumber, pageSize,dealerId, userId int,isFavorite bool) (*responses.StockResponse,  *ErrorType){

	strInter := strings.Split(tInterval,",")
	intInter := []int{}
	for _,str := range strInter{
		i,_ := strconv.Atoi(str)
		intInter = append(intInter,i)
	}

	p,err := interactors.StockRepo.SelectStocks(intInter,barcode,name,description,category,orderBy,orderAs,pageNumber, pageSize,dealerId,userId,isFavorite)
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

func (DashboardInteractor) SetFavoriteProduct(productId int, isFavorite bool) *ErrorType{

	err := interactors.StockRepo.SetFavoriteByProductId(productId,isFavorite)
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

func (DashboardInteractor) GetReceivings(tInterval string,person,status,orderBy,orderAs string,pageNumber, pageSize int,creator int) (*responses.ReceivingResponse,  *ErrorType){

	strInter := strings.Split(tInterval,",")
	intInter := []int{}
	for _,str := range strInter{
		i,_ := strconv.Atoi(str)
		intInter = append(intInter,i)
	}

	p,err := interactors.ReceivingRepo.SelectReceivings(intInter,person,status,orderBy,orderAs,pageNumber, pageSize,creator)
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

func (DashboardInteractor) GetPayments(tInterval string,person,status,orderBy,orderAs string,pageNumber, pageSize,creator int) (*responses.PaymentResponse,  *ErrorType){

	strInter := strings.Split(tInterval,",")
	intInter := []int{}
	for _,str := range strInter{
		i,_ := strconv.Atoi(str)
		intInter = append(intInter,i)
	}
	p,err := interactors.PaymentRepo.SelectPayments(intInter,person,status,orderBy,orderAs,pageNumber, pageSize,creator)
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

func (DashboardInteractor) GetExpenses(tInterval string,name,description,orderBy,orderAs string,pageNumber, pageSize int,creator int) (*responses.ExpenseResponse,  *ErrorType){

	strInter := strings.Split(tInterval,",")
	intInter := []int{}
	for _,str := range strInter{
		i,_ := strconv.Atoi(str)
		intInter = append(intInter,i)
	}

	p,err := interactors.ExpenseRepo.SelectExpenses(intInter,name,description,orderBy,orderAs,pageNumber, pageSize,creator)
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

func (d DashboardInteractor) CreateSaleBasket(p *SaleBasket) *ErrorType{

	p.CreationDate = int(time.Now().Unix())
	err := interactors.SaleBasketRepo.InsertSaleBasket(p)
	if err != nil{
		LogError(err)
		return GetError(0)
	}

	go d.CreateSaleDetail(p) // go routine

	return nil

}

func (DashboardInteractor) CreateSaleDetail(p *SaleBasket) {

	ItemStr := p.ItemsStr

	var basket []SaleBasketItem
	timeNow := int(time.Now().Unix())

	if err := json.Unmarshal([]byte(ItemStr),&basket); err != nil{
		LogError(err)
	}

	//var reportEntry SaleSummaryReport

	for _,v := range basket{

		sDetail := &SaleDetail{
			CreationDate:timeNow,
			BasketId:p.Id,
			ProductId:v.Id,
			Qty:v.Qty,
			Discount:v.Discount,
			UserId:p.UserId,
		}

		if err := interactors.SaleDetailRepo.InsertSaleDetail(sDetail); err != nil{
			LogError(err)
			// TODO: rollback
		}

		if err := interactors.StockRepo.DecrementProductFromStock(v.Id,v.Qty); err != nil{
			LogError(err)
			// TODO: rollback
		}
	}

}

func (DashboardInteractor) UpdateSaleBasket(p *SaleBasket) *ErrorType{

	p.CreationDate = int(time.Now().Unix())
	err := interactors.SaleBasketRepo.UpdateSaleBasketById(p,p.Id)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil

}

func (DashboardInteractor) GetSaleBasketById(id int) (*SaleBasket,*ErrorType){

	p,err := interactors.SaleBasketRepo.SelectSaleBasketById(id)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

func (DashboardInteractor) GetSaleBaskets(tInterval string,userId int,orderBy,orderAs string,pageNumber, pageSize int) (*responses.SaleBasketResponse,  *ErrorType){

	strInter := strings.Split(tInterval,",")
	intInter := []int{}
	for _,str := range strInter{
		i,_ := strconv.Atoi(str)
		intInter = append(intInter,i)
	}

	p,err := interactors.SaleBasketRepo.SelectSaleBaskets(intInter,userId,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}
	return p,nil

}

func (DashboardInteractor) DeleteSaleBaskets(ids []int) *ErrorType{

	err := interactors.SaleBasketRepo.DeleteSaleBaskets(ids)
	if err != nil{
		LogError(err)
		return GetError(0)
	}
	return nil
}

// ###################################################################

// # Reports #

func (DashboardInteractor) GetSaleSummaryReport(tInterval string) (*SaleSummaryReport,  *ErrorType){

	strInter := strings.Split(tInterval,",")
	intInter := []int{}
	for _,str := range strInter{
		i,_ := strconv.Atoi(str)
		intInter = append(intInter,i)
	}

	p := &SaleSummaryReport{}

	items,err := interactors.SaleSummaryReportRepo.SelectSaleSummaryReportItems(intInter)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}

	p.AsObject = items.AsObject

	for _,i := range items.AsObject.Items{

		p.GrossProfit += i.GrossProfit
		p.GrossProfits = append(p.GrossProfits, i.GrossProfit)

		p.NetProfit += i.NetProfit
		p.NetProfits = append(p.NetProfits, i.NetProfit)

		p.SaleCount += i.SaleCount
		p.SaleCounts = append(p.SaleCounts, i.SaleCount)

		p.ItemCount += i.ItemCount
		p.ItemCounts = append(p.ItemCounts, i.ItemCount)

		p.CustomerCount += i.CustomerCount
		p.CustomerCounts = append(p.CustomerCounts, i.CustomerCount)

		p.Discount += i.Discount
		p.Discounts = append(p.Discounts, i.Discount)

		p.BasketValue += i.BasketValue
		p.BasketValues = append(p.BasketValues, i.BasketValue)

		p.BasketSize += i.BasketSize
		p.BasketSizes = append(p.BasketSizes, i.BasketSize)

		p.Timestamps = append(p.Timestamps, i.Timestamp)

	}

	p.BasketValue = p.GrossProfit / float64(p.SaleCount)


	return p,nil
}

func (DashboardInteractor) GetCurrentStockReport(name,category,orderBy,orderAs string,pageNumber, pageSize int) (*responses.CurrentStockReportResponse,  *ErrorType){

	p,err := interactors.StockRepo.SelectCurrentStockReport(name,category,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}

	p.Total.Name = "Total"
	for _,v := range p.Items {
		p.Total.Qty += v.Qty
		p.Total.PurchasePrice += v.PurchasePrice
		p.Total.SalePrice += v.SalePrice
		p.Total.GrossValue += v.GrossValue
		p.Total.NetValue += v.NetValue
		p.Total.TotalProfit += v.TotalProfit
	}

	return p,nil
}

type saleInstance struct {
	Barcode string `json:"barcode"`
	Qty int			`json:"qty"`
}

func (DashboardInteractor) GetActivityLog(tInterval string,userId int)(*ActivityLogs,*ErrorType){

	res := []*ActivityLogItem{}

	strInter := strings.Split(tInterval,",")
	intInter := []int{}
	for _,str := range strInter{
		i,_ := strconv.Atoi(str)
		intInter = append(intInter,i)
	}

	// # get sales
	sales,err := interactors.SaleBasketRepo.SelectSaleBaskets(intInter,userId,"","",0,0)
	if err != nil {
		LogError(err)
		return nil,GetError(0)
	}




	for _,v := range sales.Items {

		text := ``

		var saleInst []SaleBasketItem
		err := json.Unmarshal([]byte(v.ItemsStr),&saleInst)
		if err != nil {
			LogError(err)
		}

		for _,vv := range saleInst{
			product,err := interactors.ProductRepo.SelectProductById(vv.Id)
			if err != nil {
				LogError(err)
			}
			text += product.Name + `: ` + strconv.Itoa(vv.Qty) + ` piece `
			text += ` `

		}

		temp := &ActivityLogItem{
			User:v.UserName,
			Date:v.CreationDate,
			ActivityType:"Sale",
			Description:text,
			Title: "Sale",
		}

		res = append(res,temp)

	}

	// # get stock entries
	stocks ,err := interactors.StockRepo.SelectStocks(intInter,"","","","","","",0,0,0,userId,false)
	if err != nil {
		LogError(err)
		return nil,GetError(0)
	}

	for _,v := range stocks.Items {

		text := strconv.Itoa(v.Qty) + ` piece ` + v.Product.Name
		detail := "Dealer:" + v.DealerName
		temp := &ActivityLogItem{
			User:v.UserName,
			Date:v.UpdateDate,
			ActivityType:"Stock",
			Description:text,
			Detail:detail,
			Title: "Stock Entry",

		}

		res = append(res,temp)

	}

	// # getPayments
	payments,err := interactors.PaymentRepo.SelectPayments(intInter,"","","","",0,0,userId)
	if err != nil {
		LogError(err)
		return nil,GetError(0)
	}

	for _,v := range payments.Items{
		text := v.PersonName + ` named operator, ` + strconv.FormatFloat(v.Amount,'f',2,64) + ` amount.`
		detail := `Payment Date: ` + time.Unix(int64(v.ExpectedDate),0).Format("2016-01-02 15:04:05")
		temp := &ActivityLogItem{
			User:v.UserName,
			Date:v.UpdateDate,
			ActivityType:"Payment",
			Description:text,
			Detail:detail,
			Title: "Payment Entry",

		}

		res = append(res,temp)
	}

	// # getReceivings
	receivings,err := interactors.ReceivingRepo.SelectReceivings(intInter,"","","","",0,0,userId)
	if err != nil {
		LogError(err)
		return nil,GetError(0)
	}

	for _,v := range receivings.Items{

		text := v.PersonName + ` named operator, ` + strconv.FormatFloat(v.Amount,'f',2,64) + ` amount.`
		detail := `Payment Date: ` + time.Unix(int64(v.ExpectedDate),0).Format("2016-01-02 15:04:05")
		temp := &ActivityLogItem{
			User:v.UserName,
			Date:v.UpdateDate,
			ActivityType:"Receiving",
			Description:text,
			Detail:detail,
			Title: "Receiving Entry",

		}

		res = append(res,temp)

	}

	// getExpense
	expenses,err := interactors.ExpenseRepo.SelectExpenses(intInter,"","","","",0,0,userId)
	if err != nil {
		LogError(err)
		return nil,GetError(0)
	}

	for _,v := range expenses.Items{

		text := `'` + v.Name + ` named person,` + strconv.FormatFloat(v.Price,'f',2,64) + ` amount.`
		//detail := `Ödenme Tarihi: ` + time.Unix(int64(v.ExpectedDate),0).Format("2016-01-02 15:04:05")
		temp := &ActivityLogItem{
			User:v.UserName,
			Date:v.UpdateDate,
			ActivityType:"Expense",
			Description:text,
			Title: "Expense Entry",
			//Detail:detail,
		}

		res = append(res,temp)


	}

	// # sort by date ASC
	sort.Slice(res,func(i,j int) bool {
		return res[i].Date < res[j].Date
	})

	result := &ActivityLogs{
		Items:res,
	}
	result.Count = len(res)

	return result,nil
}

func (DashboardInteractor) GetPaymentReport(tInterval string) (*PaymentReport,  *ErrorType){

	strInter := strings.Split(tInterval,",")
	intInter := []int{}
	for _,str := range strInter{
		i,_ := strconv.Atoi(str)
		intInter = append(intInter,i)
	}

	var timestamp []string
	firstDay := intInter[0]
	lastDay := intInter[1]

	firstDayTimeFormat := time.Unix(int64(firstDay),0)
	//lastDayTimeFormat := time.Unix(int64(lastDay),0)

	for timeIterator:= firstDayTimeFormat ; timeIterator.Unix() < int64(lastDay) ; {
		timeIteratorStr := timeIterator.Format("02/01") // DD/MM format
		timestamp = append(timestamp,timeIteratorStr)
		timeIterator = timeIterator.Add(24 *time.Hour)
	}

	result := &PaymentReport{}
	result.Timestamps = timestamp

	paymentsList := make([]float64, len(timestamp))
	expensesList := make([]float64, len(timestamp))
	receivingsList := make([]float64, len(timestamp))

	// # check payments
	payments,err := interactors.PaymentRepo.SelectPayments(intInter,"","","","",0,0,0)
	if err != nil {
		LogError(err)
		return nil,GetError(0)
	}

	for _,v := range payments.Items{

		if v.Status == "Bitti"{

			// find the index number
			expectedTimeStr := time.Unix(int64(v.ExpectedDate),0).Format("02/01")
			var index int
			for k,v := range timestamp {
				if v == expectedTimeStr{
					index = k
				}
			}

			result.TotalPayments += v.Amount
			paymentsList[index] += v.Amount

		}else if v.Status == "Gecikmiş" {
			result.OverduePayments += 1
		}

		paymentItem := &PaymentList{
			Person:v.PersonName,
			Amount:v.Amount,
			Timestamp:v.ExpectedDate,
			Status:v.Status,
			Detail:v.Summary,
			Type: "Payment",
		}

		if paymentItem.Status == "Gecikmiş" {
			paymentItem.Status = "Overdue"
		}else if paymentItem.Status == "Bitti"{
			paymentItem.Status = "Finished"
		}else if paymentItem.Status == "Bekliyor" {
			paymentItem.Status = "Pending"
		}

		result.ItemsAsObject = append(result.ItemsAsObject,paymentItem)

	}

	//result.Payments = append(result.Payments,paymentsList...)
	result.Payments = paymentsList


	// # Receivings
	receivings,err := interactors.ReceivingRepo.SelectReceivings(intInter,"","","","",0,0,0)
	if err != nil {
		LogError(err)
		return nil,GetError(0)
	}

	for _,v := range receivings.Items{

		if v.Status == "Bitti"{
			// find the index number
			expectedTimeStr := time.Unix(int64(v.ExpectedDate),0).Format("02/01")
			var index2 int
			for k,v := range timestamp {
				if v == expectedTimeStr{
					index2 = k
				}
			}

			result.TotalReceivings += v.Amount
			receivingsList[index2] += v.Amount

		}else if v.Status == "Gecikmiş" {
			result.OverdueReceivings += 1
		}

		paymentItem := &PaymentList{
			Person:v.PersonName,
			Amount:v.Amount,
			Timestamp:v.ExpectedDate,
			Status:v.Status,
			Detail:"Receiving",
			Type: "Receiving",
		}

		if paymentItem.Status == "Gecikmiş" {
			paymentItem.Status = "Overdue"
		}else if paymentItem.Status == "Bitti"{
			paymentItem.Status = "Finished"
		}else if paymentItem.Status == "Bekliyor" {
			paymentItem.Status = "Pending"
		}


		result.ItemsAsObject = append(result.ItemsAsObject,paymentItem)
	}

	result.Receivings = receivingsList


	// # expenses
	expenses,err := interactors.ExpenseRepo.SelectExpenses(intInter,"","","","",0,0,0)
	if err != nil {
		LogError(err)
		return nil,GetError(0)
	}

	for _,v := range expenses.Items{

		// find the index number
		expectedTimeStr := time.Unix(int64(v.UpdateDate),0).Format("02/01")
		var index3 int
		for k,v := range timestamp {
			if v == expectedTimeStr{
				index3 = k
			}
		}

		result.TotalExpenses += v.Price
		expensesList[index3] += v.Price

		paymentItem := &PaymentList{
			Person:v.UserName,
			Amount:v.Price,
			Timestamp:v.UpdateDate,
			Status:"Finished",
			Detail:v.Name,
			Type: "Expense",
		}

		result.ItemsAsObject = append(result.ItemsAsObject,paymentItem)
	}

	sort.Slice(result.ItemsAsObject, func(i, j int) bool {
		return result.ItemsAsObject[i].Timestamp < result.ItemsAsObject[j].Timestamp
	})

	result.Expenses = expensesList

	return result,nil
}

func (DashboardInteractor) GetProductReport(tInterval string,productName string,category string,userId int) (*ProductReport,  *ErrorType){

	strInter := strings.Split(tInterval,",")
	intInter := []int{}
	for _,str := range strInter{
		i,_ := strconv.Atoi(str)
		intInter = append(intInter,i)
	}

	res,err := interactors.SaleDetailRepo.SelectSaleDetails(intInter,productName,category,userId)
	if err != nil{
		LogError(err)
		return nil,GetError(0)
	}

	// TODO : profit percentage sürüm hesabı olmadı
	// profil percentage'i sql den kaldırabiliriz
	for k,v := range res.Items{
		res.Items[k].ProfitPercentage = v.NetProfit / (res.TotalNetProfit) * 100
	}

	return res,nil
}
// ########################################################################

// # Reports To Excel

func (DashboardInteractor) GetSaleSummaryReportAsExcel(tInterval string) (string,  *ErrorType){

	strInter := strings.Split(tInterval,",")
	intInter := []int{}
	for _,str := range strInter{
		i,_ := strconv.Atoi(str)
		intInter = append(intInter,i)
	}

	p := &SaleSummaryReport{}

	items,err := interactors.SaleSummaryReportRepo.SelectSaleSummaryReportItems(intInter)
	if err != nil{
		LogError(err)
		return "",GetError(0)
	}

	p.AsObject = items.AsObject

	for _,i := range items.AsObject.Items{

		p.GrossProfit += i.GrossProfit
		p.NetProfit += i.NetProfit
		p.SaleCount += i.SaleCount
		p.ItemCount += i.ItemCount
		p.CustomerCount += i.CustomerCount
		p.Discount += i.Discount
		p.BasketValue += i.BasketValue
		p.BasketSize += i.BasketSize

	}

	fileName := "Sale-Summary-Report"
	filePath := "./excelFiles/" + fileName + ".xlsx"
	SaveSaleReportAsExcelFile(p,fileName)

	return filePath,nil
}

func (DashboardInteractor) GetCurrentStockReportAsExcel(name,category,orderBy,orderAs string,pageNumber, pageSize int) (string,  *ErrorType){

	p,err := interactors.StockRepo.SelectCurrentStockReport(name,category,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		LogError(err)
		return "",GetError(0)
	}
	p.Total.Name = "Total"
	for _,v := range p.Items {
		p.Total.Qty += v.Qty
		p.Total.PurchasePrice += v.PurchasePrice
		p.Total.SalePrice += v.SalePrice
		p.Total.GrossValue += v.GrossValue
		p.Total.NetValue += v.NetValue
		p.Total.TotalProfit += v.TotalProfit
	}

	fileName := "Current-Stock-Report"
	filePath := "./excelFiles/" + fileName + ".xlsx"
	SaveCurrentStockReportAsExcelFile(p, fileName)

	return filePath,nil
}

func (DashboardInteractor) GetPaymentReportAsExcel(tInterval string) (string,  *ErrorType){

	strInter := strings.Split(tInterval,",")
	intInter := []int{}
	for _,str := range strInter{
		i,_ := strconv.Atoi(str)
		intInter = append(intInter,i)
	}

	var timestamp []string
	firstDay := intInter[0]
	lastDay := intInter[1]

	firstDayTimeFormat := time.Unix(int64(firstDay),0)
	//lastDayTimeFormat := time.Unix(int64(lastDay),0)

	for timeIterator:= firstDayTimeFormat ; timeIterator.Unix() < int64(lastDay) ; {
		timeIteratorStr := timeIterator.Format("02/01") // DD/MM format
		timestamp = append(timestamp,timeIteratorStr)
		timeIterator = timeIterator.Add(24 *time.Hour)
	}

	result := &PaymentReport{}
	result.Timestamps = timestamp

	paymentsList := make([]float64, len(timestamp))
	expensesList := make([]float64, len(timestamp))
	receivingsList := make([]float64, len(timestamp))

	// # check payments
	payments,err := interactors.PaymentRepo.SelectPayments(intInter,"","","","",0,0,0)
	if err != nil {
		LogError(err)
		return "",GetError(0)
	}

	for _,v := range payments.Items{

		if v.Status == "Bitti"{

			// find the index number
			expectedTimeStr := time.Unix(int64(v.ExpectedDate),0).Format("02/01")
			var index int
			for k,v := range timestamp {
				if v == expectedTimeStr{
					index = k
				}
			}

			result.TotalPayments += v.Amount
			paymentsList[index] += v.Amount

		}else if v.Status == "Gecikmiş" {
			result.OverduePayments += 1
		}

		paymentItem := &PaymentList{
			Person:v.PersonName,
			Amount:v.Amount,
			Timestamp:v.ExpectedDate,
			Status:v.Status,
			Detail:v.Summary,
			Type: "Ödeme",
		}

		result.ItemsAsObject = append(result.ItemsAsObject,paymentItem)

	}

	//result.Payments = append(result.Payments,paymentsList...)
	result.Payments = paymentsList


	// # Receivings
	receivings,err := interactors.ReceivingRepo.SelectReceivings(intInter,"","","","",0,0,0)
	if err != nil {
		LogError(err)
		return "",GetError(0)
	}

	for _,v := range receivings.Items{

		if v.Status == "Bitti"{
			// find the index number
			expectedTimeStr := time.Unix(int64(v.ExpectedDate),0).Format("02/01")
			var index2 int
			for k,v := range timestamp {
				if v == expectedTimeStr{
					index2 = k
				}
			}

			result.TotalReceivings += v.Amount
			receivingsList[index2] += v.Amount

		}else if v.Status == "Gecikmiş" {
			result.OverdueReceivings += 1
		}

		paymentItem := &PaymentList{
			Person:v.PersonName,
			Amount:v.Amount,
			Timestamp:v.ExpectedDate,
			Status:v.Status,
			Detail:"Tahsilat",
			Type: "Tahsilat",
		}

		result.ItemsAsObject = append(result.ItemsAsObject,paymentItem)
	}

	result.Receivings = receivingsList


	// # expenses
	expenses,err := interactors.ExpenseRepo.SelectExpenses(intInter,"","","","",0,0,0)
	if err != nil {
		LogError(err)
		return "",GetError(0)
	}

	for _,v := range expenses.Items{

		// find the index number
		expectedTimeStr := time.Unix(int64(v.UpdateDate),0).Format("02/01")
		var index3 int
		for k,v := range timestamp {
			if v == expectedTimeStr{
				index3 = k
			}
		}

		result.TotalExpenses += v.Price
		expensesList[index3] += v.Price

		paymentItem := &PaymentList{
			Person:v.UserName,
			Amount:v.Price,
			Timestamp:v.UpdateDate,
			Status:"Bitti",
			Detail:v.Name,
			Type: "Harcama",
		}

		result.ItemsAsObject = append(result.ItemsAsObject,paymentItem)
	}

	sort.Slice(result.ItemsAsObject, func(i, j int) bool {
		return result.ItemsAsObject[i].Timestamp < result.ItemsAsObject[j].Timestamp
	})

	result.Expenses = expensesList

	fileName := "Payment-Report"
	filePath := "./excelFiles/" + fileName + ".xlsx"
	SavePaymentReportAsExcelFile(result, fileName)

	return filePath,nil
}

func (d DashboardInteractor) GetProductReportAsExcel(tInterval string,productName string,category string,userId int) (string,  *ErrorType){

	result,err := d.GetProductReport(tInterval,productName,category,userId)
	if err != nil {
		LogError(err)
		return "",err
	}

	fileName := "Current-Product-Report"
	filePath := "./excelFiles/" + fileName + ".xlsx"
	SaveProductReportAsExcelFile(result, fileName)

	return filePath,nil
}

func SaveProductReportAsExcelFile(p *ProductReport,fileName string){
	file := excelize.NewFile()
	// Create a new sheet.
	//index := file.NewSheet("Sheet1")
	// Set value of a cell.
	var cols []string
	cols = append(cols, "No")
	cols = append(cols, "İsim")
	cols = append(cols, "Adet")
	cols = append(cols, "Brüt Kar")
	cols = append(cols, "Net Kar")
	cols = append(cols, "İndirim")
	cols = append(cols, "Kar Marjı (%)")
	cols = append(cols, "Kar Dilimi (%)")
	cols = append(cols, "İade Sayısı (%)")

	//print titles
	rowIndex := 1
	colIndex := 'A'
	for j:=0; j < len(cols); j++ {

		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex), cols[j])
		colIndex++
	}

	//print values
	colIndex = 'A'
	rowIndex += 1
	for k,v := range p.Items {
		rowIndex = k + 2
		colIndex = 'A'

		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),k+1)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.ProductName)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.Qty)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.GrossProfit)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.NetProfit)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.Discount)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.Markup)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.ProfitPercentage)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.NumberOfReturn)
		colIndex++
	}

	rowIndex += 1
	colIndex = 'A'

	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),"Toplam")
	colIndex++
	//productName
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),p.TotalQty)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),p.TotalGrossProfit)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),p.TotalNetProfit)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),p.TotalDiscount)
	colIndex++
	// markup
	colIndex++
	// profitPercentage
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),p.TotalNumberOfReturn)
	colIndex++

	err := file.SaveAs("./excelFiles/" + fileName + ".xlsx")
	if err != nil {
		LogError(err)
	}

}


func SaveSaleReportAsExcelFile(saleReport *SaleSummaryReport, fileName string) {
	file := excelize.NewFile()
	// Create a new sheet.
	//index := file.NewSheet("Sheet1")
	// Set value of a cell.
	var cols []string
	cols = append(cols, "No")
	cols = append(cols, "Tarih")
	cols = append(cols, "Gross Profit")
	cols = append(cols, "Net Profit")
	cols = append(cols, "Sale Count")
	cols = append(cols, "Item Count")
	cols = append(cols, "Customer Count")
	cols = append(cols, "Discount")
	cols = append(cols, "Basket Value")
	cols = append(cols, "Basket Size")



	//print titles
	rowIndex := 1
	colIndex := 'A'
	for j:=0; j < len(cols); j++ {

		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex), cols[j])
		colIndex++
	}

	//print values
	colIndex = 'A'
	rowIndex += 1
	for k,v := range saleReport.AsObject.Items {
		rowIndex = k + 2
		colIndex = 'A'

		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),k+1)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),time.Unix(int64(v.Timestamp),0).Format("02/01/2006 03:04"))
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.GrossProfit)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.NetProfit)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.SaleCount)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.ItemCount)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.CustomerCount)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.Discount)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.BasketValue)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.BasketSize)
		colIndex++

	}

	rowIndex += 1
	colIndex = 'A'

	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),"Toplam")
	colIndex++
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),saleReport.GrossProfit)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),saleReport.NetProfit)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),saleReport.SaleCount)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),saleReport.ItemCount)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),saleReport.CustomerCount)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),saleReport.Discount)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),saleReport.BasketValue)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),saleReport.BasketSize)
	colIndex++

	rowIndex += 1
	colIndex = 'A'

	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),"Ortalama")
	colIndex++
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),saleReport.GrossProfit/float64(saleReport.AsObject.Count))
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),saleReport.NetProfit/float64(saleReport.AsObject.Count))
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),saleReport.SaleCount/saleReport.AsObject.Count)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),saleReport.ItemCount/saleReport.AsObject.Count)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),saleReport.CustomerCount/saleReport.AsObject.Count)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),saleReport.Discount/float64(saleReport.AsObject.Count))
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),saleReport.BasketValue/float64(saleReport.AsObject.Count))
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),saleReport.BasketSize/float64(saleReport.AsObject.Count))
	colIndex++


	err := file.SaveAs("./excelFiles/" + fileName + ".xlsx")
	if err != nil {
		fmt.Println(err)
	}
}

func SaveCurrentStockReportAsExcelFile(stockReport *responses.CurrentStockReportResponse, fileName string) {
	file := excelize.NewFile()
	// Create a new sheet.
	//index := file.NewSheet("Sheet1")
	// Set value of a cell.
	var cols []string
	cols = append(cols, "No")
	cols = append(cols, "Name")
	cols = append(cols, "Category")
	cols = append(cols, "Qty")
	cols = append(cols, "Purchase Price")
	cols = append(cols, "Sale Price")
	cols = append(cols, "Gross Value")
	cols = append(cols, "Net Value")
	cols = append(cols, "Total Profit")


	//print titles
	rowIndex := 1
	colIndex := 'A'
	for j:=0; j < len(cols); j++ {

		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex), cols[j])
		colIndex++
	}

	//print values
	colIndex = 'A'
	rowIndex += 1

	for k,v := range stockReport.Items{
		rowIndex = k + 2
		colIndex = 'A'

		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),k+1)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.Name)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.Category)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.Qty)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.PurchasePrice)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.SalePrice)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.GrossValue)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.NetValue)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.TotalProfit)
		colIndex++

	}

	rowIndex += 1
	colIndex = 'A'

	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),"Total")
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),stockReport.Total.Name)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),stockReport.Total.Category)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),stockReport.Total.Qty)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),stockReport.Total.PurchasePrice)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),stockReport.Total.SalePrice)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),stockReport.Total.GrossValue)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),stockReport.Total.NetValue)
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),stockReport.Total.TotalProfit)
	colIndex++

	err := file.SaveAs("./excelFiles/" + fileName + ".xlsx")
	if err != nil {
		fmt.Println(err)
	}
}

func SavePaymentReportAsExcelFile(paymentReport *PaymentReport, fileName string) {
	file := excelize.NewFile()
	// Create a new sheet.
	//index := file.NewSheet("Sheet1")
	// Set value of a cell.
	var cols []string
	cols = append(cols, "No")
	cols = append(cols, "Timestamp")
	cols = append(cols, "Person")
	cols = append(cols, "Amount")
	cols = append(cols, "Status")
	cols = append(cols, "Detail")
	cols = append(cols, "Type")



	//print titles
	rowIndex := 1
	colIndex := 'A'
	for j:=0; j < len(cols); j++ {

		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex), cols[j])
		colIndex++
	}

	//print values
	colIndex = 'A'
	rowIndex += 1

	for k,v := range paymentReport.ItemsAsObject{
		rowIndex = k + 2
		colIndex = 'A'

		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),k+1)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),time.Unix(int64(v.Timestamp),0).Format("02/01/2006 03:04"))
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.Person)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.Amount)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.Status)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.Detail)
		colIndex++
		file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),v.Type)
		colIndex++

	}

	rowIndex += 2
	colIndex = 'A'

	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),"Özet")
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),"-")
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),"Harcamalar Toplamı")
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),paymentReport.TotalExpenses)
	colIndex++

	rowIndex += 1
	colIndex = 'A'

	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),"Özet")
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),"-")
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),"Ödemeler Toplamı")
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),paymentReport.TotalPayments)
	colIndex++

	rowIndex += 1
	colIndex = 'A'

	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),"Özet")
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),"+")
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),"Tahsilat Toplamı")
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),paymentReport.TotalReceivings)
	colIndex++

	rowIndex += 1
	colIndex = 'A'

	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),"Özet")
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),"₺")
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),"Net Toplam")
	colIndex++
	file.SetCellValue("Sheet1", string(colIndex) + strconv.Itoa(rowIndex),paymentReport.TotalReceivings - paymentReport.TotalPayments - paymentReport.TotalExpenses)
	colIndex++


	err := file.SaveAs("./excelFiles/" + fileName + ".xlsx")
	if err != nil {
		fmt.Println(err)
	}
}

// #########################################################################

// # Utils #

func GetError(code int) (*ErrorType){
	return &ErrorType{Code: code, Message: errorMap["tr"][code]}
}