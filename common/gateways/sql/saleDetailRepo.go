package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	. "stock/common/logger"
	. "stock/entities"
	"strconv"
	"time"
)

const stTableSaleDetail = `CREATE TABLE IF NOT EXISTS %s.sale_detail (
						  id             INT AUTO_INCREMENT PRIMARY KEY,
						  creation_date  INT     NOT NULL DEFAULT 0,
						  basket_id		 INT	DEFAULT 0,
						  product_id	 INT	DEFAULT 0,
						  qty			 INT	DEFAULT 0,
						  discount		 INT	DEFAULT 0,
						  user_id 		 INT 	DEFAULT 1,
						  FOREIGN KEY (user_id) REFERENCES %s.user (id) ON DELETE CASCADE ON UPDATE CASCADE,
						  FOREIGN KEY (basket_id) REFERENCES %s.sale_basket (id) ON DELETE CASCADE ON UPDATE CASCADE	
						)ENGINE=InnoDB DEFAULT CHARSET=utf8;`



const stSelectSaleDetailById = `SELECT id,creation_date,basket_id,product_id,qty,discount,user_id FROM %s.sale_detail
									 WHERE id=?`

const stInsertSaleDetail = `INSERT INTO %s.sale_detail (creation_date,basket_id,product_id,qty,discount,user_id)
							VALUES (?,?,?,?,?,?)`

const stUpdateSaleDetailById = `UPDATE %s.sale_detail SET creation_date=?,basket_id=?,product_id,qty=?,discount=?,user_id=?
								WHERE id=?`

const stDeleteSaleDetailById = `DELETE FROM %s.sale_detail WHERE id=?`

type SaleDetailRepo struct {}

var sldt *SaleDetailRepo
var qSelectSaleDetailById,qInsertSaleDetail,qUpdateSaleDetailById,qDeleteSaleDetailById *sql.Stmt

func GetSaleDetailRepo() *SaleDetailRepo{
	if sldt == nil {
		sldt = &SaleDetailRepo{}

		var err error
		if _, err = DB.Exec(sss(stTableSaleDetail)); err != nil {
			LogError(err)
		}

		qSelectSaleDetailById, err = DB.Prepare(s(stSelectSaleDetailById))
		if err != nil {
			LogError(err)
		}

		qInsertSaleDetail, err = DB.Prepare(s(stInsertSaleDetail))
		if err != nil {
			LogError(err)
		}

		qUpdateSaleDetailById, err = DB.Prepare(s(stUpdateSaleDetailById))
		if err != nil {
			LogError(err)
		}
		qDeleteSaleDetailById, err = DB.Prepare(s(stDeleteSaleDetailById))
		if err != nil {
			LogError(err)
		}
	}

	return sldt
}

func (sldt *SaleDetailRepo) SelectSaleDetailById(id int)(*SaleDetail,error){
	p := &SaleDetail{}
	row := qSelectSaleDetailById.QueryRow(id)
	err := row.Scan(&p.Id,&p.CreationDate,&p.BasketId,&p.ProductId,&p.Qty,&p.Discount,&p.UserId)
	if err != nil{
		LogError(err)
		return nil, err
	}
	return p,nil
}


func (sldt *SaleDetailRepo) InsertSaleDetail(p *SaleDetail)(error){

	timeNOW := int(time.Now().Unix())
	result,err := qInsertSaleDetail.Exec(timeNOW,p.BasketId,p.ProductId,p.Qty,p.Discount,p.UserId)
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

func (sldt *SaleDetailRepo) UpdateSaleDetailById(p *SaleDetail, IdToUpdate int)(error){

	timeNOW := int(time.Now().Unix())
	_,err := qUpdateSaleDetailById.Exec(timeNOW,p.BasketId,p.ProductId,p.Qty,p.Discount,p.UserId,IdToUpdate)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (sldt *SaleDetailRepo) DeleteSaleDetailById(Id int)(error){

	_,err := qDeleteSaleDetailById.Exec(Id)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (sldt *SaleDetailRepo) DeleteSaleDetails(ids []int)(error){


	stDelete := `DELETE FROM %s.sale_detail WHERE id in (`

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

func (sldt *SaleDetailRepo) SelectSaleDetails(timeInterval []int,productId, userId int,orderBy,orderAs string,pageNumber, pageSize int) {

	//response := &responses.SaleDetailResponse{}
	//items := []*SaleDetail{}
	//
	//var timeAvail bool
	//var userAvail bool
	//
	//var orderByAvail bool
	//var pageNumberAvail bool
	//var pageSizeAvail bool
	//
	//
	//if len(timeInterval) > 0{
	//	timeAvail = true
	//}
	//
	//if userId > 0 {
	//	userAvail = true
	//}
	//
	//if len(orderBy) != 0{
	//	if orderBy != "id" {
	//		orderByAvail = true
	//	}
	//}
	//if pageNumber > 0 {
	//	pageNumberAvail = true
	//}
	//if pageSize > 0 {
	//	pageSizeAvail = true
	//}
	//
	//stSelect := `SELECT s.id,s.creation_date,s.items,s.user_id,u.name
	//					FROM %s.sale_detail AS s
	//					JOIN %s.user AS u ON u.id=s.user_id`
	//stCount := `SELECT COUNT(*) FROM %s.sale_detail AS s
	//					JOIN %s.user AS u ON u.id=s.user_id`
	//
	//stSelect = ss(stSelect)
	//stCount = ss(stCount)
	//
	//filter := ``
	//
	//if  timeAvail || userAvail{
	//	filter += " WHERE "
	//
	//
	//	if timeAvail{
	//		filter += " s.creation_date > " + strconv.FormatInt(int64(timeInterval[0]),10)
	//		filter += " AND s.creation_date < " + strconv.FormatInt(int64(timeInterval[1]),10)
	//
	//		if userAvail{
	//			filter += ` AND `
	//		}
	//	}
	//
	//	if userAvail{
	//		filter += ` s.user_id = ` + strconv.FormatInt(int64(userId),10)
	//	}
	//
	//
	//}
	//
	//stSelect += filter
	//stCount += filter
	//
	//stSelect += ` ORDER BY `
	//if orderByAvail {
	//	stSelect +=  orderBy
	//}else{
	//	stSelect += ` s.id `
	//}
	//if orderAs == "asc"{
	//	stSelect += ` ASC `
	//}else{
	//	stSelect += ` DESC `
	//}
	//if pageNumberAvail && pageSizeAvail {
	//	offset := strconv.FormatInt(int64((pageNumber-1)*pageSize),10)
	//	pageSizeStr := strconv.FormatInt(int64(pageSize),10)
	//	stSelect += ` LIMIT ` + offset + `,` + pageSizeStr
	//}
	//
	//LogDebug(stSelect)
	//
	//qSelect, err := DB.Prepare(stSelect)
	//defer qSelect.Close()
	//
	//if err != nil{
	//	LogError(err)
	//	return nil, err
	//}
	//
	//rows, err := qSelect.Query()
	//if err != nil{
	//	LogError(err)
	//	return nil, err
	//}
	//
	//for rows.Next(){
	//	p := &SaleDetail{}
	//	err = rows.Scan(&p.Id,&p.CreationDate,&p.ItemsStr,&p.UserId,&p.UserName)
	//	if err != nil {
	//		LogError(err)
	//	}
	//	items = append(items, p)
	//}
	//
	//response.Items = items
	//
	//qCount, err := DB.Prepare(stCount)
	//defer qCount.Close()
	//if err != nil{
	//	LogError(err)
	//	return nil, err
	//}
	//count := qCount.QueryRow()
	//if err != nil{
	//	LogError(err)
	//	return nil, err
	//}
	//count.Scan(&response.Count)
	//
	//return response,nil
}

func (sldt *SaleDetailRepo) Close() {
	qSelectSaleDetailById.Close()
	qInsertSaleDetail.Close()
	qUpdateSaleDetailById.Close()
	qDeleteSaleDetailById.Close()
}