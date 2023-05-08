package sql

import (
	"fmt"
	"github.com/RafalSalwa/interview-app-srv/config"
	"github.com/RafalSalwa/interview-app-srv/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type GormDB struct {
	*gorm.DB
}

func NewUsersDBGorm(c config.ConfDB, l *logger.Logger) *gorm.DB {
	conString := fmt.Sprintf(dbString, c.Username, c.Password, c.Host, c.Port, c.DBName, dbParams)
	db, err := gorm.Open(mysql.Open(conString), &gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Info)})
	fmt.Println(conString)
	if err != nil {
		l.Fatal().Err(err).Msg("DB connection start failure")
	}
	sqlDB, err := db.DB()

	err = sqlDB.Ping()
	if err != nil {
		l.Error().Err(err)
	}
	return db
}
