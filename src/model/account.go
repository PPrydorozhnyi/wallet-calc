package model

import wallet "github.com/PPrydorozhnyi/wallet/proto"

type Account struct {
	Id      string         `json:"id,omitempty"`
	Wallets *wallet.Wallet `json:"wallets,omitempty"`
	Version int            `json:"version,omitempty"`
}
