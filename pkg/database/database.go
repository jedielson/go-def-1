package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBManager interface {
	InitDb() error
	GetDB() *gorm.DB
	Close() error
}

type dbManager struct {
	DBConn *gorm.DB
	Dsn    string
}

func NewDbManager(dsn string) DBManager {
	man := &dbManager{
		Dsn:    dsn,
		DBConn: nil,
	}

	return man
}

func (bd *dbManager) InitDb() (err error) {
	bd.DBConn, err = gorm.Open(sqlite.Open(bd.Dsn), &gorm.Config{})
	return err
}

func (bd *dbManager) Close() (err error) {
	db, err := bd.DBConn.DB()

	if err != nil {
		return
	}

	defer func() {
		err = db.Close()
	}()

	return
}

func (bd *dbManager) GetDB() *gorm.DB {
	return bd.DBConn
}
