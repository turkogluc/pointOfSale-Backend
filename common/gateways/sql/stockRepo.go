package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	. "stock/common/logger"
	. "stock/entities"
	"strconv"
	"stock/entities/responses"
)

const stTableStock = `CREATE TABLE IF NOT EXISTS %s.stock (
  id         INT AUTO_INCREMENT PRIMARY KEY,
  product_id INT      NOT NULL,
  qty     INT      NOT NULL,
  dealer_id  INT      NOT NULL,
  creation_date       INT NOT NULL,
  update_date INT NOT NULL,
  user_id 		 INT 	DEFAULT 1,
  favorite		BOOLEAN	DEFAULT false,
  FOREIGN KEY (user_id) REFERENCES %s.user (id) ON DELETE CASCADE ON UPDATE CASCADE,
  FOREIGN KEY (product_id) REFERENCES %s.product (id) ON DELETE CASCADE ON UPDATE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;`



const stSelectStockById = `SELECT * FROM %s.stock
									 WHERE id=?`

const stSelectStockByProductId = `SELECT * FROM %s.stock
									 WHERE product_id=?`

const stInsertStock = `INSERT INTO %s.stock (product_id,qty,dealer_id,creation_date,update_date,user_id,favorite)
							VALUES (?,?,?,?,?,?,?)`

const stUpdateStockById = `UPDATE %s.stock SET product_id=?, qty=?, dealer_id=?, update_date=?,user_id=?,favorite=?
								WHERE id=?`

const stDecrementProductFromStock = `UPDATE %s.stock SET qty = qty - ?
								WHERE product_id=?`

const stSetFavoriteByProductId = `UPDATE %s.stock SET favorite = ?
								WHERE product_id=?`

const stDeleteStockById = `DELETE FROM %s.stock WHERE id=?`

type StockRepo struct {}

var st *StockRepo
var qSelectStockById,qInsertStock,qSetFavoriteByProductId,qSelectStockByProductId,qUpdateStockById,qDecrementProductFromStock,qDeleteStockById *sql.Stmt

func GetStockRepo() *StockRepo{
	if st == nil {
		st = &StockRepo{}

		var err error
		if _, err = DB.Exec(sss(stTableStock)); err != nil {
			LogError(err)
		}

		qSelectStockById, err = DB.Prepare(s(stSelectStockById))
		if err != nil {
			LogError(err)
		}

		qSelectStockByProductId, err = DB.Prepare(s(stSelectStockByProductId))
		if err != nil {
			LogError(err)
		}

		qInsertStock, err = DB.Prepare(s(stInsertStock))
		if err != nil {
			LogError(err)
		}

		qUpdateStockById, err = DB.Prepare(s(stUpdateStockById))
		if err != nil {
			LogError(err)
		}

		qDecrementProductFromStock, err = DB.Prepare(s(stDecrementProductFromStock))
		if err != nil {
			LogError(err)
		}

		qSetFavoriteByProductId, err = DB.Prepare(s(stSetFavoriteByProductId))
		if err != nil {
			LogError(err)
		}

		qDeleteStockById, err = DB.Prepare(s(stDeleteStockById))
		if err != nil {
			LogError(err)
		}
	}

	return st
}

func (st *StockRepo) SelectStockById(id int)(*Stock,error){
	p := &Stock{}
	row := qSelectStockById.QueryRow(id)
	err := row.Scan(&p.Id,&p.ProductId,&p.Qty,&p.DealerId,&p.CreationDate,&p.UpdateDate,&p.UserId,&p.IsFavorite)
	if err != nil{
		LogError(err)
		return nil, err
	}
	return p,nil
}

func (st *StockRepo) SelectStockByProductId(productId int)(*Stock,error){
	p := &Stock{}
	row := qSelectStockByProductId.QueryRow(productId)
	err := row.Scan(&p.Id,&p.ProductId,&p.Qty,&p.DealerId,&p.CreationDate,&p.UpdateDate,&p.UserId,&p.IsFavorite)
	if err != nil{
		LogError(err)
		return nil, err
	}
	return p,nil
}


func (st *StockRepo) InsertStock(p *Stock)(error){

	result,err := qInsertStock.Exec(p.ProductId,p.Qty,p.DealerId,p.CreationDate,p.UpdateDate,p.UserId,p.IsFavorite)
	if err != nil{
		LogError(err)
		return err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		LogError(err)
		return err
	}
	p.Id = int(lastId)

	return nil
}

func (st *StockRepo) UpdateStockById(p *Stock, IdToUpdate int)(error){

	_,err := qUpdateStockById.Exec(p.ProductId,p.Qty,p.DealerId,p.UpdateDate,p.UserId,p.IsFavorite,IdToUpdate)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (st *StockRepo) DecrementProductFromStock(productId,count int)(error){

	_,err := qDecrementProductFromStock.Exec(count,productId)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (st *StockRepo) SetFavoriteByProductId(productId int,fav bool)(error){

	_,err := qSetFavoriteByProductId.Exec(fav,productId)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (st *StockRepo) DeleteStockById(Id int)(error){

	_,err := qDeleteStockById.Exec(Id)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (st *StockRepo) DeleteStocks(ids []int)(error){


	stDelete := `DELETE FROM %s.stock WHERE id in (`

	for k,v := range ids{
		stDelete += strconv.FormatInt(int64(v),10)
		if k < len(ids)-1{
			stDelete+=`,`
		}
	}
	stDelete += `)`

	stDelete = s(stDelete)
	LogDebug(stDelete)

	qDelete,err := DB.Prepare(stDelete)
	defer qDelete.Close()
	if err != nil{
		LogDebug(err)
		return err
	}

	_,err = qDelete.Exec()
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}


func (st *StockRepo) SelectStocks(timeInterval []int,barcode,name,description,category,orderBy,orderAs string,pageNumber, pageSize,dealerId int,creatorId int,isFavorite bool) (*responses.StockResponse,  error) {

	stes := &responses.StockResponse{}
	items := []*responses.StockItem{}

	var timeAvail bool

	var barAvail bool
	var nameAvail bool
	var descAvail bool
	var catAvail bool
	var dealerAvail bool
	var crtAvail bool
	var favAvail bool

	var orderByAvail bool
	var pageNumberAvail bool
	var pageSizeAvail bool


	if len(timeInterval) > 1{
		timeAvail = true
	}

	if len(barcode) > 0{
		barAvail = true
	}
	if len(name) > 0{
		nameAvail = true
	}
	if len(description) > 0{
		descAvail = true
	}
	if len(category) > 0{
		catAvail = true
	}
	if dealerId > 0 {
		dealerAvail = true
	}
	if creatorId > 0 {
		crtAvail = true
	}
	if isFavorite == true {
		favAvail = true
	}

	if len(orderBy) != 0{
		if orderBy != "name" {
			orderByAvail = true
		}
	}
	if pageNumber > 0 {
		pageNumberAvail = true
	}
	if pageSize > 0 {
		pageSizeAvail = true
	}

	stSelect := `SELECT s.id,s.product_id,s.qty,s.dealer_id,s.creation_date,s.update_date,p.barcode,p.name,p.description,p.category,p.purchase_price,p.sale_price,p.register_date,IFNULL(per.id,0),IFNULL(per.name,""),s.user_id,u.name,s.favorite,p.image_path
						FROM %s.stock AS s
						JOIN %s.product AS p ON s.product_id = p.id 
						LEFT JOIN %s.person AS per ON s.dealer_id = per.id
						JOIN %s.user AS u ON s.user_id=u.id`

	stCount := `SELECT COUNT(*) FROM %s.stock AS s
						JOIN %s.product AS p ON s.product_id = p.id 
						LEFT JOIN %s.person AS per ON s.dealer_id = per.id
						JOIN %s.user AS u ON s.user_id = u.id `

	stSelect = s4(stSelect)
	stCount = s4(stCount)

	filter := ``

	if  crtAvail || timeAvail || barAvail || nameAvail || descAvail || catAvail || dealerAvail || favAvail{
		filter += " WHERE "

		if favAvail {
			filter +=  ` s.favorite = true `

			if crtAvail || timeAvail || barAvail || nameAvail || descAvail || catAvail || dealerAvail{
				filter += " AND "
			}
		}
		if crtAvail {
			filter +=  ` s.user_id = ` + strconv.FormatInt(int64(creatorId),10)

			if timeAvail || barAvail || nameAvail || descAvail || catAvail || dealerAvail{
				filter += " AND "
			}
		}

		if timeAvail{
			filter += " s.update_date > " + strconv.FormatInt(int64(timeInterval[0]),10)
			filter += " AND s.update_date < " + strconv.FormatInt(int64(timeInterval[1]),10)

			if barAvail || nameAvail || descAvail || catAvail || dealerAvail{
				filter += " AND "
			}
		}

		if barAvail {
			filter +=  ` p.barcode LIKE ` + `'%` + barcode + `%' `

			if nameAvail || descAvail || catAvail || dealerAvail{
				filter += " AND "
			}
		}

		if nameAvail {
			filter +=  ` p.name LIKE ` + `'%` + name + `%' `

			if descAvail || catAvail || dealerAvail {
				filter += " AND "
			}

		}

		if descAvail {
			filter +=  ` p.description LIKE ` + `'%` + description + `%' `

			if catAvail || dealerAvail{
				filter += " AND "
			}

		}

		if catAvail {
			filter +=  ` p.category LIKE ` + `'%` + category + `%' `

			if dealerAvail{
				filter += " AND "
			}
		}

		if dealerAvail {
			filter +=  ` s.dealer_id = ` + strconv.FormatInt(int64(dealerId),10)
		}
	}

	stSelect += filter
	stCount += filter

	stSelect += ` ORDER BY `
	if orderByAvail {
		stSelect +=  orderBy
	}else{
		stSelect += ` p.name `
	}
	if orderAs == "desc"{
		stSelect += ` DESC `
	}else{
		stSelect += ` ASC `
	}
	if pageNumberAvail && pageSizeAvail {
		offset := strconv.FormatInt(int64((pageNumber-1)*pageSize),10)
		pageSizeStr := strconv.FormatInt(int64(pageSize),10)
		stSelect += ` LIMIT ` + offset + `,` + pageSizeStr
	}

	LogDebug(stSelect)

	qSelect, err := DB.Prepare(stSelect)
	defer qSelect.Close()

	if err != nil{
		LogError(err)
		return nil, err
	}

	rows, err := qSelect.Query()
	if err != nil{
		LogError(err)
		return nil, err
	}

	for rows.Next(){
		p := &responses.StockItem{}
		p.Product = &Product{}
		err = rows.Scan(&p.Id,&p.Product.Id,&p.Qty,&p.DealerId,&p.CreationDate,&p.UpdateDate,&p.Product.Barcode,&p.Product.Name,&p.Product.Description,&p.Product.Category,&p.Product.PurchasePrice,&p.Product.SalePrice,&p.Product.RegisterDate,&p.DealerId,&p.DealerName,&p.UserId,&p.UserName,&p.IsFavorite,&p.Product.ImagePath)
		if err != nil {
			LogError(err)
		}
		items = append(items, p)
	}

	stes.Items = items

	qCount, err := DB.Prepare(stCount)
	defer qCount.Close()
	if err != nil{
		LogError(err)
		return nil, err
	}
	count := qCount.QueryRow()
	if err != nil{
		LogError(err)
		return nil, err
	}
	count.Scan(&stes.Count)

	return stes,nil
}

func (st *StockRepo) SelectCurrentStockReport(name,category,orderBy,orderAs string,pageNumber, pageSize int) (*responses.CurrentStockReportResponse,  error) {

	stes := &responses.CurrentStockReportResponse{}
	items := []*responses.CurrentStockReportItem{}

	var nameAvail bool
	var catAvail bool

	var orderByAvail bool
	var pageNumberAvail bool
	var pageSizeAvail bool


	if len(name) > 0{
		nameAvail = true
	}

	if len(category) > 0 {
		catAvail = true
	}

	if len(orderBy) != 0{
		if orderBy != "8" {  // order by (s.qty*p.sale_price)
			orderByAvail = true
		}
	}
	if pageNumber > 0 {
		pageNumberAvail = true
	}
	if pageSize > 0 {
		pageSizeAvail = true
	}

	stSelect := `select p.id,p.name,p.category,s.qty,p.purchase_price,p.sale_price,(s.qty*p.purchase_price),(s.qty*p.sale_price),(s.qty*p.sale_price)-(s.qty*p.purchase_price)
    FROM %s.stock as s
    JOIN %s.product p ON s.product_id = p.id `

	stCount := `select COUNT(*)
    FROM %s.stock as s
    JOIN %s.product p ON s.product_id = p.id `

	stSelect = ss(stSelect)
	stCount = ss(stCount)

	filter := ``

	if  nameAvail || catAvail {
		filter += " WHERE "

		if nameAvail {
			filter +=  ` ( p.name LIKE ` + `'%` + name + `%' OR p.barcode LIKE ` + `'%` + name + `%' OR p.category LIKE '%` + name + `%' )`

			if catAvail {
				filter += ` AND `
			}
		}
		if catAvail{
			filter += ` p.category = '` + category + `' `
		}
	}

	stSelect += filter
	stCount += filter

	stSelect += ` ORDER BY `
	if orderByAvail {
		stSelect +=  orderBy
	}else{
		stSelect += ` 8 `
	}
	if orderAs == "asc"{
		stSelect += ` ASC `
	}else{
		stSelect += ` DESC `
	}
	if pageNumberAvail && pageSizeAvail {
		offset := strconv.FormatInt(int64((pageNumber-1)*pageSize),10)
		pageSizeStr := strconv.FormatInt(int64(pageSize),10)
		stSelect += ` LIMIT ` + offset + `,` + pageSizeStr
	}

	LogDebug(stSelect)

	qSelect, err := DB.Prepare(stSelect)
	defer qSelect.Close()

	if err != nil{
		LogError(err)
		return nil, err
	}

	rows, err := qSelect.Query()
	if err != nil{
		LogError(err)
		return nil, err
	}

	for rows.Next(){
		p := &responses.CurrentStockReportItem{}
		err = rows.Scan(&p.Id,&p.Name,&p.Category,&p.Qty,&p.PurchasePrice,&p.SalePrice,&p.GrossValue,&p.NetValue,&p.TotalProfit)
		if err != nil {
			LogError(err)
		}
		items = append(items, p)
	}

	stes.Items = items

	qCount, err := DB.Prepare(stCount)
	defer qCount.Close()
	if err != nil{
		LogError(err)
		return nil, err
	}
	count := qCount.QueryRow()
	if err != nil{
		LogError(err)
		return nil, err
	}
	count.Scan(&stes.Count)

	return stes,nil
}


func (st *StockRepo) Close() {
	qSelectStockById.Close()
	qInsertStock.Close()
	qUpdateStockById.Close()
	qDeleteStockById.Close()
}
