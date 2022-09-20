package models

import (
	"errors"
)

type TransactionProvider interface {
	Create() (Transaction, error)
}

type Transaction interface {
	Commit() error
	Rollback() error
}

var ErrDb = errors.New("cannot access the repository")
var ErrDbObject = errors.New("cannot access or change data in the repository")
