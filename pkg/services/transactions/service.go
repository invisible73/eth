package transactions

import (
	"database/sql"
	"errors"

	"github.com/invisible73/eth/pkg/services/eth"
	"github.com/invisible73/eth/pkg/types"
)

type service struct {
	db     *sql.DB
	repo   Repo
	client eth.Client
}

func (s *service) Send(tx types.Transaction) (string, error) {
	hash, err := s.client.SendTransaction(tx)
	if err != nil {
		return "", err
	}

	tx.Hash = hash

	if err = s.repo.Save(tx); err != nil {
		return "", err
	}

	return hash, nil
}

func (s *service) GetLast() ([]types.Transaction, error) {
	sqltx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	data, err := s.repo.GetLast(sqltx)
	if err != nil {
		sqltx.Rollback()
		return nil, err
	}

	ids := make([]uint32, len(data))

	for i, tx := range data {
		ids[i] = tx.ID
	}

	if err = s.repo.SetSended(sqltx, ids); err != nil {
		sqltx.Rollback()
		return nil, err
	}

	if err = sqltx.Commit(); err != nil {
		return nil, err
	}

	return data, nil
}

// New ...
func New(db *sql.DB, client eth.Client, repo Repo) (Service, error) {
	if db == nil {
		return nil, errors.New("empty connection")
	}

	if repo == nil {
		return nil, errors.New("empty repo")
	}

	if client == nil {
		return nil, errors.New("empty client")
	}

	return &service{
		db:     db,
		repo:   repo,
		client: client,
	}, nil
}
