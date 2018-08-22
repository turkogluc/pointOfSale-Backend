package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	. "stock/common/logger"
	. "stock/entities"
	"strconv"
	"stock/entities/responses"
)

const stTablePayment = `CREATE TABLE IF NOT EXISTS %s.payment (
		id              INT AUTO_INCREMENT PRIMARY KEY,
		person_id		INT NOT NULL,
		amount			FLOAT NOT NULL,
		creation_date	INT,
		update_date		INT,
		expected_date	INT,
		status			ENUM('Bekliyor','Bitti','Gecikmiş'),
		summary			VARCHAR(400) DEFAULT '',
		user_id 		 INT 	DEFAULT 1,
		FOREIGN KEY (user_id) REFERENCES %s.user (id) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (person_id) REFERENCES %s.person (id) ON DELETE CASCADE ON UPDATE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;`

const stSelectPaymentById = `SELECT * FROM %s.payment
									 WHERE id=?`

const stInsertPayment = `INSERT INTO %s.payment (person_id,amount,creation_date,update_date,expected_date,status,summary,user_id)
							VALUES (?,?,?,?,?,?,?,?)`

const stUpdatePaymentById = `UPDATE %s.payment SET person_id=?, amount=?, update_date=?, expected_date=?, status=?, summary=?, user_id=?
								WHERE id=?`

const stPaymentStatus = `UPDATE %s.payment SET status=?
								WHERE id=?`

const stDeletePaymentById = `DELETE FROM %s.payment WHERE id=?`

type PaymentRepo struct {}

var pym *PaymentRepo
var qSelectPaymentById,qInsertPayment,qUpdatePaymentById,qPaymentStatus,qDeletePaymentById *sql.Stmt

func GetPaymentRepo() *PaymentRepo{
	if pym == nil {
		pym = &PaymentRepo{}

		var err error
		if _, err = DB.Exec(sss(stTablePayment)); err != nil {
			LogError(err)
		}

		qSelectPaymentById, err = DB.Prepare(s(stSelectPaymentById))
		if err != nil {
			LogError(err)
		}

		qInsertPayment, err = DB.Prepare(s(stInsertPayment))
		if err != nil {
			LogError(err)
		}

		qUpdatePaymentById, err = DB.Prepare(s(stUpdatePaymentById))
		if err != nil {
			LogError(err)
		}
		qPaymentStatus, err = DB.Prepare(s(stPaymentStatus))
		if err != nil {
			LogError(err)
		}
		qDeletePaymentById, err = DB.Prepare(s(stDeletePaymentById))
		if err != nil {
			LogError(err)
		}
	}

	return pym
}

func (pym *PaymentRepo) SelectPaymentById(id int)(*Payment,error){
	p := &Payment{}
	row := qSelectPaymentById.QueryRow(id)
	err := row.Scan(&p.Id,&p.PersonId,&p.Amount,&p.CreationDate,&p.UpdateDate,&p.ExpectedDate,&p.Status,&p.Summary,&p.UserId)
	if err != nil{
		LogError(err)
		return nil, err
	}
	return p,nil
}


func (pym *PaymentRepo) InsertPayment(p *Payment)(error){

	result,err := qInsertPayment.Exec(p.PersonId,p.Amount,p.CreationDate,p.UpdateDate,p.ExpectedDate,p.Status,p.Summary,p.UserId)
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

func (pym *PaymentRepo) UpdatePaymentById(p *Payment, IdToUpdate int)(error){

	_,err := qUpdatePaymentById.Exec(p.PersonId,p.Amount,p.UpdateDate,p.ExpectedDate,p.Status,p.Summary,p.UserId,IdToUpdate)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (pym *PaymentRepo) SetPaymentStatus(status string, IdToUpdate int)(error){

	_,err := qPaymentStatus.Exec(status,IdToUpdate)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (pym *PaymentRepo) DeletePaymentById(Id int)(error){

	_,err := qDeletePaymentById.Exec(Id)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (pym *PaymentRepo) DeletePayments(ids []int)(error){


	stDelete := `DELETE FROM %s.payment WHERE id in (`

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

func (pym *PaymentRepo) SelectPayments(person,status,orderBy,orderAs string,pageNumber, pageSize int) (*responses.PaymentResponse,  error) {

	pymes := &responses.PaymentResponse{}
	items := []*PaymentsItem{}

	var personAvail bool
	var statusAvail bool

	var orderByAvail bool
	var pageNumberAvail bool
	var pageSizeAvail bool



	if len(person) > 0{
		personAvail = true
	}
	if len(status) > 0{
		statusAvail = true
	}

	if len(orderBy) != 0{
		if orderBy != "r.expected_date" {
			orderByAvail = true
		}
	}
	if pageNumber > 0 {
		pageNumberAvail = true
	}
	if pageSize > 0 {
		pageSizeAvail = true
	}

	stSelect := `SELECT r.id,r.person_id,r.amount,r.creation_date,r.update_date,r.expected_date,r.status,p.name,p.phone,r.summary,r.user_id,u.name
						FROM %s.payment as r
						JOIN %s.person as p ON r.person_id = p.id
						JOIN %s.user as u ON r.user_id = u.id`

	stCount := `SELECT COUNT(*) FROM %s.payment as r
						JOIN %s.person as p ON r.person_id = p.id
						JOIN %s.user as u ON r.user_id = u.id`

	stSelect = sss(stSelect)
	stCount = sss(stCount)

	filter := ``

	if  personAvail || statusAvail {
		filter += " WHERE "


		if personAvail {
			filter +=  ` p.name LIKE ` + `'%` + person + `%' `

			if statusAvail {
				filter += " AND "
			}
		}

		if statusAvail {
			filter +=  ` r.status LIKE ` + `'%` + status + `%' `
		}

	}

	stSelect += filter
	stCount += filter

	stSelect += ` ORDER BY `
	if orderByAvail {
		stSelect +=  orderBy
	}else{
		stSelect += ` r.expected_date `
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
		p := &PaymentsItem{}
		err = rows.Scan(&p.Id,&p.PersonId,&p.Amount,&p.CreationDate,&p.UpdateDate,&p.ExpectedDate,&p.Status,&p.PersonName,&p.PersonPhone,&p.Summary,&p.UserId,&p.UserName)
		if err != nil {
			LogError(err)
		}
		items = append(items, p)
	}

	pymes.Items = items

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
	count.Scan(&pymes.Count)

	return pymes,nil
}

func (pym *PaymentRepo) Close() {
	qSelectPaymentById.Close()
	qInsertPayment.Close()
	qUpdatePaymentById.Close()
	qDeletePaymentById.Close()
}