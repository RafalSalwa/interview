package sql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type GormDB struct {
	*gorm.DB
}

func NewGormConnection(cfg MySQL) (*gorm.DB, error) {
	conString := fmt.Sprintf(dbString, cfg.Username, cfg.Password, cfg.Addr, cfg.DBName, dbParams)
	db, err := gorm.Open(mysql.Open(conString), &gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Info)})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
