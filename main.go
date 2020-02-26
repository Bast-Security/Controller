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

func cardValidate(card, door string) (valid bool) {
	valid = false

	rows, err := db.Query(`SELECT Users.name
		FROM Users
		INNER JOIN UserRole ON UserRole.userid = Users.id
		INNER JOIN Permissions ON Permissions.role = UserRole.role
		INNER JOIN AuthTypes ON Permissions.door = AuthTypes.door
		WHERE Permissions.door = ?
		AND Users.cardno = ?
		AND (AuthTypes.authType = 2 OR AuthTypes.authType = -3)`, door, card)

	if err != nil {
		log.Println(err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var user string
			if err := rows.Scan(&user); err != nil {
				log.Println(err)
				return
			}
			valid = true
		}
		if err := rows.Err(); err != nil {
			log.Println(err)
		}
	}
	return
}

func pinValidate(pin, door string) (valid bool) {
	valid = false

	rows, err := db.Query(`SELECT Users.name
		FROM Users
		INNER JOIN UserRole ON UserRole.userid = Users.id
		INNER JOIN Permissions ON Permissions.role = UserRole.role
		INNER JOIN AuthTypes ON Permissions.door = AuthTypes.door
		WHERE Permissions.door = ?
		AND Users.pin = ?
		AND (AuthTypes.authType = 1 OR AuthTypes.authType = -3)`, door, pin)

	if err != nil {
		log.Println(err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var user string
			if err := rows.Scan(&user); err != nil {
				log.Println(err)
				return
			}
			valid = true
		}
		if err := rows.Err(); err != nil {
			log.Println(err)
		}
	}
	return
}

func handlePin(client mqtt.Client, message mqtt.Message) {
	log.Println(string(message.Payload()))
}

func handleCard(client mqtt.Client, message mqtt.Message) {
	log.Println(string(message.Payload()))
}

