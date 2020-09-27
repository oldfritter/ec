package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"ec/config"
)

type GormDB struct {
	*gorm.DB
	gdbDone bool
}

func (c *GormDB) DbCommit() {
	if c.gdbDone {
		return
	}
	tx := c.Commit()
	c.gdbDone = true
	if err := tx.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
}

func (c *GormDB) DbRollback() {
	if c.gdbDone {
		return
	}
	tx := c.Rollback()
	c.gdbDone = true
	if err := tx.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
}

func MainDbBegin() *GormDB {
	txn := config.MainDb.Begin()
	if txn.Error != nil {
		panic(txn.Error)
	}
	return &GormDB{txn, false}
}

func LogDbBegin() *GormDB {
	txn := config.LogDb.Begin()
	if txn.Error != nil {
		panic(txn.Error)
	}
	return &GormDB{txn, false}
}
