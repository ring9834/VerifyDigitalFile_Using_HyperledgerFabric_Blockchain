package utils

import (
	"fmt"
	"hzx/goSqlHelper"
	"time"

	"github.com/astaxie/beego"
)

const (
	CONN_LIVE_TIME = 24 //连接使用时间 小时
)

func OpenDB() *goSqlHelper.SqlHelper {
	host := beego.AppConfig.String("mssql_host")
	port, err := beego.AppConfig.Int("mssql_port")
	if err != nil {
		port = 1433
	}

	user := beego.AppConfig.String("user")
	password := beego.AppConfig.String("password")
	dbName := beego.AppConfig.String("db_name")

	connString := fmt.Sprintf("server=%s;port%d;database=%s;user id=%s;password=%s", host, port, dbName, user, password)
	db, err2 := goSqlHelper.MssqlOpen("mssql", connString)
	db.Connection.SetMaxOpenConns(200)
	db.Connection.SetMaxIdleConns(50)
	db.Connection.SetConnMaxLifetime(time.Duration(CONN_LIVE_TIME) * time.Hour)

	if err2 != nil {
		panic(err2) //stop exec
	}
	return db
}
