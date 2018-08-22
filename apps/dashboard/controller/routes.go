package controllers

import (
	"github.com/gin-gonic/gin"
	."stock/entities"
	"strconv"
	"strings"
	"stock/entities/responses"
	. "stock/common/logger"
	"github.com/dgrijalva/jwt-go"
)


func InitRoutes(public,private *gin.RouterGroup) {

	public.POST("login", handleLogin)
	private.GET("me", handleGetMe)
	private.GET("fillProductTable", fillProductTable)

	private.POST("createProduct", createProduct)
	private.POST("updateProduct", updateProduct)
	private.GET("getProductById", getProductById)
	private.GET("deleteProducts", deleteProducts)
	private.GET("getProducts",getProducts)

	private.POST("createStock", createStock)
	private.POST("updateStock", updateStock)
	private.GET("getStockById", getStockById)
	private.GET("deleteStocks", deleteStocks)
	private.GET("getStocks",getStocks)

	private.POST("createPerson", createPerson)
	private.POST("updatePerson", updatePerson)
	private.GET("getPersonById", getPersonById)
	private.GET("deletePeople", deletePeople)
	private.GET("getPeople",getPeople)

	private.POST("createReceiving", createReceiving)
	private.POST("updateReceiving", updateReceiving)
	private.GET("getReceivingById", getReceivingById)
	private.GET("deleteReceivings", deleteReceivings)
	private.GET("getReceivings",getReceivings)
	private.GET("setReceivingStatus", setReceivingStatus)

	private.POST("createPayment", createPayment)
	private.POST("updatePayment", updatePayment)
	private.GET("getPaymentById", getPaymentById)
	private.GET("deletePayments", deletePayments)
	private.GET("getPayments",getPayments)
	private.GET("setPaymentStatus", setPaymentStatus)

	private.POST("createExpense", createExpense)
	private.POST("updateExpense", updateExpense)
	private.GET("getExpenseById", getExpenseById)
	private.GET("deleteExpenses", deleteExpenses)
	private.GET("getExpenses",getExpenses)

	private.POST("createUser", createUser)
	private.POST("updateUser", updateUser)
	private.GET("getUserById", getUserById)
	private.GET("deleteUsers", deleteUsers)
	private.GET("getUsers",getUsers)

	private.POST("createSale", createSale)
	private.POST("updateSale", updateSale)
	private.GET("getSaleById", getSaleById)
	private.GET("deleteSales", deleteSales)
	private.GET("getSales",getSales)


}

func handleGetMe(c *gin.Context){
	var err *ErrorType

	uId := getUserIdFromToken(c)

	ur, err := UseCase.GetUserById(uId)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(LoginResponse{User: ur}))

}

func fillProductTable(c *gin.Context){

	err := UseCase.FillProductTable()
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))

}

func getUserIdFromToken(c *gin.Context) int {
	v, _ := c.Get("token-claims")
	LogDebug(v)
	claims := v.(jwt.MapClaims)
	userId := claims["userId"].(float64)
	return int(userId)
}


func createProduct (c *gin.Context){
	p := Product{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.CreateProduct(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func updateProduct (c *gin.Context){
	p := Product{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.UpdateProduct(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getProductById (c *gin.Context){
	id,_ := strconv.Atoi(c.Query("id"))

	p, err := UseCase.GetProductById(id)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func deleteProducts(c *gin.Context){

	idList := strings.Split(c.Query("ids"),",")

	var ids []int
	for _,id := range idList{
		i,_ := strconv.Atoi(id)
		ids = append(ids,i)
	}

	err := UseCase.DeleteProducts(ids)
	if err != nil {
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

func getProducts (c *gin.Context){

	barcode := c.Query("barcode")
	name := c.Query("name")
	description := c.Query("description")
	category := c.Query("category")

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")
	isDropdown,_ := strconv.ParseBool(c.Query("isDropdown"))

	p, err := UseCase.GetProducts(barcode,name,description,category,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	if isDropdown {
		ResponseList := []*responses.ProductDropdownResponse{}
		for _,product := range p.Items{
			r := &responses.ProductDropdownResponse{}
			r.Id = product.Id
			r.Name = product.Name
			r.Price = product.SalePrice
			ResponseList = append(ResponseList, r)
		}
		c.JSON(200, generateSuccessResponse(ResponseList))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

// ########################################################

func createStock (c *gin.Context){
	p := Stock{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.CreateStock(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func updateStock (c *gin.Context){
	p := Stock{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.UpdateStock(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getStockById (c *gin.Context){
	id,_ := strconv.Atoi(c.Query("id"))

	p, err := UseCase.GetStockById(id)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func deleteStocks(c *gin.Context){

	idList := strings.Split(c.Query("ids"),",")

	var ids []int
	for _,id := range idList{
		i,_ := strconv.Atoi(id)
		ids = append(ids,i)
	}

	err := UseCase.DeleteStocks(ids)
	if err != nil {
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

func getStocks (c *gin.Context){

	barcode := c.Query("barcode")
	name := c.Query("name")
	description := c.Query("description")
	category := c.Query("category")
	dealerId,_ := strconv.Atoi(c.Query("dealerId"))

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")
	//isDropdown,_ := strconv.ParseBool(c.Query("isDropdown"))

	p, err := UseCase.GetStocks(barcode,name,description,category,orderBy,orderAs,pageNumber, pageSize,dealerId)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

// ###########################################

func createPerson (c *gin.Context){
	p := Person{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.CreatePerson(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func updatePerson (c *gin.Context){
	p := Person{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.UpdatePerson(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getPersonById (c *gin.Context){
	id,_ := strconv.Atoi(c.Query("id"))

	p, err := UseCase.GetPersonById(id)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func deletePeople(c *gin.Context){

	idList := strings.Split(c.Query("ids"),",")

	var ids []int
	for _,id := range idList{
		i,_ := strconv.Atoi(id)
		ids = append(ids,i)
	}

	err := UseCase.DeletePersons(ids)
	if err != nil {
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

func getPeople (c *gin.Context){

	name := c.Query("name")
	pType := c.Query("pType")

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")
	isDropdown,_ := strconv.ParseBool(c.Query("isDropdown"))

	p, err := UseCase.GetPeople(name,pType,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	if isDropdown {
		ResponseList := []*responses.PersonDropdownResponse{}
		for _,per := range p.Items{
			r := &responses.PersonDropdownResponse{}
			r.Id = per.Id
			r.Name = per.Name
			r.Type = per.Type
			ResponseList = append(ResponseList, r)
		}
		c.JSON(200, generateSuccessResponse(ResponseList))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

// #############################################################

func createReceiving (c *gin.Context){
	p := Receiving{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.CreateReceiving(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func updateReceiving (c *gin.Context){
	p := Receiving{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.UpdateReceiving(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getReceivingById (c *gin.Context){
	id,_ := strconv.Atoi(c.Query("id"))

	p, err := UseCase.GetReceivingById(id)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func deleteReceivings(c *gin.Context){

	idList := strings.Split(c.Query("ids"),",")

	var ids []int
	for _,id := range idList{
		i,_ := strconv.Atoi(id)
		ids = append(ids,i)
	}

	err := UseCase.DeleteReceivings(ids)
	if err != nil {
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

func getReceivings (c *gin.Context){

	person := c.Query("person")
	status := c.Query("status")

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")

	p, err := UseCase.GetReceivings(person,status,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func setReceivingStatus (c *gin.Context){

	id,_ := strconv.Atoi(c.Query("id"))
	status := c.Query("status")

	err := UseCase.SetReceivingStatus(status,id)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

// #######################################################

func createPayment (c *gin.Context){
	p := Payment{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.CreatePayment(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func updatePayment (c *gin.Context){
	p := Payment{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.UpdatePayment(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getPaymentById (c *gin.Context){
	id,_ := strconv.Atoi(c.Query("id"))

	p, err := UseCase.GetPaymentById(id)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func deletePayments(c *gin.Context){

	idList := strings.Split(c.Query("ids"),",")

	var ids []int
	for _,id := range idList{
		i,_ := strconv.Atoi(id)
		ids = append(ids,i)
	}

	err := UseCase.DeletePayments(ids)
	if err != nil {
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

func getPayments (c *gin.Context){

	person := c.Query("person")
	status := c.Query("status")

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")

	p, err := UseCase.GetPayments(person,status,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func setPaymentStatus (c *gin.Context){

	id,_ := strconv.Atoi(c.Query("id"))
	status := c.Query("status")

	err := UseCase.SetPaymentStatus(status,id)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

// ###############################################################


func createExpense (c *gin.Context){
	p := Expense{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.CreateExpense(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func updateExpense (c *gin.Context){
	p := Expense{}
	c.BindJSON(&p)

	p.UserId = getUserIdFromToken(c)
	err := UseCase.UpdateExpense(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getExpenseById (c *gin.Context){
	id,_ := strconv.Atoi(c.Query("id"))

	p, err := UseCase.GetExpenseById(id)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func deleteExpenses(c *gin.Context){

	idList := strings.Split(c.Query("ids"),",")

	var ids []int
	for _,id := range idList{
		i,_ := strconv.Atoi(id)
		ids = append(ids,i)
	}

	err := UseCase.DeleteExpenses(ids)
	if err != nil {
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

func getExpenses (c *gin.Context){

	name := c.Query("name")
	description := c.Query("description")

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")
	//isDropdown,_ := strconv.ParseBool(c.Query("isDropdown"))

	p, err := UseCase.GetExpenses(name,description,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

// ###############################################################

func createUser (c *gin.Context){
	p := User{}
	c.BindJSON(&p)

	err := UseCase.CreateUser(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func updateUser (c *gin.Context){
	p := User{}
	c.BindJSON(&p)

	err := UseCase.UpdateUser(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getUserById (c *gin.Context){
	id,_ := strconv.Atoi(c.Query("id"))

	p, err := UseCase.GetUserById(id)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func deleteUsers(c *gin.Context){

	idList := strings.Split(c.Query("ids"),",")

	var ids []int
	for _,id := range idList{
		i,_ := strconv.Atoi(id)
		ids = append(ids,i)
	}

	err := UseCase.DeleteUsers(ids)
	if err != nil {
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

func getUsers (c *gin.Context){

	name := c.Query("name")
	email := c.Query("email")

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")
	isDropdown,_ := strconv.ParseBool(c.Query("isDropdown"))

	p, err := UseCase.GetUsers(name,email,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	if isDropdown {
		ResponseList := []*responses.UserDropdownResponse{}
		for _,per := range p.Items{
			r := &responses.UserDropdownResponse{}
			r.Id = per.Id
			r.Name = per.Name
			r.Email = per.Email
			r.Phone = per.Phone
			ResponseList = append(ResponseList, r)
		}
		c.JSON(200, generateSuccessResponse(ResponseList))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

// #########################################################

func createSale (c *gin.Context){
	p := Sale{}
	c.BindJSON(&p)

	err := UseCase.CreateSale(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func updateSale (c *gin.Context){
	p := Sale{}
	c.BindJSON(&p)

	err := UseCase.UpdateSale(&p)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func getSaleById (c *gin.Context){
	id,_ := strconv.Atoi(c.Query("id"))

	p, err := UseCase.GetSaleById(id)

	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

func deleteSales(c *gin.Context){

	idList := strings.Split(c.Query("ids"),",")

	var ids []int
	for _,id := range idList{
		i,_ := strconv.Atoi(id)
		ids = append(ids,i)
	}

	err := UseCase.DeleteSales(ids)
	if err != nil {
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse("ok"))
}

func getSales (c *gin.Context){

	tInterval := c.Query("timeInterval")

	pageNumber,_ := strconv.Atoi(c.Query("pageNumber"))
	pageSize,_ := strconv.Atoi(c.Query("pageSize"))

	orderBy := c.Query("orderBy")
	orderAs := c.Query("orderAs")
	//isDropdown,_ := strconv.ParseBool(c.Query("isDropdown"))

	p, err := UseCase.GetSales(tInterval,orderBy,orderAs,pageNumber, pageSize)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(p))
}

// ###############################################################


func handleLogin (c *gin.Context){
	p := LoginParams{}
	c.BindJSON(&p)

	u,t,err := UseCase.Login(p.Email,p.Password,secret)
	if err != nil{
		c.JSON(200, generateFailResponse(err))
		return
	}

	c.JSON(200, generateSuccessResponse(LoginResponse{Token: t, User: u}))

}

func generateSuccessResponse(data interface{}) (map[string]interface{}) {

	return gin.H{"data": data, "success": true, "errorCode": 0, "errorMessage": ""}
}

func generateFailResponse( err *ErrorType) (map[string]interface{}){
	return gin.H{"data": nil , "success": false, "errorCode": err.Code, "errorMessage": err.Message}
}