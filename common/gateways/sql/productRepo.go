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
						  register_date  INT     NOT NULL DEFAULT 0,
						  user_id 		 INT 	DEFAULT 1,
						  image_path	TEXT NOT NULL,
  						  FOREIGN KEY (user_id) REFERENCES %s.user (id) ON DELETE CASCADE ON UPDATE CASCADE
						)ENGINE=InnoDB DEFAULT CHARSET=utf8;`



const stSelectProductById = `SELECT id,barcode,name,description,category,purchase_price,sale_price,register_date,user_id,image_path FROM %s.product
									 WHERE id=?`

const stSelectProductCategories = `SELECT category FROM %s.product
									GROUP BY category`

const stInsertProduct = `INSERT INTO %s.product (barcode,name,description,category,purchase_price,sale_price,register_date,user_id,image_path)
							VALUES (?,?,?,?,?,?,?,?,?)`

const stUpdateProductById = `UPDATE %s.product SET barcode=?, name=?, description=?, category=?, purchase_price=?, sale_price=?, register_date=?, user_id=?,image_path=?
								WHERE id=?`

const stDeleteProductById = `DELETE FROM %s.product WHERE id=?`

type ProductRepo struct {}

var pr *ProductRepo
var qSelectProductById,qSelectProductCategories,qInsertProduct,qUpdateProductById,qDeleteProductById *sql.Stmt

func GetProductRepo() *ProductRepo{
	if pr == nil {
		pr = &ProductRepo{}

		var err error
		if _, err = DB.Exec(ss(stTableProduct)); err != nil {
			LogError(err)
		}

		qSelectProductById, err = DB.Prepare(s(stSelectProductById))
		if err != nil {
			LogError(err)
		}

		qSelectProductCategories, err = DB.Prepare(s(stSelectProductCategories))
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
	err := row.Scan(&p.Id,&p.Barcode,&p.Name,&p.Description,&p.Category,&p.PurchasePrice,&p.SalePrice,&p.RegisterDate,&p.UserId,&p.ImagePath)
	if err != nil{
		LogError(err)
		return nil, err
	}
	return p,nil
}

func (pr *ProductRepo) SelectProductCategories()([]string,error){

	var resp []string
	rows,err := qSelectProductCategories.Query()
	if err != nil{
		LogError(err)
		return nil, err
	}

	for rows.Next(){
		var temp string
		err = rows.Scan(&temp)
		if err != nil {
			LogError(err)
			return nil, err
		}
		resp = append(resp, temp)
	}

	return resp,nil

}


func (pr *ProductRepo) InsertProduct(p *Product)(error){

	result,err := qInsertProduct.Exec(p.Barcode,p.Name,p.Description,p.Category,p.PurchasePrice,p.SalePrice,p.RegisterDate,p.UserId,p.ImagePath)
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

	_,err := qUpdateProductById.Exec(p.Barcode,p.Name,p.Description,p.Category,p.PurchasePrice,p.SalePrice,p.RegisterDate,p.UserId,p.ImagePath,IdToUpdate)
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
	items := []*responses.ProductItem{}

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

	stSelect := `SELECT p.id, p.barcode, p.name, p.description, p.category, p.purchase_price, p.sale_price, p.register_date, p.user_id, u.name,p.image_path
				FROM %s.product as p
				JOIN %s.user AS u ON p.user_id = u.id`
	stCount := `SELECT COUNT(*) FROM %s.product as p
				JOIN %s.user AS u ON p.user_id = u.id`

	stSelect = ss(stSelect)
	stCount = ss(stCount)

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
			filter +=  ` p.name LIKE ` + `'%` + name + `%' `

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
		p := &responses.ProductItem{}
		err = rows.Scan(&p.Id,&p.Barcode,&p.Name,&p.Description,&p.Category,&p.PurchasePrice,&p.SalePrice,&p.RegisterDate,&p.UserId,&p.UserName,&p.ImagePath)
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