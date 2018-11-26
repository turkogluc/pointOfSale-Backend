package main

import (
	"os"
	"log"
	"github.com/spf13/viper"
	"github.com/fsnotify/fsnotify"
	"stock/apps/dashboard/controller"
	. "stock/common/logger"
	"stock/common/gateways/sql"
	common "stock/common/projectArch/interactors"
)

type Specification struct {
	SqlHost     string
	SqlPort     string
	SqlDB       string
	SqlUser     string
	SqlPass     string
	Debug		bool
	LogFile		string
}

type Environment struct {
	Env string
}

var s Specification
var e Environment

var SaleReporterJobInterval string

var logFile *os.File

func main() {
	var err error
	viper.SetConfigName("config")

	// Init configuration
	viper.AddConfigPath("./")
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Fatal(err)
	}

	SaleReporterJobInterval = viper.GetString("app.reporterJobInterval")

	s.SqlHost = viper.GetString("app.sql.host")
	s.SqlPort = viper.GetString("app.sql.port")
	s.SqlUser = viper.GetString("app.sql.user")
	s.SqlPass = viper.GetString("app.sql.pass")
	s.SqlDB = viper.GetString("app.sql.db")

	s.Debug = viper.GetBool("app.log.debug")
	s.LogFile = viper.GetString("app.log.file")

	// setup watching for config file changes
	viper.WatchConfig()
	log.Println("APP:", viper.Get("app.name"))
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name, "APP:", viper.Get("app.name"))
	})

	// Init image path
	controllers.SetImagePath(viper.GetString("app.image.path"))

	// Init logging
	logFile, err = os.OpenFile(s.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(logFile)
	InitLogger(logFile, logFile, logFile, s.Debug)


	sql.Init(s.SqlHost,s.SqlPort,s.SqlDB, s.SqlUser,s.SqlPass)
	defer sql.Close()


	common.ProductRepo = sql.GetProductRepo()
	common.PersonRepo = sql.GetPersonRepo()
	common.StockRepo = sql.GetStockRepo()
	common.ReceivingRepo = sql.GetReceivingRepo()
	common.PaymentRepo = sql.GetPaymentRepo()
	common.ExpenseRepo = sql.GetExpenseRepo()
	common.UserRepo = sql.GetUserRepo()
	common.SaleBasketRepo = sql.GetSaleBasketRepo()
	common.SaleDetailRepo = sql.GetSaleDetailRepo()
	common.SaleSummaryReportRepo = sql.GetSaleSummaryReportRepo()

	controllers.StartApplicationBackend(SaleReporterJobInterval)
}