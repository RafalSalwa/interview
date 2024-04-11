package sql

import (
	"database/sql"
	"fmt"
)

type MySQL struct {
	Addr     string `mapstructure:"addr"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbName"`
}
type DB struct {
	*sql.DB
}

const (
	driver   = "mysql"
	dbString = "%s:%s@tcp(%s)/%s?%s"
	dbParams = "parseTime=true&loc=Europe%2FWarsaw&charset=utf8&collation=utf8_polish_ci"
)

func NewMySQLConnection(c MySQL) (*DB, error) {
	con := fmt.Sprintf(dbString, c.Username, c.Password, c.Addr, c.DBName, dbParams)
	db, err := sql.Open(driver, con)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
