package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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
	Id        string          `json:"id"`
	Type      string          `json:"type"`
	Vertical  string          `json:"vertical"`
	Amount    decimal.Decimal `json:"amount"`
	CreatedAt int64           `json:"createdTs"`
	UpdatedAt int64           `json:"updatedTs"`
}

type CommandResponse struct {
	Id          uuid.UUID     `json:"commandProcessingId"`
	Actions     []*OutcomeDto `json:"actions"`
	ProcessedAt int64         `json:"processedTs"`
}

type OutcomeDto struct {
	BalanceId    string          `json:"balanceId"`
	BalanceAfter decimal.Decimal `json:"balanceAfter"`
	Currency     string          `json:"unit"`
	ActionId     string          `json:"actionId"`
	BalanceType  *string         `json:"balanceType,omitempty"`
	Vertical     *string         `json:"vertical,omitempty"`
}
