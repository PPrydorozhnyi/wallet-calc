package main

import (
	"errors"
	"github.com/PPrydorozhnyi/wallet/db"
	"github.com/PPrydorozhnyi/wallet/handler"
	"github.com/PPrydorozhnyi/wallet/util"
	"log"
	"net"
	"net/http"
	"time"
)

//TIP Simple web application

func main() {

	start := time.Now()

	initConnectionPool()
	defer db.CloseConnectionPool()

	addr := util.GetStringEnv("SERVER_PORT", ":8081")
	log.Printf("Starting Server on port %s\n", addr)

	go testReadiness(start, addr)

	mux := http.NewServeMux()
	initMappings(mux)
	startServer(mux, addr)
}

func initMappings(mux *http.ServeMux) {
	mux.HandleFunc("/", handler.Handle)
	mux.HandleFunc("/posts", handler.PostsHandle)
	mux.HandleFunc("/accounts/", handler.AccountHandle)

	//wallets
	mux.HandleFunc("/api/v1/accounts/{accountId}", handler.WalletsHandle)

	//commands
	mux.HandleFunc("/api/v1/accounts/{accountId}/transactions", handler.TransactionsHandle)
}

func initConnectionPool() {
	start := time.Now()

	err := db.CreateConnectionPool()
	if err != nil {
		log.Fatalf("Failed to create connection pool %s\n", err)
	}

	err = db.HealthCheck()
	if err != nil {
		log.Fatalf("Database is not reachable at the very beginning %s\n", err)
	}

	log.Printf("Connection pool is initialized in %s\n", time.Since(start))
}

func startServer(mux *http.ServeMux, addr string) {
	err := http.ListenAndServe(addr, mux)

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed.")
	} else if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}

func testReadiness(startTime time.Time, addr string) {
	for {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			e := conn.Close()
			if e != nil {
				log.Println("Cannot close test connection")
			}

			log.Printf("Server is ready to accept requests at %s. Started in %s", time.Now().Format(time.RFC3339),
				time.Since(startTime))
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}
