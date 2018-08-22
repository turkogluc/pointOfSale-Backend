package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	. "stock/common/logger"
	. "stock/entities"
	"strconv"
	"stock/entities/responses"
)

const stTableReceiving = `CREATE TABLE IF NOT EXISTS %s.receiving (
		id              INT AUTO_INCREMENT PRIMARY KEY,
		person_id		INT NOT NULL,
		amount			FLOAT NOT NULL,
		creation_date	INT,
		update_date		INT,
		expected_date	INT,
		product_ids		VARCHAR(200) DEFAULT '',
		status			ENUM('Bekliyor','Bitti','Gecikmiş'),
		user_id 		 INT 	DEFAULT 1,
  		FOREIGN KEY (user_id) REFERENCES %s.user (id) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (person_id) REFERENCES %s.person (id) ON DELETE CASCADE ON UPDATE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;`



const stSelectReceivingById = `SELECT * FROM %s.receiving
									 WHERE id=?`

const stInsertReceiving = `INSERT INTO %s.receiving (person_id,amount,creation_date,update_date,expected_date,product_ids,status,user_id)
							VALUES (?,?,?,?,?,?,?,?)`

const stUpdateReceivingById = `UPDATE %s.receiving SET person_id=?, amount=?, update_date=?, expected_date=?,product_ids=?, status=?, user_id=?
								WHERE id=?`

const stSetStatus = `UPDATE %s.receiving SET status=?
								WHERE id=?`

const stDeleteReceivingById = `DELETE FROM %s.receiving WHERE id=?`

type ReceivingRepo struct {}

var rcv *ReceivingRepo
var qSelectReceivingById,qInsertReceiving,qSetStatus,qUpdateReceivingById,qDeleteReceivingById *sql.Stmt

func GetReceivingRepo() *ReceivingRepo{
	if rcv == nil {
		rcv = &ReceivingRepo{}

		var err error
		if _, err = DB.Exec(sss(stTableReceiving)); err != nil {
			LogError(err)
		}

		qSelectReceivingById, err = DB.Prepare(s(stSelectReceivingById))
		if err != nil {
			LogError(err)
		}

		qInsertReceiving, err = DB.Prepare(s(stInsertReceiving))
		if err != nil {
			LogError(err)
		}

		qUpdateReceivingById, err = DB.Prepare(s(stUpdateReceivingById))
		if err != nil {
			LogError(err)
		}

		qSetStatus, err = DB.Prepare(s(stSetStatus))
		if err != nil {
			LogError(err)
		}
		qDeleteReceivingById, err = DB.Prepare(s(stDeleteReceivingById))
		if err != nil {
			LogError(err)
		}
	}

	return rcv
}

func (rcv *ReceivingRepo) SelectReceivingById(id int)(*Receiving,error){
	p := &Receiving{}
	row := qSelectReceivingById.QueryRow(id)
	err := row.Scan(&p.Id,&p.PersonId,&p.Amount,&p.CreationDate,&p.UpdateDate,&p.ExpectedDate,&p.ProductIds,&p.Status,&p.UserId)
	if err != nil{
		LogError(err)
		return nil, err
	}
	return p,nil
}


func (rcv *ReceivingRepo) InsertReceiving(p *Receiving)(error){

	result,err := qInsertReceiving.Exec(p.PersonId,p.Amount,p.CreationDate,p.UpdateDate,p.ExpectedDate,p.ProductIds,p.Status,p.UserId)
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

func (rcv *ReceivingRepo) UpdateReceivingById(p *Receiving, IdToUpdate int)(error){

	_,err := qUpdateReceivingById.Exec(p.PersonId,p.Amount,p.UpdateDate,p.ExpectedDate,p.ProductIds,p.Status,p.UserId,IdToUpdate)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (rcv *ReceivingRepo) SetStatus(status string, IdToUpdate int)(error){

	_,err := qSetStatus.Exec(status,IdToUpdate)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (rcv *ReceivingRepo) DeleteReceivingById(Id int)(error){

	_,err := qDeleteReceivingById.Exec(Id)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (rcv *ReceivingRepo) DeleteReceivings(ids []int)(error){


	stDelete := `DELETE FROM %s.receiving WHERE id in (`

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

func (rcv *ReceivingRepo) SelectReceivings(person,status,orderBy,orderAs string,pageNumber, pageSize int) (*responses.ReceivingResponse,  error) {

	rcves := &responses.ReceivingResponse{}
	items := []*ReceivingsItem{}

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

	stSelect := `SELECT r.id,r.person_id,r.amount,r.creation_date,r.update_date,r.expected_date,r.status,p.name,p.phone,r.product_ids,r.user_id,u.name
						FROM %s.receiving as r
						JOIN %s.person as p ON r.person_id = p.id
						JOIN %s.user as u ON u.id = r.user_id`

	stCount := `SELECT COUNT(*) FROM %s.receiving as r
						JOIN %s.person as p ON r.person_id = p.id
						JOIN %s.user as u ON u.id = r.user_id`

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
		p := &ReceivingsItem{}
		err = rows.Scan(&p.Id,&p.PersonId,&p.Amount,&p.CreationDate,&p.UpdateDate,&p.ExpectedDate,&p.Status,&p.PersonName,&p.PersonPhone,&p.ProductIds,&p.UserId,&p.UserName)
		if err != nil {
			LogError(err)
		}
		items = append(items, p)
	}

	rcves.Items = items

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
	count.Scan(&rcves.Count)

	return rcves,nil
}

func (rcv *ReceivingRepo) Close() {
	qSelectReceivingById.Close()
	qInsertReceiving.Close()
	qUpdateReceivingById.Close()
	qDeleteReceivingById.Close()
}