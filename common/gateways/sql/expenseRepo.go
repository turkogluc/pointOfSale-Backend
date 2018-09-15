package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	. "stock/common/logger"
	. "stock/entities"
	"strconv"
	"stock/entities/responses"
)

const stTableExpense = `CREATE TABLE IF NOT EXISTS %s.expense (
						  id             INT AUTO_INCREMENT PRIMARY KEY,
						  name           VARCHAR(50) NULL DEFAULT '',
						  description    VARCHAR(200) NULL DEFAULT '',
						  price 		 FLOAT    NOT NULL DEFAULT 0,
						  creation_date  INT     NOT NULL DEFAULT 0,
						  update_date    INT     NOT NULL DEFAULT 0,
						  user_id 		 INT 	DEFAULT 1,
						  FOREIGN KEY (user_id) REFERENCES %s.user (id) ON DELETE CASCADE ON UPDATE CASCADE
						)ENGINE=InnoDB DEFAULT CHARSET=utf8;`



const stSelectExpenseById = `SELECT id,name,description,price,creation_date,update_date,user_id FROM %s.expense
									 WHERE id=?`

const stInsertExpense = `INSERT INTO %s.expense (name, description , price , creation_date , update_date,user_id)
							VALUES (?,?,?,?,?,?)`

const stUpdateExpenseById = `UPDATE %s.expense SET name=? ,description=? ,price=?,update_date=?,user_id=?
								WHERE id=?`

const stDeleteExpenseById = `DELETE FROM %s.expense WHERE id=?`

type ExpenseRepo struct {}

var exp *ExpenseRepo
var qSelectExpenseById,qInsertExpense,qUpdateExpenseById,qDeleteExpenseById *sql.Stmt

func GetExpenseRepo() *ExpenseRepo{
	if exp == nil {
		exp = &ExpenseRepo{}

		var err error
		if _, err = DB.Exec(ss(stTableExpense)); err != nil {
			LogError(err)
		}

		qSelectExpenseById, err = DB.Prepare(s(stSelectExpenseById))
		if err != nil {
			LogError(err)
		}

		qInsertExpense, err = DB.Prepare(s(stInsertExpense))
		if err != nil {
			LogError(err)
		}

		qUpdateExpenseById, err = DB.Prepare(s(stUpdateExpenseById))
		if err != nil {
			LogError(err)
		}
		qDeleteExpenseById, err = DB.Prepare(s(stDeleteExpenseById))
		if err != nil {
			LogError(err)
		}
	}

	return exp
}

func (exp *ExpenseRepo) SelectExpenseById(id int)(*Expense,error){
	p := &Expense{}
	row := qSelectExpenseById.QueryRow(id)
	err := row.Scan(&p.Id,&p.Name,&p.Description,&p.Price,&p.CreateDate,&p.UpdateDate,&p.UserId)
	if err != nil{
		LogError(err)
		return nil, err
	}
	return p,nil
}


func (exp *ExpenseRepo) InsertExpense(p *Expense)(error){

	result,err := qInsertExpense.Exec(p.Name,p.Description,p.Price,p.CreateDate,p.UpdateDate,p.UserId)
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

func (exp *ExpenseRepo) UpdateExpenseById(p *Expense, IdToUpdate int)(error){

	_,err := qUpdateExpenseById.Exec(p.Name,p.Description,p.Price,p.UpdateDate,p.UserId,IdToUpdate)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (exp *ExpenseRepo) DeleteExpenseById(Id int)(error){

	_,err := qDeleteExpenseById.Exec(Id)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (exp *ExpenseRepo) DeleteExpenses(ids []int)(error){


	stDelete := `DELETE FROM %s.expense WHERE id in (`

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

func (exp *ExpenseRepo) SelectExpenses(timeInterval []int,name,description,orderBy,orderAs string,pageNumber, pageSize,creatorId int) (*responses.ExpenseResponse,  error) {

	expes := &responses.ExpenseResponse{}
	items := []*responses.ExpenseItem{}

	var timeAvail bool
	var crtAvail bool

	var nameAvail bool
	var descAvail bool

	var orderByAvail bool
	var pageNumberAvail bool
	var pageSizeAvail bool


	if len(timeInterval) > 1{
		timeAvail = true
	}
	if creatorId > 0 {
		crtAvail = true
	}

	if len(name) > 0{
		nameAvail = true
	}
	if len(description) > 0{
		descAvail = true
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

	stSelect := `SELECT e.id, e.name, e.description, e.price, e.creation_date, e.update_date, e.user_id, u.name 
				FROM %s.expense as e JOIN %s.user AS u
				ON e.user_id = u.id`
	stCount := `SELECT COUNT(*) FROM %s.expense as e JOIN %s.user AS u ON e.user_id = u.id`

	stSelect = ss(stSelect)
	stCount = ss(stCount)

	filter := ``

	if  timeAvail || crtAvail || nameAvail || descAvail {
		filter += " WHERE "

		if crtAvail {
			filter +=  ` e.user_id = ` + strconv.FormatInt(int64(creatorId),10)

			if timeAvail || nameAvail || descAvail {
				filter += " AND "
			}
		}

		if timeAvail{
			filter += " e.update_date > " + strconv.FormatInt(int64(timeInterval[0]),10)
			filter += " AND e.update_date < " + strconv.FormatInt(int64(timeInterval[1]),10)

			if nameAvail || descAvail{
				filter += " AND "
			}
		}

		if nameAvail {
			filter +=  ` u.name LIKE ` + `'%` + name + `%' `

			if descAvail {
				filter += " AND "
			}
		}

		if descAvail {
			filter +=  ` description LIKE ` + `'%` + description + `%' `
		}
	}

	stSelect += filter
	stCount += filter

	stSelect += ` ORDER BY `
	if orderByAvail {
		stSelect +=  orderBy
	}else{
		stSelect += ` u.name `
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
		p := &responses.ExpenseItem{}
		err = rows.Scan(&p.Id,&p.Name,&p.Description,&p.Price,&p.CreateDate,&p.UpdateDate,&p.UserId,&p.UserName)
		if err != nil {
			LogError(err)
		}
		items = append(items, p)
	}

	expes.Items = items

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
	count.Scan(&expes.Count)

	return expes,nil
}

func (exp *ExpenseRepo) Close() {
	qSelectExpenseById.Close()
	qInsertExpense.Close()
	qUpdateExpenseById.Close()
	qDeleteExpenseById.Close()
}