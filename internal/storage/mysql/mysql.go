package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func CreateConnection(mysqlConnStr string) (*sqlx.DB, error) {
	mysqlConn, err := sqlx.Connect("mysql", mysqlConnStr)
	if err != nil {
		return nil, err
	}

	return mysqlConn, nil
}
