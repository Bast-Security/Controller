package main

import (
	"github.com/grandcat/zeroconf"
	_ "github.com/mattn/go-sqlite3"

	"database/sql"
	"io/ioutil"
	"log"
	"net"
)

var (
	db *sql.DB
	name string
)

func main() {
	var (
		mdnsServer *zeroconf.Server
		err error
	)

	log.Println("Opening SQL connection.")
	if db, err = sql.Open("sqlite3", "./bast.db"); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for dbOk := false; !dbOk; {
		log.Println("Loading system settings.")
		if name, err = getSetting("name"); err != nil {
			log.Println(err)
			log.Println("Creating initial database")

			query, err := ioutil.ReadFile("migrations/create_tables.sql")
			if err != nil {
				log.Fatal("Failed to read sql script", err)	
			}

			if _, err := db.Exec(string(query)); err != nil {
				log.Fatal("Failed to execute sql script", err)
			}
		} else {
			dbOk = true
		}
	}

	log.Println("Starting MDNS server.")
	if ifaces, err := net.Interfaces(); err == nil {
		mdnsServer, err = zeroconf.Register(name, SERVICE, DOMAIN, PORT, nil, ifaces)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer mdnsServer.Shutdown()

	log.Println("Starting REST server.")
	httpServer := server()
	httpServer.ListenAndServeTLS("pki/bast-root.crt", "pki/bast-root.key")
}

