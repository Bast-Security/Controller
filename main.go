package main

import (
	"github.com/grandcat/zeroconf"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"

	"database/sql"
	"log"
	"net"
)

var (
	db *sql.DB
)

func main() {
	var (
		mdnsServer *zeroconf.Server
		err error
	)

	log.Println("Starting MDNS server.")
	if ifaces, err := net.Interfaces(); err == nil {
		mdnsServer, err = zeroconf.Register("Bast Controller", SERVICE, DOMAIN, PORT, nil, ifaces)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer mdnsServer.Shutdown()

	log.Println("Opening SQL connection.")
	if db, err = sql.Open("mysql", "bast:bast@/bast"); err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	log.Println("Connecting to MQTT broker.")
	mqttClient := mqtt.NewClient(mqtt.NewClientOptions().AddBroker("tcp://localhost:1883"))
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	select {}
}

func cardValidate(card, door string) {

}

func pinValidate(pin, door string) {

}

func handlePin(client mqtt.Client, message mqtt.Message) {
	log.Println(string(message.Payload()))
}

func handleCard(client mqtt.Client, message mqtt.Message) {
	log.Println(string(message.Payload()))
}

