package models

import (
	_ "github.com/denisenkom/go-mssqldb"
)

type FilePathInfo struct {
	ID         int
	Name       string
	FullName   string
	IsDir      bool
	CreateTime string
	Extension  string
}

// func GetPagedPayments(pageindex int, pagesize int) ([]goSqlHelper.HelperRow, int, int) {
// 	tbName := "pay_demo"
// 	fields := "*"
// 	where := " 1=1 "
// 	sort := " Id ASC "
// 	rows, pagecount, recordcount := utils.GetPagedDataTable(tbName, fields, where, sort, pageindex, pagesize)
// 	return rows, pagecount, recordcount
// }
