package main

import (
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"net/http"
	"log"
)

var (
	db *sql.DB
	addr string
)

func main() {
	var (
		certFile string
		keyFile string
		ssl bool = true
		valid bool
		err error
	)

	if certFile, valid = os.LookupEnv("BAST_CERT"); valid {
		keyFile, ssl = os.LookupEnv("BAST_KEY")
	}

	if addr, valid = os.LookupEnv("BAST_ADDR"); !valid {
		addr = ":8080"
	}

	dbUser := os.Getenv("BAST_DB_USER")
	dbPass := os.Getenv("BAST_DB_PASS")
	dbDB := os.Getenv("BAST_DB_DB")

	if db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPass, dbDB)); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting REST server.")
	errChan := make(chan error)

	go func() {
		if ssl {
			httpServer := server()
			errChan <- httpServer.ListenAndServeTLS(certFile,keyFile)
		} else {
			errChan <- http.ListenAndServe(addr, router())
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)

	select {
	case err = <-errChan:
		log.Println("HTTP Server Closed with error: ", err)
		log.Println(certFile)
		log.Println(keyFile)
	case s := <-sig:
		log.Println("Received Signal: ", s)
	}

	db.Close()
}

