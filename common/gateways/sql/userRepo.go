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

const stTableUser = `CREATE TABLE IF NOT EXISTS %s.user (
  id      INT AUTO_INCREMENT PRIMARY KEY,
  name    VARCHAR(70) NOT NULL DEFAULT '',
  phone   VARCHAR(20) DEFAULT '',
  address VARCHAR(250) DEFAULT '',
  email   VARCHAR(100) DEFAULT '' UNIQUE,
  password  VARCHAR(30) DEFAULT '',
  token   VARBINARY(100) DEFAULT '',
  register_date    INT
)ENGINE=InnoDB DEFAULT CHARSET=utf8;`



const stSelectUserById = `SELECT * FROM %s.user
									 WHERE id=?`

const stSelectUserByEmail = `SELECT * FROM %s.user
									 WHERE email=?`


const stInsertUser = `INSERT INTO %s.user (name,phone,address,email,password,token,register_date)
							VALUES (?,?,?,?,?,?,?)`

const stUpdateUserById = `UPDATE %s.user SET name=?,phone=?,address=?,email=?,password=?,token=?
								WHERE id=?`

const stDeleteUserById = `DELETE FROM %s.user WHERE id=?`

type UserRepo struct {}

var usr *UserRepo
var qSelectUserById,qSelectUserByEmail,qInsertUser,qUpdateUserById,qDeleteUserById *sql.Stmt

func GetUserRepo() *UserRepo{
	if usr == nil {
		usr = &UserRepo{}

		var err error
		if _, err = DB.Exec(s(stTableUser)); err != nil {
			LogError(err)
		}

		qSelectUserById, err = DB.Prepare(s(stSelectUserById))
		if err != nil {
			LogError(err)
		}

		qSelectUserByEmail, err = DB.Prepare(s(stSelectUserByEmail))
		if err != nil {
			LogError(err)
		}

		qInsertUser, err = DB.Prepare(s(stInsertUser))
		if err != nil {
			LogError(err)
		}

		qUpdateUserById, err = DB.Prepare(s(stUpdateUserById))
		if err != nil {
			LogError(err)
		}
		qDeleteUserById, err = DB.Prepare(s(stDeleteUserById))
		if err != nil {
			LogError(err)
		}
	}

	return usr
}

func (usr *UserRepo) SelectUserById(id int)(*User,error){
	p := &User{}
	var temp string
	row := qSelectUserById.QueryRow(id)
	err := row.Scan(&p.Id,&p.Name,&p.Phone,&p.Address,&p.Email,&temp,&p.Token,&p.RegisterDate)
	if err != nil{
		LogError(err)
		return nil, err
	}
	return p,nil
}

func (usr *UserRepo) SelectUserByEmail(email string)(*User,error){
	p := &User{}
	var temp string
	row := qSelectUserByEmail.QueryRow(email)
	err := row.Scan(&p.Id,&p.Name,&p.Phone,&p.Address,&p.Email,&temp,&p.Token,&p.RegisterDate)
	if err != nil{
		LogError(err)
		return nil, err
	}
	return p,nil
}


func (usr *UserRepo) InsertUser(p *User)(error){

	result,err := qInsertUser.Exec(p.Name,p.Phone,p.Address,p.Email,p.Password,p.Token,p.RegisterDate)
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

func (usr *UserRepo) UpdateUserById(p *User, IdToUpdate int)(error){

	_,err := qUpdateUserById.Exec(p.Name,p.Phone,p.Address,p.Email,p.Password,p.Token,IdToUpdate)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (usr *UserRepo) DeleteUserById(Id int)(error){

	_,err := qDeleteUserById.Exec(Id)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (usr *UserRepo) DeleteUsers(ids []int)(error){


	stDelete := `DELETE FROM %s.user WHERE id in (`

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

func (usr *UserRepo) SelectUsers(name,email,orderBy,orderAs string,pageNumber, pageSize int) (*responses.UserResponse,  error) {

	usres := &responses.UserResponse{}
	items := []*responses.UserItem{}

	var nameAvail bool
	var emailAvail bool

	var orderByAvail bool
	var pageNumberAvail bool
	var pageSizeAvail bool

	if len(name) > 0{
		nameAvail = true
	}
	if len(email) > 0{
		emailAvail = true
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

	stSelect := `SELECT * FROM %s.user`
	stCount := `SELECT COUNT(*) FROM %s.user`

	stSelect = s(stSelect)
	stCount = s(stCount)

	filter := ``

	if  nameAvail || emailAvail {
		filter += " WHERE "

		if nameAvail {
			filter +=  ` name LIKE ` + `'%` + name + `%' `

			if emailAvail {
				filter += " AND "
			}
		}

		if emailAvail {
			filter +=  ` email LIKE ` + `'%` + email + `%' `
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

	var temp string
	for rows.Next(){
		p := &responses.UserItem{}
		err = rows.Scan(&p.Id,&p.Name,&p.Phone,&p.Address,&p.Email,&temp,&p.Token,&p.RegisterDate)
		if err != nil {
			LogError(err)
		}
		items = append(items, p)
	}

	usres.Items = items

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
	count.Scan(&usres.Count)

	return usres,nil
}

func (usr *UserRepo) Close() {
	qSelectUserById.Close()
	qInsertUser.Close()
	qUpdateUserById.Close()
	qDeleteUserById.Close()
}