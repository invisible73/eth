package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"
)

// Transaction ...
type Transaction struct {
	ID            uint32
	Date          time.Time
	From          string
	To            string
	Amount        *big.Int
	Hash          string
	BlockHash     string
	BlockNumber   uint32
	Confirmations uint32
}

// MarshalJSON ...
func (t Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.newJSON())
}

type jsonTransaction struct {
	Date          time.Time `json:"date"`
	Address       string    `json:"address"`
	Amount        string    `json:"amount"`
	Confirmations uint32    `json:"confirmations"`
}

func (t Transaction) newJSON() jsonTransaction {
	return jsonTransaction{
		Date:          t.Date,
		Address:       t.To,
		Amount:        t.toEth(),
		Confirmations: t.Confirmations,
	}
}

func (t Transaction) toEth() string {
	s := t.Amount.String()
	if len(s) <= 18 {
		return fmt.Sprintf("0.%018s", s)
	}
	div := s[:len(s)-18]
	mod := s[len(s)-18:]

	return fmt.Sprintf("%s.%s", div, mod)
}
