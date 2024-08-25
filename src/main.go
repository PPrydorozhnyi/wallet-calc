package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"time"
	"wallet/handler"
)

//TIP Simple web application

func main() {

	addr := os.Getenv("PORT")
	if addr == "" {
		addr = ":8081"
	}

	start := time.Now()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.Handle)
	mux.HandleFunc("/posts", handler.PostsHandle)
	mux.HandleFunc("/posts/", handler.PostHandle)

	log.Printf("Starting Server on port %s\n", addr)

	go testReadiness(start, addr)

	err := http.ListenAndServe(addr, mux)

	if errors.Is(err, http.ErrServerClosed) {
		log.Println("Server closed.")
	} else if err != nil {
		log.Fatalf("Error starting server: %s", err)
	} else {
		log.Printf("Started server in %s.\n", time.Since(start))
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

			log.Printf("Server is ready to accept connections at %s. Started in %s", time.Now().Format(time.RFC3339),
				time.Since(startTime))
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
}
