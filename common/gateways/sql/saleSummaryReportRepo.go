package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	. "stock/common/logger"
	. "stock/entities"
	"strconv"
)

const stTableSaleSummaryReport = `CREATE TABLE IF NOT EXISTS %s.sale_summary_report (
	id             	INT AUTO_INCREMENT PRIMARY KEY,
	gross_profit   	FLOAT DEFAULT 0,
	net_profit		FLOAT DEFAULT 0,
	sale_count		INT DEFAULT 0,
	item_count		INT DEFAULT 0,
	customer_count	INT DEFAULT 0,
	discount		FLOAT DEFAULT 0,
	basket_value	FLOAT DEFAULT 0,
	basket_size		FLOAT DEFAULT 0,
	timestamp		INT DEFAULT 0
	)ENGINE=InnoDB DEFAULT CHARSET=utf8;`



const stSelectSaleSummaryReportById = `SELECT id,gross_profit,net_profit,sale_count,item_count,customer_count,discount,basket_value,basket_size,timestamp 
										FROM %s.sale_summary_report
									 	WHERE id=?`

const stSelectSaleSummaryReportByDate = `SELECT id,gross_profit,net_profit,sale_count,item_count,customer_count,discount,basket_value,basket_size,timestamp 
										FROM %s.sale_summary_report
									 	WHERE timestamp=?`

const stInsertSaleSummaryReport = `INSERT INTO %s.sale_summary_report (gross_profit,net_profit,sale_count,item_count,customer_count,discount,basket_value,basket_size,timestamp)
							VALUES (?,?,?,?,?,?,?,?,?)`

const stUpdateSaleSummaryReportById = `UPDATE %s.sale_summary_report SET
								gross_profit=?, net_profit=?, sale_count=?, item_count=?, customer_count=?, discount=?, basket_value=?, basket_size=?, timestamp=?
								WHERE id=?`

const stDeleteSaleSummaryReportById = `DELETE FROM %s.sale_summary_report WHERE id=?`

type SaleSummaryReportRepo struct {}

var slrp *SaleSummaryReportRepo
var qSelectSaleSummaryReportById, qInsertSaleSummaryReport, qUpdateSaleSummaryReportById, qDeleteSaleSummaryReportById *sql.Stmt

func GetSaleSummaryReportRepo() *SaleSummaryReportRepo {
	if slrp == nil {
		slrp = &SaleSummaryReportRepo{}

		var err error
		if _, err = DB.Exec(s(stTableSaleSummaryReport)); err != nil {
			LogError(err)
		}

		qSelectSaleSummaryReportById, err = DB.Prepare(s(stSelectSaleSummaryReportById))
		if err != nil {
			LogError(err)
		}

		qInsertSaleSummaryReport, err = DB.Prepare(s(stInsertSaleSummaryReport))
		if err != nil {
			LogError(err)
		}

		qUpdateSaleSummaryReportById, err = DB.Prepare(s(stUpdateSaleSummaryReportById))
		if err != nil {
			LogError(err)
		}
		qDeleteSaleSummaryReportById, err = DB.Prepare(s(stDeleteSaleSummaryReportById))
		if err != nil {
			LogError(err)
		}
	}

	return slrp
}

func (slrp *SaleSummaryReportRepo) SelectSaleSummaryReportById(id int)(*SaleSummaryReport,error){
	p := &SaleSummaryReport{}
	row := qSelectSaleSummaryReportById.QueryRow(id)
	err := row.Scan(&p.Id,&p.GrossProfit,&p.NetProfit,&p.SaleCount,&p.ItemCount,&p.CustomerCount,&p.Discount,&p.BasketValue,&p.BasketSize,&p.Timestamp)
	if err != nil{
		LogError(err)
		return nil, err
	}
	return p,nil
}

func (slrp *SaleSummaryReportRepo) SelectSaleSummaryReportByDate(id int)(*SaleSummaryObjectItem,error){
	p := &SaleSummaryObjectItem{}
	row := qSelectSaleSummaryReportById.QueryRow(id)
	err := row.Scan(&p.Id,&p.GrossProfit,&p.NetProfit,&p.SaleCount,&p.ItemCount,&p.CustomerCount,&p.Discount,&p.BasketValue,&p.BasketSize,&p.Timestamp)
	if err != nil{
		LogError(err)
		return nil, err
	}
	return p,nil
}


func (slrp *SaleSummaryReportRepo) InsertSaleSummaryReport(p *SaleSummaryObjectItem)(error){

	result,err := qInsertSaleSummaryReport.Exec(p.GrossProfit,p.NetProfit,p.SaleCount,p.ItemCount,p.CustomerCount,p.Discount,p.BasketValue,p.BasketSize,p.Timestamp)
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

func (slrp *SaleSummaryReportRepo) UpdateSaleSummaryReportById(p *SaleSummaryObjectItem, IdToUpdate int)(error){

	_,err := qUpdateSaleSummaryReportById.Exec(p.GrossProfit,p.NetProfit,p.SaleCount,p.ItemCount,p.CustomerCount,p.Discount,p.BasketValue,p.BasketSize,p.Timestamp,IdToUpdate)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (slrp *SaleSummaryReportRepo) DeleteSaleSummaryReportById(Id int)(error){

	_,err := qDeleteSaleSummaryReportById.Exec(Id)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (slrp *SaleSummaryReportRepo) DeleteSaleSummaryReport(ids []int)(error){


	stDelete := `DELETE FROM %s.sale_summary_report WHERE id in (`

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


func (slrp *SaleSummaryReportRepo) SelectSaleSummaryReportItems(timeInterval []int) (*SaleSummaryReport,  error) {

	var timeAvail bool
	objectItems := []*SaleSummaryObjectItem{}

	if len(timeInterval) > 0{
		timeAvail = true
	}


	stSelectItems := `SELECT gross_profit,net_profit,sale_count,item_count,customer_count,discount,basket_value,basket_size,timestamp
						FROM %s.sale_summary_report`

	stSelectItems = s(stSelectItems)

	filter := ``

	if  timeAvail {
		filter += " WHERE "

		filter += " timestamp > " + strconv.FormatInt(int64(timeInterval[0]),10)
		filter += " AND timestamp < " + strconv.FormatInt(int64(timeInterval[1]),10)
	}

	filter += " ORDER BY timestamp"

	stSelectItems += filter

	LogDebug(stSelectItems)


	qSelectItems, err := DB.Prepare(stSelectItems)
	defer qSelectItems.Close()

	if err != nil{
		LogError(err)
		return nil, err
	}

	rows, err := qSelectItems.Query()
	if err != nil{
		LogError(err)
		return nil, err
	}

	p := &SaleSummaryReport{}

	for rows.Next(){
		var gp,np,dis,bv,bs float64
		var sc,ic,cc,ts int


		err = rows.Scan(&gp,&np,&sc,&ic,&cc,&dis,&bv,&bs,&ts)
		if err != nil {
			LogError(err)
		}

		obj := &SaleSummaryObjectItem{
			GrossProfit:gp,
			NetProfit:np,
			SaleCount:sc,
			ItemCount:ic,
			CustomerCount:cc,
			Discount:dis,
			BasketValue:bv,
			BasketSize:bs,
			Timestamp:ts,
		}

		objectItems = append(objectItems,obj)

	}

	p.AsObject = &SaleSummaryObject{}
	p.AsObject.Items = objectItems
	p.AsObject.Count = len(objectItems)

	return p,nil
}

func (slrp *SaleSummaryReportRepo) Close() {
	qSelectSaleSummaryReportById.Close()
	qInsertSaleSummaryReport.Close()
	qUpdateSaleSummaryReportById.Close()
	qDeleteSaleSummaryReportById.Close()
}