package goSqlHelper

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/bobby96333/commonLib/stackError"
)

type SqlHelper struct {
	Connection       *sql.DB
	context          context.Context
	tx               *sql.Tx
	debugMod         bool
	stckErrorPowerId int
}

const QUERY_BUFFER_SIZE = 20

/**
@todo no sql

	var obj=new(tb1)
	con.Insert(obj)
	obj.setup(conn)
	obj.Select("id,val").Where("id=2").QueryList()
	sqlHelper.Select("id,val").Where("id=2").QueryList()

*/
func MssqlOpen(driver string, connectionStr string) (*SqlHelper, *stackError.StackError) {

	sqlHelper := new(SqlHelper)
	err := sqlHelper.Init(driver, connectionStr)
	if err != nil {
		return nil, err
	}
	return sqlHelper, nil
}

// func New(connectionStr string) (*SqlHelper, *stackError.StackError) {
// 	return MssqlOpen(connectionStr)
// }

/**
  open db
*/
func (ths *SqlHelper) Init(driver, connectionStr string) *stackError.StackError {
	if DefaultDebugModel {
		ths.OpenDebug()
	} else {
		ths.stckErrorPowerId = -1
	}

	var err error
	//	sql.Open
	ths.Connection, err = sql.Open(driver, connectionStr)
	if err != nil {
		return stackError.New(fmt.Sprintf("db connected failed:%s", err.Error()))
	}
	err = ths.Connection.Ping()
	if err != nil {
		return stackError.NewFromError(err, ths.stckErrorPowerId)
	}
	return nil
}

/**
begin context
*/
func (ths *SqlHelper) BeginContext(ctx context.Context) *SqlHelperRunner {
	runner := new(SqlHelperRunner)
	runner.SetDB(ths.Connection)
	runner.SetContext(ctx)
	return runner
}

/**
begin a trasnaction
*/
func (ths *SqlHelper) Begin() *SqlHelperRunner {
	runner := new(SqlHelperRunner)
	runner.SetDB(ths.Connection)
	runner.Begin()
	return runner
}

/**
print sql and parameter at prepare exeucting
*/
func (ths *SqlHelper) OpenDebug() {
	ths.debugMod = true
	//this.stckErrorPowerId = stackError.GetPowerKey()
	//stackError.SetPower(true, this.stckErrorPowerId)

}

/**
begin a trasnaction
*/
func (ths *SqlHelper) BeginTx(ctx context.Context, opts *sql.TxOptions) (*SqlHelperRunner, *stackError.StackError) {
	runner := new(SqlHelperRunner)
	runner.SetDB(ths.Connection)
	err := runner.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return runner, nil
}

/**
set db object
*/
func (ths *SqlHelper) SetDB(conn *sql.DB) {
	ths.Connection = conn
}

/**
get Querying handler
*/
func (ths *SqlHelper) Querying(sql string, args ...interface{}) (*Querying, *stackError.StackError) {
	if ths.debugMod {
		fmt.Println(sql)
		fmt.Println(args...)
	}
	var rows, err = ths.query(sql, args...)
	if err != nil {
		return nil, err
	}
	querying := NewQuerying(rows, ths.stckErrorPowerId)
	return querying, nil
}

/**
  read a int value
*/
func (ths *SqlHelper) QueryScalar(val interface{}, sql string, args ...interface{}) *stackError.StackError {
	if ths.debugMod {
		fmt.Println(sql)
		fmt.Println(args...)
	}
	var err error
	rows, _ := ths.query(sql, args...)
	if err != nil {
		return stackError.NewFromError(err, ths.stckErrorPowerId)
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(val)
		return stackError.NewFromError(err, ths.stckErrorPowerId)
	}
	return NoFoundError
}

/**
  read a int value
*/
func (ths *SqlHelper) QueryScalarInt(sql string, args ...interface{}) (int, *stackError.StackError) {
	var val int
	err := ths.QueryScalar(&val, sql, args...)
	return val, err
}

/**
  read a int value
*/
func (ths *SqlHelper) QueryScalarString(sql string, args ...interface{}) (string, *stackError.StackError) {
	var val string
	err := ths.QueryScalar(&val, sql, args...)
	return val, err
}

/*
execute sql
*/
func (ths *SqlHelper) Exec(sql string, args ...interface{}) (sql.Result, *stackError.StackError) {
	if ths.debugMod {
		fmt.Println(sql)
		fmt.Println(args...)
	}
	var err error
	stmt, _ := ths.prepare(sql)
	if err != nil {
		return nil, stackError.NewFromError(err, ths.stckErrorPowerId)
	}
	defer stmt.Close()
	result, err := stmt.Exec(args...)
	if err != nil {
		return nil, stackError.NewFromError(err, ths.stckErrorPowerId)
	}
	return result, nil
}

/*
execute insert sql
*/
func (ths *SqlHelper) ExecInsert(sql string, args ...interface{}) (int64, *stackError.StackError) {
	result, err := ths.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()
	// if err2 != nil { //因为mssql暂不支持LastInsertId(),so slashed the three lines 20210528 by hhy
	// 	return 0, stackError.NewFromError(err2, ths.stckErrorPowerId)
	// }
	return id, nil
}

/*
execute update or delete sql
*/
func (ths *SqlHelper) ExecUpdateOrDel(sql string, args ...interface{}) (int64, *stackError.StackError) {
	result, err := ths.Exec(sql, args...)
	if err != nil {
		return 0, err
	}

	cnt, err2 := result.RowsAffected()
	if err2 != nil {
		return 0, stackError.NewFromError(err2, ths.stckErrorPowerId)
	}
	return cnt, nil
}

/*
   close db pool
*/
func (ths *SqlHelper) Close() *stackError.StackError {
	err := ths.Connection.Close()
	return stackError.NewFromError(err, ths.stckErrorPowerId)
}

// get auto sql
func (ths *SqlHelper) Auto() *AutoSql {
	return NewAutoSql(ths)
}
