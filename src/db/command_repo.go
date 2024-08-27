package db

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func ApplyTransaction() error {
	// TODO figure out which context should be used here
	// context.WithTimeout(context.Background(), 120*time.Second)

	batch := &pgx.Batch{}
	batch.Queue("insert into ledger(description, amount) values($1, $2)", "q1", 1)
	batch.Queue("insert into bedger(description, amount) values($1, $2)", "q2", 2)
	batch.Queue("update ledger set description = $1 where id = $2", "q3", -1)

	return batchUpsert(context.Background(), batch)
}
