package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"sync"
	"time"
)

const (
	defaultMaxConns          = int32(250)
	defaultMinConns          = int32(10)
	defaultMaxConnLifetime   = time.Hour
	defaultMaxConnIdleTime   = time.Minute * 30
	defaultHealthCheckPeriod = time.Minute
	defaultConnectTimeout    = time.Second * 5
)

var (
	connectionString = os.Getenv("DATABASE_URL")
	connectionPool   *pgxpool.Pool
	connectionPoolMu sync.Mutex
)

func init() {
	if connectionString == "" {
		connectionString = "postgres://wc_user:wc_password@localhost:5432/wc"
	}
}

func config() (*pgxpool.Config, error) {
	dbConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	// todo add possibility to fetch configs from env variables
	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

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
		err := CreateConnectionPool()
		if err != nil {
			return err
		}
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
