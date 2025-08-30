package database

import (
	"gorm.io/gorm"
)

type TransactionManager interface {
	RunTransaction(fn func(tx *gorm.DB) error) error
}

type TransactionDBWrapper struct {
	*gorm.DB
}

func (db *TransactionDBWrapper) RunTransaction(fn func(tx *gorm.DB) error) error {
	return db.DB.Transaction(fn)
}

type Transactional interface {
	WithTx(tx *gorm.DB) interface{}
}
