package main

import (
	"github.com/grandcat/zeroconf"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"database/sql"
	"log"
	"net"
	"net/http"
	"fmt"
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
		log.Println(token.Error())
	}

	log.Println("Starting REST server.")
	router := chi.NewRouter()

	router.Post("/addUser", func(res http.ResponseWriter, req *http.Request) {
		var user User
		fmt.Fprintln(res, "/addUser")
		render.DecodeJSON(req.Body, &user)

		_result, err := db.Exec(`INSERT INTO Users (name, email, pin, cardno) VALUES (?, ?, ?, ?);`, user.Name, user.Email, user.Pin, user.Card)

		if err != nil {
			res.WriteHeader(400)
		} else {
			res.WriteHeader(200)
		}
	})

	router.Post("/addRole", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(res, "/addRole")
	})

	router.Post("/addLock", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(res, "/addLock")
	})

	router.Get("/listUsers", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(res, "/listUsers")
	})

	router.Get("/listRoles", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(res, "/listRoles")
	})

	router.Get("/listLocks", func(res http.ResponseWriter, req *http.Request) {
		fmt.Fprint(res, "/listLocks")
	})

	http.ListenAndServe(":8080", router)
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

