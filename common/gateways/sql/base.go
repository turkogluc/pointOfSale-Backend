package sql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	sql "github.com/rmulley/go-fast-sql"
	. "stock/common/logger"
)

var DB *sql.DB
var database string

func Init(host, port, dbName,user, pass string) {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?readTimeout=10m&parseTime=true", user, pass, host, port, "")

	DB, err = sql.Open("mysql", dsn, 100)
	if err != nil {
		LogError(err)
	}

	// Open doesn't open a connection. Validate DSN data:
	err = DB.Ping()
	if err != nil {
		LogError(err)
	}

	DB.SetMaxIdleConns(0)
	DB.SetMaxOpenConns(5000)

	database = dbName

}

func Close() {
	//if vr != nil {
	//	vr.close()
	//}
	//if vdr != nil {
	//	vdr.close()
	//}
	//if fer != nil {
	//	fer.close()
	//}
	//if vetr != nil {
	//	vetr.close()
	//}
	//if vfr != nil {
	//	vfr.close()
	//}
	//if vser != nil {
	//	vser.close()
	//}
	//if vsr != nil {
	//	vsr.close()
	//}
	//if vwtr != nil {
	//	vwtr.close()
	//}

}


func s(stmt string) string {
	return fmt.Sprintf(stmt, database)
}

func ss(stmt string) string {
	return fmt.Sprintf(stmt, database, database)
}

func sss(stmt string) string {
	return fmt.Sprintf(stmt, database, database, database)
}

func s4(stmt string) string {
	return fmt.Sprintf(stmt, database, database, database,database)
}


//func ss(stmt string) string {
//	return fmt.Sprintf(stmt, databaseFleet, databaseFleet)
//}
//
//func sss(stmt string) string {
//	return fmt.Sprintf(stmt, databaseFleet, databaseFleet, databaseFleet)
//}
//
//func s4(stmt string) string {
//	return fmt.Sprintf(stmt, databaseFleet, databaseFleet, databaseFleet, databaseFleet)
//}
//
//func s5(stmt string) string {
//	return fmt.Sprintf(stmt, databaseFleet, databaseFleet, databaseFleet, databaseFleet, databaseFleet)
//}
//
//func s6(stmt string) string {
//	return fmt.Sprintf(stmt, databaseFleet,databaseFleet,databaseFleet, databaseFleet, databaseFleet, databaseFleet)
//}
//
//func s7(stmt string) string {
//	return fmt.Sprintf(stmt, databaseFleet,databaseFleet,databaseFleet, databaseFleet, databaseFleet, databaseFleet, databaseFleet)
//}
//
//func s8(stmt string) string {
//	return fmt.Sprintf(stmt, databaseFleet,databaseFleet,databaseFleet, databaseFleet, databaseFleet, databaseFleet, databaseFleet, databaseFleet)
//}
