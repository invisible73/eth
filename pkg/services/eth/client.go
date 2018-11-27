package eth

import "github.com/invisible73/eth/pkg/types"

// Client ...
type Client interface {
	SendTransaction(types.Transaction) (string, error)
	ListAccounts() (types.Accounts, error)
}
