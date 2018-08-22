package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	. "stock/common/logger"
	. "stock/entities"
	"strconv"
	//"stock/entities/responses"
	"stock/entities/responses"
)

const stTablePerson = `CREATE TABLE IF NOT EXISTS %s.person (
  id      INT AUTO_INCREMENT PRIMARY KEY,
  name    VARCHAR(50) NOT NULL DEFAULT '',
  phone   VARCHAR(20) DEFAULT '',
  email   VARCHAR(50) DEFAULT '',
  address VARCHAR(200) DEFAULT '',
  p_type    ENUM('Tedarikçi','Müşteri') NOT NULL,
  creation_date    INT   NOT NULL,
  user_id 		 INT 	DEFAULT 1,
  FOREIGN KEY (user_id) REFERENCES %s.user (id) ON DELETE CASCADE ON UPDATE CASCADE
)ENGINE=InnoDB DEFAULT CHARSET=utf8;`



const stSelectPersonById = `SELECT * FROM %s.person
									 WHERE id=?`

const stInsertPerson = `INSERT INTO %s.person (name,phone,email,address,p_type,creation_date,user_id)
							VALUES (?,?,?,?,?,?,?)`

const stUpdatePersonById = `UPDATE %s.person SET name=?, phone=?, email=?, address=?, p_type=?, creation_date=?,user_id=?
								WHERE id=?`

const stDeletePersonById = `DELETE FROM %s.person WHERE id=?`

type PersonRepo struct {}

var prsn *PersonRepo
var qSelectPersonById,qInsertPerson,qUpdatePersonById,qDeletePersonById *sql.Stmt

func GetPersonRepo() *PersonRepo{
	if prsn == nil {
		prsn = &PersonRepo{}

		var err error
		if _, err = DB.Exec(ss(stTablePerson)); err != nil {
			LogError(err)
		}

		qSelectPersonById, err = DB.Prepare(s(stSelectPersonById))
		if err != nil {
			LogError(err)
		}

		qInsertPerson, err = DB.Prepare(s(stInsertPerson))
		if err != nil {
			LogError(err)
		}

		qUpdatePersonById, err = DB.Prepare(s(stUpdatePersonById))
		if err != nil {
			LogError(err)
		}
		qDeletePersonById, err = DB.Prepare(s(stDeletePersonById))
		if err != nil {
			LogError(err)
		}
	}

	return prsn
}

func (prsn *PersonRepo) SelectPersonById(id int)(*Person,error){
	p := &Person{}
	row := qSelectPersonById.QueryRow(id)
	err := row.Scan(&p.Id,&p.Name,&p.Phone,&p.Email,&p.Address,&p.Type,&p.CreationDate,&p.UserId)
	if err != nil{
		LogError(err)
		return nil, err
	}
	return p,nil
}


func (prsn *PersonRepo) InsertPerson(p *Person)(error){

	result,err := qInsertPerson.Exec(p.Name,p.Phone,p.Email,p.Address,p.Type,p.CreationDate,p.UserId)
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

func (prsn *PersonRepo) UpdatePersonById(p *Person, IdToUpdate int)(error){

	_,err := qUpdatePersonById.Exec(p.Name,p.Phone,p.Email,p.Address,p.Type,p.CreationDate,p.UserId,IdToUpdate)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (prsn *PersonRepo) DeletePersonById(Id int)(error){

	_,err := qDeletePersonById.Exec(Id)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (prsn *PersonRepo) DeletePersons(ids []int)(error){


	stDelete := `DELETE FROM %s.person WHERE id in (`

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

func (prsn *PersonRepo) SelectPeople(name,pType,orderBy,orderAs string,pageNumber, pageSize int) (*responses.PersonResponse,  error) {

	prsnes := &responses.PersonResponse{}
	items := []*responses.PersonItem{}

	var nameAvail bool
	var typeAvail bool

	var orderByAvail bool
	var pageNumberAvail bool
	var pageSizeAvail bool

	if len(name) > 0{
		nameAvail = true
	}
	if len(pType) > 0{
		typeAvail = true
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

	stSelect := `SELECT p.id,p.name,p.phone,p.email,p.address,p.p_type,p.creation_date,p.user_id,u.name
					FROM %s.person AS p
					JOIN %s.user AS u ON p.user_id = u.id`
	stCount := `SELECT COUNT(*) FROM %s.person AS p
					JOIN %s.user AS u ON p.user_id = u.id`

	stSelect = ss(stSelect)
	stCount = ss(stCount)

	filter := ``

	if  nameAvail || typeAvail {
		filter += " WHERE "

		if nameAvail {
			filter +=  ` p.name LIKE ` + `'%` + name + `%' `

			if typeAvail {
				filter += " AND "
			}
		}

		if typeAvail {
			filter +=  ` p_type LIKE ` + `'%` + pType + `%' `
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
		p := &responses.PersonItem{}
		err = rows.Scan(&p.Id,&p.Name,&p.Phone,&p.Email,&p.Address,&p.Type,&p.CreationDate,&p.UserId,&p.UserName)
		if err != nil {
			LogError(err)
		}
		items = append(items, p)
	}

	prsnes.Items = items

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
	count.Scan(&prsnes.Count)

	return prsnes,nil
}

func (prsn *PersonRepo) Close() {
	qSelectPersonById.Close()
	qInsertPerson.Close()
	qUpdatePersonById.Close()
	qDeletePersonById.Close()
}