package interactors

import (
	"golang.org/x/crypto/bcrypt"
	. "stock/common/logger"
	"github.com/robfig/cron"
	"fmt"
	"stock/common/projectArch/interactors"
	"time"
	"stock/entities"
)


func comparePasswords(hashedPwd []byte, plainPwd []byte) bool {

	err := bcrypt.CompareHashAndPassword(hashedPwd, plainPwd)
	if err != nil {
		return false
	}

	return true
}

func hashAndSalt(pwd []byte) []byte {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		LogError(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return hash
}

func StartReceivingCheckCronJob(){
	c := cron.New()
	c.AddFunc("@every 1h", func() {
		fmt.Println("StartReceivingCheckCronJob ...")
		receivings,err := interactors.ReceivingRepo.SelectReceivings([]int{},"","","","",0,0,0)
		if err != nil {
			LogError(err)
			return
		}
		timeNow := int(time.Now().Unix())
		for _,item := range receivings.Items{
			if item.ExpectedDate < timeNow{
				interactors.ReceivingRepo.SetStatus("Gecikmiş",item.Id)
			}
		}

	})
	c.Start()
}

func SaleReporterJob(interval string){
	c := cron.New()
	c.AddFunc(interval, func() {
		fmt.Println("SaleReporterJob ...")

		// when a sale is executed a new record with "processed = false" parameter is added
		// to saleBasket repo. Select the ones not processed yet
		// Processing means writing these sale records to report table

		records,err := interactors.SaleBasketRepo.RetrieveNotProcessedRecords()
		if err != nil {
			LogError(err)
			return
		}

		for _,record := range records.Items{

			// to have 1 type of timestamp format (1 timestamp for 1 day)
			// use the timestamp of YYYY-MM-DD 13:00:00 as the identifier of that specific day

			t := time.Unix(int64(record.Timestamp), 0)		// convert timestamp to time
			roundedTime := time.Date(t.Year(), t.Month(), t.Day(), 13, 0, 0, 0, t.Location()) // set the hour 13:00
			roundedTimestamp := int(roundedTime.Unix())

			// customer_id is returned to custome_count, always 1 customer is expected for 1 sale
			// therefore it might be either 0 or a non zero integer which will be a customer id and means 1 customer
			if record.CustomerCount > 0{
				record.CustomerCount = 1
			}

			existentRecord,err := interactors.SaleSummaryReportRepo.SelectSaleSummaryReportByDate(roundedTimestamp)
			if err != nil && existentRecord == nil {
				// means there is no record about this day already written
				// so we should create the first record for this day
				LogDebug(err)


				err = interactors.SaleSummaryReportRepo.InsertSaleSummaryReport(record)
				if err != nil{
					LogError(err)
					fmt.Println("Fatal error:", err)
				}else{
					// set processed true

					err := interactors.SaleBasketRepo.SetSaleBasketIsProcessedStatus(record.Id,true)
					if err != nil{
						LogError(err)
						fmt.Println("Fatal error:", err)
					}
				}


			}else if existentRecord != nil{
				// there is already a record and we should update it by summing up

				temp := &entities.SaleSummaryObjectItem{}

				temp.GrossProfit = record.GrossProfit + existentRecord.GrossProfit
				temp.NetProfit 	 = record.NetProfit + existentRecord.NetProfit
				temp.SaleCount   = record.SaleCount + existentRecord.SaleCount
				temp.ItemCount   = record.ItemCount + existentRecord.ItemCount
				temp.Discount    = record.Discount + existentRecord.Discount
				temp.CustomerCount = record.CustomerCount + existentRecord.CustomerCount
				temp.BasketValue = temp.GrossProfit / float64(temp.SaleCount)
				temp.BasketSize = float64(temp.ItemCount) / float64(temp.SaleCount)
				temp.Timestamp = roundedTimestamp

				err = interactors.SaleSummaryReportRepo.UpdateSaleSummaryReportById(temp,existentRecord.Id)
				if err != nil {
					LogError(err)
					fmt.Println("Fatal error:", err)
				}else{
					// set processed true

					err := interactors.SaleBasketRepo.SetSaleBasketIsProcessedStatus(record.Id,true)
					if err != nil{
						LogError(err)
						fmt.Println("Fatal error:", err)
					}
				}

			}

		}

	})
	c.Start()
}