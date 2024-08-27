package model

import wallet "github.com/PPrydorozhnyi/wallet/proto"

type Account struct {
	Id          string         `json:"id,omitempty"`
	WalletState *wallet.Wallet `json:"walletState,omitempty"`
	Version     int            `json:"version,omitempty"`
}
