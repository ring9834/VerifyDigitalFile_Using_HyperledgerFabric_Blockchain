package goSqlHelper

import (
	"fmt"
	"strconv"
)

func (ths *AutoSql) generateSelectSql() string {
	field := "*"
	if ths.fieldSql != "" {
		field = ths.fieldSql
	}
	sql := "SELECT " + field
	if ths.tbname != "" {
		sql += " FROM " + ths.tbname
	}
	for _, join := range ths.joins {
		sql += " " + join
	}
	if ths.where != "" {
		sql += " WHERE " + ths.where
	}
	if ths.groupBy != "" {
		sql += " GROUP BY " + ths.groupBy
	}
	if ths.having != "" {
		sql += " HAVING " + ths.having
	}
	if ths.orderby != "" {
		sql += " ORDER BY " + ths.orderby
	}
	if ths.limit != 0 {
		sql += " LIMIT " + strconv.Itoa(ths.limit)
	}
	return sql
}

func (ths *AutoSql) generateUpdateSql() string {

	sql := fmt.Sprintf("UPDATE %s set %s ", ths.tbname, ths.set)
	if ths.where != "" {
		sql += " WHERE " + ths.where
	}
	if ths.orderby != "" {
		sql += " ORDER BY " + ths.orderby
	}
	if ths.limit != 0 {
		sql += " LIMIT " + strconv.Itoa(ths.limit)
	}
	return sql
}

func (ths *AutoSql) generateDeleteSql() string {

	sql := fmt.Sprintf("DELETE FROM %s ", ths.tbname)
	if ths.where != "" {
		sql += " WHERE " + ths.where
	}
	if ths.orderby != "" {
		sql += " ORDER BY " + ths.orderby
	}
	if ths.limit != 0 {
		sql += " LIMIT " + strconv.Itoa(ths.limit)
	}
	return sql
}
func (ths *AutoSql) generateInsertSql() string {

	sql := fmt.Sprintf("INSERT INTO %s ", ths.tbname)
	if ths.set != "" {
		//sql += " SET " + ths.set
		sql += ths.set
	}
	return sql
}
