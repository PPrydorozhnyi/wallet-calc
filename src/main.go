package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"wallet/db"
	"wallet/handler"
)

//TIP Simple web application

var (
	addr = os.Getenv("PORT")
)

// init suits mainly for verifying correctness of the vars
func init() {
	if addr == "" {
		addr = ":8081"
	}
}

func main() {

	start := time.Now()

	initConnectionPool()
	defer db.CloseConnectionPool()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Handle)
	mux.HandleFunc("/posts", handler.PostsHandle)
	mux.HandleFunc("/posts/", handler.PostHandle)

	log.Printf("Starting Server on port %s\n", addr)

	go testReadiness(start, addr)

	startServer(mux)
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

func startServer(mux *http.ServeMux) {
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
