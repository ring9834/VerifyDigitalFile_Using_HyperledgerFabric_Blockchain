package models

import (
	"hzx/goSqlHelper"
	"hzx/utils"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	//db *sql.DB = nil //全局数据库连接
	db *goSqlHelper.SqlHelper
)

type PaymentRecord struct {
	Id           int64
	AccountID    int64
	PartnerID    string
	UserID       string
	CreateTime   string
	Amount       float64
	OuterTradeNo string
	Remark       string
	Status       int
	Msg          string
}
type PaymentRecordStr struct {
	AccountID    string
	PartnerID    string
	UserID       string
	CreateTime   string
	Amount       string
	OuterTradeNo string
	Remark       string
}

func AddPaymenRec(rec PaymentRecordStr) PaymentRecord {
	// isql := "INSERT INTO pay_demo(account_id,partner_id,user_id,amount,outer_tradeno,remark) VALUES (?,?,?,?,?,?)"
	// accountID, _ := strconv.ParseInt(rec.AccountID, 10, 64)
	// amount, _ := strconv.ParseFloat(rec.Amount, 64)
	// response := PaymentRecord{0, accountID, rec.PartnerID, rec.UserID, rec.CreateTime, amount, rec.OuterTradeNo, rec.Remark, 0, ""}
	// rst, err := db.Exec(isql, accountID, rec.PartnerID, rec.UserID, amount, rec.OuterTradeNo, rec.Remark)
	// if err == nil {
	// 	response.Id, _ = rst.LastInsertId()
	// 	response.Status = 1
	// 	response.Msg = "已生效"
	// 	return response
	// }
	// return response
	db = utils.OpenDB()
	accountID, _ := strconv.ParseInt(rec.AccountID, 10, 64)
	amount, _ := strconv.ParseFloat(rec.Amount, 64)
	record := goSqlHelper.HelperRow{
		"account_id":    accountID,
		"partner_id":    rec.PartnerID,
		"user_id":       rec.UserID,
		"amount":        amount,
		"outer_tradeno": rec.OuterTradeNo,
		"remark":        rec.Remark,
	}

	response := PaymentRecord{0, accountID, rec.PartnerID, rec.UserID, rec.CreateTime, amount, rec.OuterTradeNo, rec.Remark, 0, ""}
	rst, err := db.Auto().Insert("pay_demo").SetRow(&record).ExecInsert()
	if err == nil {
		response.Id = rst
		response.Status = 1
		response.Msg = "已生效"
		return response
	}
	return response
}

func GetPaymenRec(accountID int64) PaymentRecord {
	isql := "SELECT * FROM pay_demo WHERE account_id=?"
	var response PaymentRecord
	response.Msg = "失败"
	rows, err := db.Querying(isql, accountID)
	if err == nil {
		err = rows.Scan(&response.Id, &response.AccountID, &response.PartnerID, &response.UserID, &response.Amount, &response.OuterTradeNo, &response.Remark)
		if err == nil {
			response.Status = 2
			response.Msg = "成功"
			return response
		}
	}
	return response
}

func GetPayments() []goSqlHelper.HelperRow {
	if db == nil {
		db = utils.OpenDB()
	}
	dt, err := db.Auto().Select("*").From("pay_demo").QueryRows()
	if err == nil {
		return dt
	}
	return nil
}

func GetPagedPayments(pageindex int, pagesize int) ([]goSqlHelper.HelperRow, int, int) {
	tbName := "pay_demo"
	fields := "*"
	where := " 1=1 "
	sort := " Id ASC "
	rows, pagecount, recordcount := utils.GetPagedDataTable(tbName, fields, where, sort, pageindex, pagesize)
	return rows, pagecount, recordcount
}
