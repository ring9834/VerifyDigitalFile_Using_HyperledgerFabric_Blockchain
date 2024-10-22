package goSqlHelper

import (
	"github.com/bobby96333/commonLib/stackError"
)

/**
  orm read data
*/
func (ths *SqlHelper) QueryOrm(orm IEntity, sql string, args ...interface{}) *stackError.StackError {

	var err error
	rows, _ := ths.query(sql, args...)
	if err != nil {
		return stackError.NewFromError(err, ths.stckErrorPowerId)
	}
	defer rows.Close()
	if !rows.Next() {
		return NoFoundError
	}
	cols, err := rows.Columns()
	if err != nil {
		return stackError.NewFromError(err, ths.stckErrorPowerId)
	}
	vals := orm.MapFields(cols)
	for i, val := range vals {
		if val == nil {
			var empty interface{}
			vals[i] = &empty
		}
	}
	err = rows.Scan(vals...)

	if err != nil {

		return stackError.NewFromError(err, ths.stckErrorPowerId)
	}
	return nil
}

/*
execute insert sql
*/
func (ths *SqlHelper) OrmInsert(orm IEntity) (int64, *stackError.StackError) {
	sql := "INSERT INTO " + orm.TableName() + " SET "
	cols := orm.MapColumn()
	i := -1
	vals := make([]interface{}, len(cols))
	for key, val := range cols {
		i++
		if i > 0 {
			sql += ","
		}
		sql += key + "=?"
		vals[i] = val
	}
	if i < 0 {
		return 0, stackError.New("no found insert data")
	}
	return ths.ExecInsert(sql, vals...)
}

func (ths *SqlHelper) OrmDelete(orm IEntity) (int64, *stackError.StackError) {
	sql := "DELETE FROM " + orm.TableName() + " WHERE "
	cols := orm.MapColumn()
	keys := orm.PrimaryKeys()
	if len(keys) < 0 {
		return 0, stackError.New("no found insert data")
	}

	vals := make([]interface{}, len(keys))
	for i, key := range keys {
		if i > 0 {
			sql += " AND "
		}
		sql += key + "=?"
		vals[i] = cols[key]
	}
	return ths.ExecUpdateOrDel(sql, vals...)
}

func (ths *SqlHelper) OrmUpdate(orm IEntity) (int64, *stackError.StackError) {
	sql := "UPDATE " + orm.TableName() + " SET "

	cols := orm.MapColumn()
	keys := orm.PrimaryKeys()
	if len(keys) < 0 {
		return 0, stackError.New("no found insert data")
	}
	vals := make([]interface{}, 0, len(cols))

	i := -1
	for key, val := range cols {
		has := false
		for _, primaryKey := range keys {
			if primaryKey == key {
				has = true
				break
			}
		}
		if has {
			continue //priimarykey
		}
		i++
		if i > 0 {
			sql += " ,"
		}
		sql += key + "=?"
		vals = append(vals, val)
	}
	sql += " WHERE "
	for i, key := range keys {
		if i > 0 {
			sql += " AND "
		}
		sql += key + "=?"
		vals = append(vals, cols[key])
	}
	return ths.ExecUpdateOrDel(sql, vals...)
}
