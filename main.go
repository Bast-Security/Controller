package main

import (
	"github.com/grandcat/zeroconf"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"log"
	"net"
	"fmt"
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
	if db, err = sql.Open("mysql", "bast:bast@/bast"); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Loading system settings.")
	if name, err = getSetting("name"); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting MDNS server.")
	if ifaces, err := net.Interfaces(); err == nil {
		mdnsServer, err = zeroconf.Register(name, SERVICE, DOMAIN, PORT, nil, ifaces)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer mdnsServer.Shutdown()

	log.Println("Connecting to MQTT broker.")
	mqttClient := mqtt.NewClient(mqtt.NewClientOptions().AddBroker("tcp://localhost:1883"))
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	} else {
		if token := mqttClient.Subscribe(fmt.Sprintf("bast/%s/+/pin", name), 0, handlePin); token.Wait() && token.Error() != nil {
			log.Println(token.Error())
		}

		if token := mqttClient.Subscribe(fmt.Sprintf("bast/%s/+/card", name), 0, handleCard); token.Wait() && token.Error() != nil {
			log.Println(token.Error())
		}
	}

	log.Println("Starting REST server.")
	httpServer := server()
	httpServer.ListenAndServeTLS("pki/bast-root.crt", "pki/bast-root.key")
}

