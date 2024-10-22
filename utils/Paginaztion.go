package utils

import (
	"hzx/goSqlHelper"
	"strconv"
)

var (
	//db *sql.DB = nil //全局数据库连接
	db *goSqlHelper.SqlHelper
)

func GetPagedDataTable(tableStr string, fieldStr string, whereStr string, sortStr string, pageIndex int, pageSize int) ([]goSqlHelper.HelperRow, int, int) {
	if db == nil {
		db = OpenDB()
	}
	isql := "IF OBJECT_ID(N'" + tableStr + "',N'U') IS NOT NULL \r\n"
	isql += "BEGIN \r\n"
	isql += "  SELECT * FROM (SELECT ROW_NUMBER() \r\n"
	isql += "  OVER(ORDER BY " + sortStr + ") AS rownum, " + fieldStr + " FROM " + tableStr + " WHERE " + whereStr + " ) AS DWHERE \r\n"
	isql += "  WHERE rownum BETWEEN CAST((" + strconv.Itoa(pageIndex) + "*" + strconv.Itoa(pageSize) + " + 1) as nvarchar(20)) \r\n"
	isql += "  AND cast(((" + strconv.Itoa(pageIndex) + "+1)*" + strconv.Itoa(pageSize) + ") as nvarchar(20)) \r\n"
	isql += "END \r\n"
	dt, err := db.QueryRows(isql)
	if err != nil {
		return nil, 0, 0
	}

	//isql = "IF OBJECT_ID(N'" + tableStr + "',N'U') IS NOT NULL \r\n"
	isql = "  SELECT COUNT(*) cnt FROM " + tableStr + " WHERE " + whereStr
	var cnt int = 0
	cnt, err = db.QueryScalarInt(isql)
	recordCount := cnt

	if err != nil {
		return nil, 0, 0
	}

	mod := recordCount % pageSize
	var pageCout int
	if mod == 0 {
		pageCout = recordCount / pageSize
	} else {
		pageCout = recordCount/pageSize + 1
	}
	return dt, pageCout, recordCount
}
