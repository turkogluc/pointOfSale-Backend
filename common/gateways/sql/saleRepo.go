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
						  items			 TEXT,
						  user_id 		 INT 	DEFAULT 1,
						  FOREIGN KEY (user_id) REFERENCES %s.user (id) ON DELETE CASCADE ON UPDATE CASCADE	
						)ENGINE=InnoDB DEFAULT CHARSET=utf8;`



const stSelectSaleById = `SELECT id,creation_date,items,user_id FROM %s.sale
									 WHERE id=?`

const stInsertSale = `INSERT INTO %s.sale (creation_date,items,user_id)
							VALUES (?,?,?)`

const stUpdateSaleById = `UPDATE %s.sale SET creation_date=?,items=?,user_id=?
								WHERE id=?`

const stDeleteSaleById = `DELETE FROM %s.sale WHERE id=?`

type SaleRepo struct {}

var sl *SaleRepo
var qSelectSaleById,qInsertSale,qUpdateSaleById,qDeleteSaleById *sql.Stmt

func GetSaleRepo() *SaleRepo{
	if sl == nil {
		sl = &SaleRepo{}

		var err error
		if _, err = DB.Exec(ss(stTableSale)); err != nil {
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
	err := row.Scan(&p.Id,&p.CreationDate,&p.ItemsStr,&p.UserId)
	if err != nil{
		LogError(err)
		return nil, err
	}
	return p,nil
}


func (sl *SaleRepo) InsertSale(p *Sale)(error){

	timeNOW := int(time.Now().Unix())
	result,err := qInsertSale.Exec(timeNOW,p.ItemsStr,p.UserId)
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
	_,err := qUpdateSaleById.Exec(timeNOW,p.ItemsStr,p.UserId,IdToUpdate)
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

func (sl *SaleRepo) SelectSales(timeInterval []int,userId int,orderBy,orderAs string,pageNumber, pageSize int) (*responses.SaleResponse,  error) {

	response := &responses.SaleResponse{}
	items := []*Sale{}

	var timeAvail bool
	var userAvail bool

	var orderByAvail bool
	var pageNumberAvail bool
	var pageSizeAvail bool


	if len(timeInterval) > 0{
		timeAvail = true
	}

	if userId > 0 {
		userAvail = true
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

	stSelect := `SELECT s.id,s.creation_date,s.items,s.user_id,u.name 
						FROM %s.sale AS s
						JOIN %s.user AS u ON u.id=s.user_id`
	stCount := `SELECT COUNT(*) FROM %s.sale AS s
						JOIN %s.user AS u ON u.id=s.user_id`

	stSelect = ss(stSelect)
	stCount = ss(stCount)

	filter := ``

	if  timeAvail || userAvail{
		filter += " WHERE "


		if timeAvail{
			filter += " s.creation_date > " + strconv.FormatInt(int64(timeInterval[0]),10)
			filter += " AND s.creation_date < " + strconv.FormatInt(int64(timeInterval[1]),10)

			if userAvail{
				filter += ` AND `
			}
		}

		if userAvail{
			filter += ` s.user_id = ` + strconv.FormatInt(int64(userId),10)
		}


	}

	stSelect += filter
	stCount += filter

	stSelect += ` ORDER BY `
	if orderByAvail {
		stSelect +=  orderBy
	}else{
		stSelect += ` s.id `
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
		err = rows.Scan(&p.Id,&p.CreationDate,&p.ItemsStr,&p.UserId,&p.UserName)
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