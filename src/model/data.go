package model

import (
	wallet "github.com/PPrydorozhnyi/wallet/proto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/proto"
	"time"
)

type Account struct {
	Id          string         `json:"id,omitempty"`
	WalletState *wallet.Wallet `json:"walletState,omitempty"`
	Version     int            `json:"version,omitempty"`
}

func (a *Account) ScanRow(r pgx.Row) error {

	var w []byte

	if err := r.Scan(&w, &a.Version); err != nil {
		return err
	}

	walletProto := &wallet.Wallet{}

	if err := proto.Unmarshal(w, walletProto); err != nil {
		return err
	}

	a.WalletState = walletProto

	return nil
}

type Ledger struct {
	Id                     uuid.UUID            `json:"commandProcessingId"`
	AccountId              string               `json:"accountId"`
	LedgerRecord           *wallet.LedgerRecord `json:"ledgerRecord"`
	CreatedAt              time.Time            `json:"createdAt"`
	CommandId              string               `json:"commandId"`
	RefCommandProcessingId *uuid.UUID           `json:"refCommandProcessingId"`
	ClientId               int                  `json:"clientId"`
	CommandType            string               `json:"commandType"`
}
