package sqldb

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type SqlInterface interface {
	DB() *sql.DB
}

type sqlStruct struct {
	db *sql.DB
}

func InitSql(driver string, host string, port string, username string, password string, database string) (SqlInterface, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, host, port, database)
	db, err := sql.Open(driver, connectionString)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)

	return &sqlStruct{
		db: db,
	}, nil

}

func (m *sqlStruct) DB() *sql.DB {
	return m.db
}
