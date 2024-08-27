package model

import (
	wallet "github.com/PPrydorozhnyi/wallet/proto"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/proto"
)

type Account struct {
	Id          string         `json:"id,omitempty"`
	WalletState *wallet.Wallet `json:"walletState,omitempty"`
	Version     int            `json:"version,omitempty"`
}

func (a *Account) ScanRow(r pgx.Row) error {

	var w []byte

	if err := r.Scan(&a.Id, w, &a.Version); err != nil {
		return err
	}

	if err := proto.Unmarshal(w, a.WalletState); err != nil {
		return err
	}

	return nil
}
