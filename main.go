package main

import (
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"fmt"
	"os"
	"os/signal"
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
		valid bool
		err error

		dbUser string
		dbPass string
		dbDB   string
	)

	if certFile, valid = os.LookupEnv("BAST_CERT"); !valid {
		certFile = "pki/bast.crt"
	}

	if keyFile, valid = os.LookupEnv("BAST_KEY"); !valid {
		keyFile = "pki/bast.key"
	}

	if addr, valid = os.LookupEnv("BAST_ADDR"); !valid {
		addr = ":8080"
	}

	if dbUser, valid = os.LookupEnv("BAST_DB_USER"); !valid {
		dbUser = "bast"
	}

	if dbPass, valid = os.LookupEnv("BAST_DB_PASS"); !valid {
		dbPass = "bast"
	}

	if dbDB, valid = os.LookupEnv("BAST_DB_DB"); !valid {
		dbDB = "bast"
	}

	if db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s?parseTime=true", dbUser, dbPass, dbDB)); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting REST server.")
	errChan := make(chan error)

	go func() {
		httpServer := server()
		errChan <- httpServer.ListenAndServeTLS(certFile,keyFile)
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

