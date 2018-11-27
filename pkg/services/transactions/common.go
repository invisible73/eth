package transactions

import (
	"database/sql"

	"github.com/invisible73/eth/pkg/types"
)

// Service ...
type Service interface {
	Send(types.Transaction) (string, error)
	GetLast() ([]types.Transaction, error)
}

// Repo ...
type Repo interface {
	Save(types.Transaction) error
	GetLast(*sql.Tx) ([]types.Transaction, error)
	SetSended(*sql.Tx, []uint32) error
}
