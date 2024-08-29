package model

import (
	"github.com/shopspring/decimal"
)

type TransactionRequest struct {
	CommandId string    `json:"commandId"`
	Actions   *[]Action `json:"actions"`
	Reason    *Reason   `json:"reason"`
}

type Action struct {
	Currency        string          `json:"unit"`
	BalanceId       string          `json:"balanceId"`
	TransactionType string          `json:"direction"` // todo add validation
	Amount          decimal.Decimal `json:"amount"`
	AllowNegative   bool            `json:"allowNegative"`
}
