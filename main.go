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
		name string
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

	if name, err = getSetting("name"); err != nil {
		log.Fatal(err)
	}

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
	router := chi.NewRouter()

	router.Post("/addUser", func(res http.ResponseWriter, req *http.Request) {
		var user User
		fmt.Fprintln(res, "/addUser")
		render.DecodeJSON(req.Body, &user)

		_, err := db.Exec(`INSERT INTO Users (name, email, pin, cardno) VALUES (?, ?, ?, ?);`, user.Name, user.Email, user.Pin, user.Card)

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

