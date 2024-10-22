package goSqlHelper

import (
	"database/sql"

	"github.com/bobby96333/commonLib/stackError"
)

func (ths *SqlHelper) query(sqlStr string, args ...interface{}) (*sql.Rows, *stackError.StackError) {

	var err error
	var rows *sql.Rows
	if ths.tx != nil {
		if ths.context == nil {
			rows, err = ths.tx.QueryContext(ths.context, sqlStr, args)
		} else {
			rows, err = ths.tx.Query(sqlStr, args)
		}
	} else if ths.context != nil {
		rows, err = ths.Connection.QueryContext(ths.context, sqlStr, args)
	} else {
		rows, err = ths.Connection.Query(sqlStr, args...)
	}
	return rows, stackError.NewFromError(err, ths.stckErrorPowerId)
}

func (ths *SqlHelper) prepare(sqlStr string) (*sql.Stmt, *stackError.StackError) {
	var smt *sql.Stmt
	var err error
	if ths.tx != nil {
		if ths.context == nil {
			smt, err = ths.tx.PrepareContext(ths.context, sqlStr)
		} else {
			smt, err = ths.tx.Prepare(sqlStr)
		}
	} else if ths.context != nil {
		smt, err = ths.Connection.PrepareContext(ths.context, sqlStr)
	} else {
		smt, err = ths.Connection.Prepare(sqlStr)
	}
	return smt, stackError.NewFromError(err, ths.stckErrorPowerId)
}
