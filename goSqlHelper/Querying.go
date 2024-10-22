package goSqlHelper

import (
	"database/sql"

	"github.com/bobby96333/commonLib/stackError"
)

type Querying struct {
	rows         *sql.Rows
	_cols        []string
	stackErrorId int
}

func (ths *Querying) Close() {
	ths.rows.Close()
	//this._cols=nil
}

func (ths *Querying) Columns() ([]string, *stackError.StackError) {
	if ths._cols == nil {
		var err error
		ths._cols, err = ths.rows.Columns()
		if err != nil {
			return nil, stackError.NewFromError(err, ths.stackErrorId)
		}
	}
	return ths._cols, nil
}

func (ths Querying) QueryRow() (HelperRow, *stackError.StackError) {

	cols, err := ths.Columns()
	if err != nil {
		return nil, err
	}
	scanArgs := make([]interface{}, len(cols))
	values := make([]interface{}, len(cols))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	if ths.rows.Next() {

		err := ths.rows.Scan(scanArgs...)
		if err != nil {
			return nil, stackError.NewFromError(err, ths.stackErrorId)
		}
		record := make(HelperRow)
		for i, col := range values {
			if col == nil {
				continue
			}
			switch col.(type) {
			case []byte:
				record[cols[i]] = string(col.([]byte))
			default:
				record[cols[i]] = col
			}
		}

		return record, nil
	}
	return nil, NoFoundError
}

func (ths Querying) Scan(vals ...interface{}) *stackError.StackError {

	if ths.rows.Next() {
		err := ths.rows.Scan(vals...)
		if err != nil {
			return stackError.NewFromError(err, ths.stackErrorId)
		}
		return nil
	}
	return NoFoundError
}

func NewQuerying(rows *sql.Rows, stackErrorId int) *Querying {
	querying := new(Querying)
	querying.rows = rows
	querying.stackErrorId = stackErrorId
	return querying
}
