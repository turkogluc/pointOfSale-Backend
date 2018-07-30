package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	. "stock/common/logger"
	. "stock/entities"
	"strconv"
	"stock/entities/responses"
)

const stTableProduct = `CREATE TABLE IF NOT EXISTS %s.product (
						  id             INT AUTO_INCREMENT PRIMARY KEY,
						  barcode        VARCHAR(50) NOT NULL DEFAULT '' UNIQUE,
						  name           VARCHAR(50) NULL DEFAULT '',
						  description    VARCHAR(200) NULL DEFAULT '',
						  category       VARCHAR(100) NULL DEFAULT '',
						  purchase_price FLOAT    NOT NULL DEFAULT 0,
						  sale_price     FLOAT    NOT NULL DEFAULT 0,
						  register_date  INT     NOT NULL DEFAULT 0
						)ENGINE=InnoDB DEFAULT CHARSET=utf8;`



const stSelectProductById = `SELECT id,barcode,name,description,category,purchase_price,sale_price,register_date FROM %s.product
									 WHERE id=?`

const stInsertProduct = `INSERT INTO %s.product (barcode,name,description,category,purchase_price,sale_price,register_date)
							VALUES (?,?,?,?,?,?,?)`

const stUpdateProductById = `UPDATE %s.product SET barcode=?, name=?, description=?, category=?, purchase_price=?, sale_price=?, register_date=?
								WHERE id=?`

const stDeleteProductById = `DELETE FROM %s.product WHERE id=?`

type ProductRepo struct {}

var pr *ProductRepo
var qSelectProductById,qInsertProduct,qUpdateProductById,qDeleteProductById *sql.Stmt

func GetProductRepo() *ProductRepo{
	if pr == nil {
		pr = &ProductRepo{}

		var err error
		if _, err = DB.Exec(s(stTableProduct)); err != nil {
			LogError(err)
		}

		qSelectProductById, err = DB.Prepare(s(stSelectProductById))
		if err != nil {
			LogError(err)
		}

		qInsertProduct, err = DB.Prepare(s(stInsertProduct))
		if err != nil {
			LogError(err)
		}

		qUpdateProductById, err = DB.Prepare(s(stUpdateProductById))
		if err != nil {
			LogError(err)
		}
		qDeleteProductById, err = DB.Prepare(s(stDeleteProductById))
		if err != nil {
			LogError(err)
		}
	}

	return pr
}

func (pr *ProductRepo) SelectProductById(id int)(*Product,error){
	p := &Product{}
	row := qSelectProductById.QueryRow(id)
	err := row.Scan(&p.Id,&p.Barcode,&p.Name,&p.Description,&p.Category,&p.PurchasePrice,&p.SalePrice,&p.RegisterDate)
	if err != nil{
		LogError(err)
		return nil, err
	}
	return p,nil
}


func (pr *ProductRepo) InsertProduct(p *Product)(error){

	result,err := qInsertProduct.Exec(p.Barcode,p.Name,p.Description,p.Category,p.PurchasePrice,p.SalePrice,p.RegisterDate)
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

func (pr *ProductRepo) UpdateProductById(p *Product, IdToUpdate int)(error){

	_,err := qUpdateProductById.Exec(p.Barcode,p.Name,p.Description,p.Category,p.PurchasePrice,p.SalePrice,p.RegisterDate,IdToUpdate)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (pr *ProductRepo) DeleteProductById(Id int)(error){

	_,err := qDeleteProductById.Exec(Id)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (pr *ProductRepo) DeleteProducts(ids []int)(error){


	stDelete := `DELETE FROM %s.product WHERE id in (`

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

func (pr *ProductRepo) SelectProducts(barcode,name,description,category,orderBy,orderAs string,pageNumber, pageSize int) (*responses.ProductResponse,  error) {

	pres := &responses.ProductResponse{}
	items := []*Product{}

	var barAvail bool
	var nameAvail bool
	var descAvail bool
	var catAvail bool

	var orderByAvail bool
	var pageNumberAvail bool
	var pageSizeAvail bool



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

	stSelect := `SELECT * FROM %s.product`
	stCount := `SELECT COUNT(*) FROM %s.product`

	stSelect = s(stSelect)
	stCount = s(stCount)

	filter := ``

	if  barAvail || nameAvail || descAvail || catAvail {
		filter += " WHERE "


		if barAvail {
			filter +=  ` barcode LIKE ` + `'%` + barcode + `%' `

			if nameAvail || descAvail || catAvail {
				filter += " AND "
			}
		}

		if nameAvail {
			filter +=  ` name LIKE ` + `'%` + name + `%' `

			if descAvail || catAvail {
				filter += " AND "
			}

		}

		if descAvail {
			filter +=  ` description LIKE ` + `'%` + description + `%' `

			if catAvail {
				filter += " AND "
			}

		}

		if catAvail {
			filter +=  ` category LIKE ` + `'%` + category + `%' `
		}
	}

	stSelect += filter
	stCount += filter

	stSelect += ` ORDER BY `
	if orderByAvail {
		stSelect +=  orderBy
	}else{
		stSelect += ` name `
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
		p := &Product{}
		err = rows.Scan(&p.Id,&p.Barcode,&p.Name,&p.Description,&p.Category,&p.PurchasePrice,&p.SalePrice,&p.RegisterDate)
		if err != nil {
			LogError(err)
		}
		items = append(items, p)
	}

	pres.Items = items

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
	count.Scan(&pres.Count)

	return pres,nil
}

func (pr *ProductRepo) Close() {
	qSelectProductById.Close()
	qInsertProduct.Close()
	qUpdateProductById.Close()
	qDeleteProductById.Close()
}