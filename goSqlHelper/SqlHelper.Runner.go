package goSqlHelper

import (
	"context"
	"database/sql"

	"github.com/bobby96333/commonLib/stackError"
)

type SqlHelperRunner struct {
	SqlHelper
}

func (ths *SqlHelperRunner) SetContext(ctx context.Context) {
	ths.context = ctx
}

func (ths *SqlHelperRunner) BeginTx(ctx context.Context, opts *sql.TxOptions) *stackError.StackError {
	tx, err := ths.Connection.BeginTx(ctx, opts)
	if err != nil {
		return stackError.NewFromError(err, ths.stckErrorPowerId)
	}
	ths.tx = tx
	return nil
}

func (ths *SqlHelperRunner) Begin() *stackError.StackError {
	tx, err := ths.Connection.Begin()
	if err != nil {
		return stackError.NewFromError(err, ths.stckErrorPowerId)
	}
	ths.tx = tx
	return nil
}
