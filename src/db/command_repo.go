package db

import (
	"context"
	"github.com/PPrydorozhnyi/wallet/model"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/proto"
)

const (
	ledgerInsertQuery = `
						INSERT INTO ledgers (command_processing_id, account_id, ledger, command_id, created_ts, client_id, ref_command_processing_id, command_type)
                           VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
                           `
	walletUpdateQuery = `
						UPDATE accounts SET wallets = $3, version = $4
                        	WHERE account_id = $1 AND version = $2;
`
	walletInsertQuery = `
						INSERT INTO accounts (account_id, wallets, version) VALUES ($1, $2, 0);
`
	InitWalletVersion = -1
)

func PersistCommandResult(acc *model.Account, ledger *model.Ledger) error {
	// TODO figure out which context should be used here
	// https://github.com/jackc/pgx/issues/1223
	// context.WithTimeout(context.Background(), 120*time.Second)

	ledgerRecord, err := proto.Marshal(ledger.LedgerRecord)

	if err != nil {
		return err
	}

	walletState, err := proto.Marshal(acc.WalletState)

	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
	batch.Queue(ledgerInsertQuery, ledger.Id, acc.Id, ledgerRecord, ledger.CommandId,
		ledger.CreatedAt, ledger.ClientId, ledger.RefCommandProcessingId, ledger.CommandType)

	if acc.Version == InitWalletVersion {
		batch.Queue(walletInsertQuery, acc.Id, walletState)
	} else {
		batch.Queue(walletUpdateQuery, acc.Id, acc.Version, walletState, acc.Version+1)
	}

	return batchUpsert(context.Background(), batch)
}
