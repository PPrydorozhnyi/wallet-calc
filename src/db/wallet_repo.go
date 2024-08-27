package db

import (
	"context"
	"github.com/PPrydorozhnyi/wallet/model"
)

const (
	selectWalletQuery = "SELECT wallets, version FROM accounts WHERE account_id = $1"
)

func GetWallet(id string) (*model.Account, error) {
	// TODO figure out which context should be used here
	a := &model.Account{}
	err := read(context.Background(), a, selectWalletQuery, id)

	return a, err
}
