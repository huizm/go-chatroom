package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	dbUsername   = ""
	dbPassword   = ""
	dbHost       = "localhost:3306"
	dbName       = ""
	maxIdleConns = 10
	maxOpenConns = 30
)

var DBEngine *gorm.DB

func ShouldLoadDB() error {
	var err error
	if DBEngine, err = newDBEngine(); err != nil {
		return err
	}

	if err = DBEngine.AutoMigrate(
		&User{},
	); err != nil {
		return err
	}

	return nil
}

func newDBEngine() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUsername, dbPassword, dbHost, dbName)

	cfg := &gorm.Config{Logger: logger.Default.LogMode(logger.Info)}

	db, err := gorm.Open(mysql.Open(dsn), cfg)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)

	return db, nil
}
