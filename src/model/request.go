package model

import (
	"github.com/shopspring/decimal"
)

type command struct {
	CommandId string  `json:"commandId"`
	Reason    *Reason `json:"reason"`
}

type CreateWalletRequest struct {
	command
	AccountId string              `json:"accountId"`
	Wallets   []*WalletDefinition `json:"wallets"`
}

type WalletDefinition struct {
	Currency string               `json:"unit"`
	Balances []*BalanceDefinition `json:"balances"`
}

type BalanceDefinition struct {
	BalanceType string          `json:"type"`
	Vertical    string          `json:"vertical"`
	Amount      decimal.Decimal `json:"amount"`
}

type TransactionRequest struct {
	command
	Actions *[]Action `json:"actions"`
}

type Action struct {
	Currency        string          `json:"unit"`
	BalanceId       string          `json:"balanceId"`
	TransactionType string          `json:"direction"` // todo add validation
	Amount          decimal.Decimal `json:"amount"`
	AllowNegative   bool            `json:"allowNegative"`
}
