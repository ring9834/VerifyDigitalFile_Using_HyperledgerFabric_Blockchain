package goSqlHelper

import (
	"database/sql"

	"github.com/bobby96333/commonLib/stackError"
)

const (
	SQL_SELECT = "SELECT"
	SQL_UPDATE = "UPDATE"
	SQL_DELETE = "DELETE"
	SQL_INSERT = "INSERT"
)

func NewAutoSql(helper *SqlHelper) *AutoSql {
	var orm = new(AutoSql)
	orm.joins = make([]string, 0)
	orm.sqlHelper = helper
	return orm
}

type AutoSql struct {
	act       string
	sqlHelper *SqlHelper
	fieldSql  string
	tbname    string
	where     string
	groupBy   string
	orderby   string
	having    string
	limit     int
	joins     []string
	set       string
	setVals   []interface{}
}

func (ths *AutoSql) Select(fieldSql string) *AutoSql {
	ths.act = SQL_SELECT
	ths.fieldSql = fieldSql
	return ths
}
func (ths *AutoSql) Delete(tbname string) *AutoSql {
	ths.act = SQL_DELETE
	ths.tbname = tbname
	return ths
}
func (ths *AutoSql) Set(setSql string) *AutoSql {
	ths.set = setSql
	return ths
}
func (ths *AutoSql) SetRow(row *HelperRow) *AutoSql {
	sql := ""
	fields := "("
	values := "("
	vals := make([]interface{}, len(*row))
	i := -1
	for key, val := range *row {
		i++
		// if i>0{
		// 	sql+=","
		// }
		// sql+=key+"=?"
		if i > 0 {
			fields += ","
			values += ","
		}
		fields += key
		values += "?"
		vals[i] = val
	}
	fields += ") "
	values += ") "
	sql += fields + " VALUES " + values
	ths.set = sql
	ths.setVals = vals
	return ths
}
func (ths *AutoSql) Update(tbname string) *AutoSql {
	ths.act = SQL_UPDATE
	ths.tbname = tbname
	return ths
}
func (ths *AutoSql) Insert(tbname string) *AutoSql {
	ths.act = SQL_INSERT
	ths.tbname = tbname
	return ths
}
func (ths *AutoSql) From(tbname string) *AutoSql {
	ths.tbname = tbname
	return ths
}
func (ths *AutoSql) Where(where string) *AutoSql {
	ths.where = where
	return ths
}
func (ths *AutoSql) Join(joinSql string) *AutoSql {
	ths.joins = append(ths.joins, joinSql)
	return ths
}
func (ths *AutoSql) Groupby(groupBySql string) *AutoSql {
	ths.groupBy = groupBySql
	return ths
}
func (ths *AutoSql) Orderby(OrderbySql string) *AutoSql {
	ths.orderby = OrderbySql
	return ths
}
func (ths *AutoSql) Having(having string) *AutoSql {
	ths.having = having
	return ths
}
func (ths *AutoSql) Limit(limit int) *AutoSql {
	ths.limit = limit
	return ths
}
func (ths *AutoSql) GenerateSql() string {
	switch ths.act {
	case SQL_SELECT:
		return ths.generateSelectSql()
	case SQL_INSERT:
		return ths.generateInsertSql()
	case SQL_UPDATE:
		return ths.generateUpdateSql()
	case SQL_DELETE:
		return ths.generateDeleteSql()
	default:
		return ths.generateSelectSql()
	}
	//panic("no found act:" + this.act)
}

func (ths *AutoSql) QueryRows(args ...interface{}) ([]HelperRow, *stackError.StackError) {
	sql := ths.GenerateSql()
	return ths.sqlHelper.QueryRows(sql, args...)
}

func (ths *AutoSql) QueryTable(args ...interface{}) (*HelperTable, *stackError.StackError) {
	sql := ths.GenerateSql()
	return ths.sqlHelper.QueryTable(sql, args...)
}

func (ths *AutoSql) Querying(args ...interface{}) (*Querying, *stackError.StackError) {
	sql := ths.GenerateSql()
	return ths.sqlHelper.Querying(sql, args...)
}

func (ths *AutoSql) QueryRow(args ...interface{}) (HelperRow, *stackError.StackError) {
	sql := ths.GenerateSql()
	return ths.sqlHelper.QueryRow(sql, args...)
}

func (ths *AutoSql) QueryScalar(val interface{}, args ...interface{}) (interface{}, *stackError.StackError) {
	sql := ths.GenerateSql()
	err := ths.sqlHelper.QueryScalar(val, sql, args...)
	return val, err
}

// /**
//   read a int value
// */
// func (ths *AutoSql) QueryScalarInt(sql string, args ...interface{}) (int, *stackError.StackError) {
// 	return ths.QueryScalarInt(sql, args...)
// }

// /**
//   read a int value
// */
// func (ths *AutoSql) QueryScalarString(sql string, args ...interface{}) (string, *stackError.StackError) {
// 	return ths.QueryScalarString(sql, args...)
// }

func (ths *AutoSql) QueryOrm(orm IEntity, args ...interface{}) *stackError.StackError {
	sql := ths.GenerateSql()
	return ths.sqlHelper.QueryOrm(orm, sql, args...)
}

func (ths *AutoSql) Exec(args ...interface{}) (sql.Result, *stackError.StackError) {
	sql := ths.GenerateSql()
	if ths.setVals != nil {
		args = append(ths.setVals, args...)
	}
	return ths.sqlHelper.Exec(sql, args...)
}

/*
execute insert sql
*/
func (ths *AutoSql) ExecInsert(args ...interface{}) (int64, *stackError.StackError) {
	sql := ths.GenerateSql()
	if ths.setVals != nil {

		args = append(ths.setVals, args...)
	}
	return ths.sqlHelper.ExecInsert(sql, args...)
}

/*
execute update or delete sql
*/
func (ths *AutoSql) ExecUpdateOrDel(args ...interface{}) (int64, *stackError.StackError) {
	sql := ths.GenerateSql()
	if ths.setVals != nil {
		args = append(ths.setVals, args...)
	}
	return ths.sqlHelper.ExecUpdateOrDel(sql, args...)
}
