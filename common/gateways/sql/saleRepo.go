package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	. "stock/common/logger"
	. "stock/entities"
	"strconv"
	"stock/entities/responses"
	"time"
)

const stTableSale = `CREATE TABLE IF NOT EXISTS %s.sale (
						  id             INT AUTO_INCREMENT PRIMARY KEY,
						  creation_date  INT     NOT NULL DEFAULT 0,
						  items			 TEXT
						)ENGINE=InnoDB DEFAULT CHARSET=utf8;`



const stSelectSaleById = `SELECT id,creation_date,items FROM %s.sale
									 WHERE id=?`

const stInsertSale = `INSERT INTO %s.sale (creation_date,items)
							VALUES (?,?)`

const stUpdateSaleById = `UPDATE %s.sale SET creation_date=?,items=?
								WHERE id=?`

const stDeleteSaleById = `DELETE FROM %s.sale WHERE id=?`

type SaleRepo struct {}

var sl *SaleRepo
var qSelectSaleById,qInsertSale,qUpdateSaleById,qDeleteSaleById *sql.Stmt

func GetSaleRepo() *SaleRepo{
	if sl == nil {
		sl = &SaleRepo{}

		var err error
		if _, err = DB.Exec(s(stTableSale)); err != nil {
			LogError(err)
		}

		qSelectSaleById, err = DB.Prepare(s(stSelectSaleById))
		if err != nil {
			LogError(err)
		}

		qInsertSale, err = DB.Prepare(s(stInsertSale))
		if err != nil {
			LogError(err)
		}

		qUpdateSaleById, err = DB.Prepare(s(stUpdateSaleById))
		if err != nil {
			LogError(err)
		}
		qDeleteSaleById, err = DB.Prepare(s(stDeleteSaleById))
		if err != nil {
			LogError(err)
		}
	}

	return sl
}

func (sl *SaleRepo) SelectSaleById(id int)(*Sale,error){
	p := &Sale{}
	row := qSelectSaleById.QueryRow(id)
	err := row.Scan(&p.Id,&p.CreationDate,&p.ItemsStr)
	if err != nil{
		LogError(err)
		return nil, err
	}
	return p,nil
}


func (sl *SaleRepo) InsertSale(p *Sale)(error){

	timeNOW := int(time.Now().Unix())
	result,err := qInsertSale.Exec(timeNOW,p.ItemsStr)
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

func (sl *SaleRepo) UpdateSaleById(p *Sale, IdToUpdate int)(error){

	timeNOW := int(time.Now().Unix())
	_,err := qUpdateSaleById.Exec(timeNOW,p.ItemsStr,IdToUpdate)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (sl *SaleRepo) DeleteSaleById(Id int)(error){

	_,err := qDeleteSaleById.Exec(Id)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (sl *SaleRepo) DeleteSales(ids []int)(error){


	stDelete := `DELETE FROM %s.sale WHERE id in (`

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

func (sl *SaleRepo) SelectSales(timeInterval []int,orderBy,orderAs string,pageNumber, pageSize int) (*responses.SaleResponse,  error) {

	response := &responses.SaleResponse{}
	items := []*Sale{}

	var timeAvail bool

	var orderByAvail bool
	var pageNumberAvail bool
	var pageSizeAvail bool


	if len(timeInterval) > 0{
		timeAvail = true
	}

	if len(orderBy) != 0{
		if orderBy != "id" {
			orderByAvail = true
		}
	}
	if pageNumber > 0 {
		pageNumberAvail = true
	}
	if pageSize > 0 {
		pageSizeAvail = true
	}

	stSelect := `SELECT * FROM %s.sale`
	stCount := `SELECT COUNT(*) FROM %s.sale`

	stSelect = s(stSelect)
	stCount = s(stCount)

	filter := ``

	if  timeAvail {
		filter += " WHERE "

		filter += " creation_date > " + strconv.FormatInt(int64(timeInterval[0]),10)
		filter += " AND creation_date < " + strconv.FormatInt(int64(timeInterval[1]),10)
	}

	stSelect += filter
	stCount += filter

	stSelect += ` ORDER BY `
	if orderByAvail {
		stSelect +=  orderBy
	}else{
		stSelect += ` id `
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
		p := &Sale{}
		err = rows.Scan(&p.Id,&p.CreationDate,&p.ItemsStr)
		if err != nil {
			LogError(err)
		}
		items = append(items, p)
	}

	response.Items = items

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
	count.Scan(&response.Count)

	return response,nil
}

func (sl *SaleRepo) Close() {
	qSelectSaleById.Close()
	qInsertSale.Close()
	qUpdateSaleById.Close()
	qDeleteSaleById.Close()
}