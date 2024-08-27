package model

import "math/big"

type AccountResponse struct {
	AccountId string                `json:"accountId"`
	Wallets   map[string]*WalletDto `json:"wallets"`
}

type WalletDto struct {
	Currency string        `json:"unit"`
	Balances []*BalanceDto `json:"balances"`
}

type BalanceDto struct {
	Id        string     `json:"id"`
	Type      string     `json:"type"`
	Vertical  string     `json:"vertical"`
	Amount    *big.Float `json:"amount"`
	CreatedAt uint64     `json:"createdTs"`
	UpdatedAt uint64     `json:"updatedTs"`
}
