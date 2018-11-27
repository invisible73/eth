package postgresql

import (
	"database/sql"
	"errors"
	"math/big"
	"time"

	"github.com/invisible73/eth/pkg/services/transactions"
	"github.com/invisible73/eth/pkg/types"
	"github.com/lib/pq"
)

type repo struct {
	db *sql.DB
}

func (r *repo) SetSended(sqltx *sql.Tx, ids []uint32) error {
	stmt, err := r.db.Prepare(`UPDATE transactions SET sended = true WHERE id = ANY ($1)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(pq.Array(ids))

	return err
}

func (r *repo) GetLast(sqltx *sql.Tx) ([]types.Transaction, error) {
	rows, err := r.db.Query(`
		SELECT id, date, "to", amount, confirmations
		FROM transactions
		WHERE confirmations < 3
		AND NOT invalidated
		AND not sended
		FOR UPDATE
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		date              time.Time
		result            []types.Transaction
		to, amount        string
		id, confirmations uint32
	)

	for rows.Next() {
		err = rows.Scan(&id, &date, &to, &amount, &confirmations)
		if err != nil {
			return nil, err
		}

		a := new(big.Int)
		a.SetString(amount, 10)

		result = append(result, types.Transaction{
			ID:            id,
			Date:          date,
			To:            to,
			Amount:        a,
			Confirmations: confirmations,
		})
	}

	return result, nil
}

func (r *repo) Save(tx types.Transaction) error {
	stmt, err := r.db.Prepare(`INSERT INTO transactions("from", "to", "amount", "hash") VALUES($1, $2, $3, $4)`)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(tx.From, tx.To, tx.Amount.String(), tx.Hash)

	return err
}

// New ...
func New(db *sql.DB) (transactions.Repo, error) {
	if db == nil {
		return nil, errors.New("empty connection")
	}

	return &repo{db: db}, nil
}
