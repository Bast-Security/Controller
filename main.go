package main

import (
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"os"
	"os/signal"
	"net/http"
	"log"
	"flag"
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
		err error
	)

	flag.StringVar(&addr, "addr", ":8080", "Address to bind to, formatted as 'address:port'")
	flag.StringVar(&certFile, "cert-file", "", "Server certificate if using SSL")
	flag.StringVar(&keyFile, "key-file", "", "Server private key if using SSL")
	flag.Parse()

	if certFile == "" || keyFile == "" {
		ssl = false
	}

	if db, err = sql.Open("mysql", "bast:bast@/bast"); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting REST server.")
	if ssl {
		httpServer := server()
		go httpServer.ListenAndServeTLS(certFile,keyFile)
	} else {
		go http.ListenAndServe(addr, router())
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt)
	<-sig

	db.Close()
}

