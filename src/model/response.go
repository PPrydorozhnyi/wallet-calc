package model

import (
	"github.com/google/uuid"
	"math/big"
)

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

type TransactionResponse struct {
	Id          uuid.UUID     `json:"commandProcessingId"`
	Actions     []*OutcomeDto `json:"actions"`
	ProcessedAt int64         `json:"processedTs"`
}

type OutcomeDto struct {
	BalanceId    string     `json:"balanceId"`
	BalanceAfter *big.Float `json:"balanceAfter"`
	Currency     string     `json:"unit"`
	ActionId     string     `json:"actionId"`
}
