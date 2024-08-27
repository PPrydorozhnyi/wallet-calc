package db

import (
	"context"
	"errors"
	"github.com/PPrydorozhnyi/wallet/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"sync"
	"time"
)

type RowScanner interface {
	ScanRow(r pgx.Row) error
}

var (
	connectionPool   *pgxpool.Pool
	connectionPoolMu sync.Mutex
)

func config() (*pgxpool.Config, error) {
	connectionString := util.GetStringEnv("DATABASE_URL", "postgres://wc_user:wc_password@localhost:5432/wc")
	dbConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	dbConfig.MaxConns = int32(util.GetIntEnv("DATABASE_CONNECTIONS_MAX", 10))
	dbConfig.MinConns = int32(util.GetIntEnv("DATABASE_CONNECTIONS_MIN", 10))
	dbConfig.MaxConnLifetime = util.GetDurationEnv("DATABASE_CONNECTIONS_MAX_LIFETIME", time.Hour)
	dbConfig.MaxConnIdleTime = util.GetDurationEnv("DATABASE_CONNECTIONS_MAX_IDLE", 30*time.Minute)
	dbConfig.HealthCheckPeriod = util.GetDurationEnv("DATABASE_CONNECTIONS_HEALTH_CHECK_PERIOD", 5*time.Minute)
	dbConfig.ConnConfig.ConnectTimeout = util.GetDurationEnv("DATABASE_CONNECTIONS_CONN_TIMEOUT", 20*time.Second)

	dbConfig.BeforeClose = func(c *pgx.Conn) {
		log.Println("Closed the connection to the database")
	}

	return dbConfig, nil
}

func CreateConnectionPool() error {
	connectionPoolMu.Lock()
	defer connectionPoolMu.Unlock()

	if connectionPool != nil {
		return nil
	}

	dbConfig, err := config()
	if err != nil {
		return err
	}

	connPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)

	if err != nil {
		return err
	}
	connectionPool = connPool

	return nil
}

func CloseConnectionPool() {
	connectionPoolMu.Lock()
	defer connectionPoolMu.Unlock()

	if connectionPool != nil {
		connectionPool.Close()
		connectionPool = nil
	}
}

func HealthCheck() error {
	if connectionPool == nil {
		return errors.New("connectionPool is not initialized")
	}

	connection, err := connectionPool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer connection.Release()

	err = connection.Ping(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func read(ctx context.Context, rs RowScanner, sql string, args ...any) error {
	connection, err := connectionPool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer connection.Release()

	if row := connection.QueryRow(ctx, sql, args); rs.ScanRow(row) != nil {
		return err
	}

	return nil
}
