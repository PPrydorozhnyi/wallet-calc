package db

import (
	"context"
	"github.com/PPrydorozhnyi/wallet/model"
	wallet "github.com/PPrydorozhnyi/wallet/proto"
	"google.golang.org/protobuf/proto"
)

func GetWallet(id string) (*model.Account, error) {
	// todo read more about context
	ctx := context.Background()
	connection, err := getConnection(ctx)
	if err != nil {
		return nil, err
	}
	defer connection.Release()

	var w []byte
	var version int

	if err = connection.QueryRow(ctx, "SELECT wallets, version FROM accounts WHERE account_id = $1",
		id).Scan(&w, &version); err != nil {
		return nil, err
	}

	walletProto := &wallet.Wallet{}
	if err = proto.Unmarshal(w, walletProto); err != nil {
		return nil, err
	}

	return &model.Account{
		Id:          id,
		WalletState: walletProto,
		Version:     version,
	}, nil
}
