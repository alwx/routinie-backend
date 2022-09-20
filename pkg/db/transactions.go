package db

import (
	"gorm.io/gorm"
	"habiko-go/pkg/models"
)

type TransactionProvider struct {
	DBConn          *gorm.DB
	TimescaleDBConn *gorm.DB
}

type Transaction struct {
	tx *gorm.DB
}

func (p *TransactionProvider) Create() (models.Transaction, error) {
	tx := p.DBConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &Transaction{tx: tx}, nil
}

func (t *Transaction) GetTx() *gorm.DB {
	return t.tx
}

func (t *Transaction) Commit() error {
	return t.tx.Commit().Error
}

func (t *Transaction) Rollback() error {
	return t.tx.Rollback().Error
}
