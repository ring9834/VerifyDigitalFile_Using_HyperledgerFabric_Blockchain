package goSqlHelper

type HelperTable struct {
	rows    []HelperRow
	columns []string
}

func NewTable(rows []HelperRow, columns []string) *HelperTable {

	helper := HelperTable{
		rows:    rows,
		columns: columns,
	}
	return &helper

}

func (ths *HelperTable) Rows() []HelperRow {
	return ths.rows
}
func (ths *HelperTable) Columns() []string {
	return ths.columns
}
