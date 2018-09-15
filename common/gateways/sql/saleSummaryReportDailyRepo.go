package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	. "stock/common/logger"
	. "stock/entities"
	"strconv"
)

const stTableSaleSummaryReportDaily = `CREATE TABLE IF NOT EXISTS %s.sale_summary_report_daily (
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



const stSelectSaleSummaryReportDailyById = `SELECT id,gross_profit,net_profit,sale_count,item_count,customer_count,discount,basket_value,basket_size,timestamp 
										FROM %s.sale_summary_report_daily
									 	WHERE id=?`

const stInsertSaleSummaryReportDaily = `INSERT INTO %s.sale_summary_report_daily (gross_profit,net_profit,sale_count,item_count,customer_count,discount,basket_value,basket_size,timestamp)
							VALUES (?,?,?,?,?,?,?,?,?)`

const stUpdateSaleSummaryReportDailyById = `UPDATE %s.sale_summary_report_daily SET
								gross_profit=?, net_profit=?, sale_count=?, item_count=?, customer_count=?, discount=?, basket_value=?, basket_size=?, timestamp=?
								WHERE id=?`

const stDeleteSaleSummaryReportDailyById = `DELETE FROM %s.sale_summary_report_daily WHERE id=?`

type SaleSummaryReportDailyRepo struct {}

var slrp *SaleSummaryReportDailyRepo
var qSelectSaleSummaryReportDailyById,qInsertSaleSummaryReportDaily,qUpdateSaleSummaryReportDailyById,qDeleteSaleSummaryReportDailyById *sql.Stmt

func GetSaleSummaryReportDailyRepo() *SaleSummaryReportDailyRepo{
	if slrp == nil {
		slrp = &SaleSummaryReportDailyRepo{}

		var err error
		if _, err = DB.Exec(s(stTableSaleSummaryReportDaily)); err != nil {
			LogError(err)
		}

		qSelectSaleSummaryReportDailyById, err = DB.Prepare(s(stSelectSaleSummaryReportDailyById))
		if err != nil {
			LogError(err)
		}

		qInsertSaleSummaryReportDaily, err = DB.Prepare(s(stInsertSaleSummaryReportDaily))
		if err != nil {
			LogError(err)
		}

		qUpdateSaleSummaryReportDailyById, err = DB.Prepare(s(stUpdateSaleSummaryReportDailyById))
		if err != nil {
			LogError(err)
		}
		qDeleteSaleSummaryReportDailyById, err = DB.Prepare(s(stDeleteSaleSummaryReportDailyById))
		if err != nil {
			LogError(err)
		}
	}

	return slrp
}

func (slrp *SaleSummaryReportDailyRepo) SelectSaleSummaryReportDailyById(id int)(*SaleSummaryReportDaily,error){
	p := &SaleSummaryReportDaily{}
	row := qSelectSaleSummaryReportDailyById.QueryRow(id)
	err := row.Scan(&p.Id,&p.GrossProfit,&p.NetProfit,&p.SaleCount,&p.ItemCount,&p.CustomerCount,&p.Discount,&p.BasketValue,&p.BasketSize,&p.Timestamp)
	if err != nil{
		LogError(err)
		return nil, err
	}
	return p,nil
}


func (slrp *SaleSummaryReportDailyRepo) InsertSaleSummaryReportDaily(p *SaleSummaryReportDaily)(error){

	result,err := qInsertSaleSummaryReportDaily.Exec(p.GrossProfit,p.NetProfit,p.SaleCount,p.ItemCount,p.CustomerCount,p.Discount,p.BasketValue,p.BasketSize,p.Timestamp)
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

func (slrp *SaleSummaryReportDailyRepo) UpdateSaleSummaryReportDailyById(p *SaleSummaryReportDaily, IdToUpdate int)(error){

	_,err := qUpdateSaleSummaryReportDailyById.Exec(p.GrossProfit,p.NetProfit,p.SaleCount,p.ItemCount,p.CustomerCount,p.Discount,p.BasketValue,p.BasketSize,p.Timestamp,IdToUpdate)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (slrp *SaleSummaryReportDailyRepo) DeleteSaleSummaryReportDailyById(Id int)(error){

	_,err := qDeleteSaleSummaryReportDailyById.Exec(Id)
	if err != nil{
		LogError(err)
		return err
	}

	return nil
}

func (slrp *SaleSummaryReportDailyRepo) DeleteSaleSummaryReportDaily(ids []int)(error){


	stDelete := `DELETE FROM %s.sale_summary_report_daily WHERE id in (`

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


func (slrp *SaleSummaryReportDailyRepo) SelectSaleSummaryReportDailyItems(timeInterval []int) (*SaleSummaryReportDaily,  error) {

	var timeAvail bool
	objectItems := []*SaleSummaryObjectItem{}

	if len(timeInterval) > 0{
		timeAvail = true
	}


	stSelectItems := `SELECT gross_profit,net_profit,sale_count,item_count,customer_count,discount,basket_value,basket_size,timestamp
						FROM %s.sale_summary_report_daily`

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

	p := &SaleSummaryReportDaily{}

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

func (slrp *SaleSummaryReportDailyRepo) Close() {
	qSelectSaleSummaryReportDailyById.Close()
	qInsertSaleSummaryReportDaily.Close()
	qUpdateSaleSummaryReportDailyById.Close()
	qDeleteSaleSummaryReportDailyById.Close()
}