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
		//variable that will be added into the database
		var role string

		fmt.Fprintln(res, "/addRole")
		render.DecodeJSON(req.Body, &role)

		_, err := db.Exec(`INSERT INTO Roles (name) VALUES (?);`, role)

		if err != nil{
			res.WriteHeader(400)
		}else{
			res.WriteHeader(200)
		}
	})

	router.Post("/addLock", func(res http.ResponseWriter, req *http.Request) {
		//variable that will be added into the database
		var door string

		fmt.Fprintln(res, "/addLock")
		render.DecodeJSON(req.Body, &door)

		_, err := db.Exec(`INSERT INTO Door (name) VALUES (?);`, door)

		if err != nil{
			res.WriteHeader(400)
		}else{
			res.WriteHeader(200)
		}
	})


	router.Get("/listUsers", func(res http.ResponseWriter, req *http.Request) {
		//array that will save the each user from the database
		var users []User

		//variable will save the querry command
		rows, err := db.Query(`select Users.id, Users.name, Users.email, Users.pin, Users.cardno from Users`)
		
		//if statement makes to sure that query was a success; if successful then each row in the Users scheme is read
		if err != nil {
			log.Println(err)
		} else {
			defer rows.Close()
			for rows.Next() {
				//variable to save the user
				var user User

				if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Pin, &user.Card); err != nil {
					log.Println(err)
					return
				}
				
				//a new user is added to the array
				users = append(users, user)
			}
			if err := rows.Err(); err != nil {
				log.Println(err)
			}
		}
		
		//converts array into a JSON and sends it to requestor
		render.JSON(res, req, users)
	})

	router.Get("/listRoles", func(res http.ResponseWriter, req *http.Request) {
		//array that will save the each role from the database
		var roles []Role

		//variable will save the querry command
		rows, err := db.Query(`select Roles.name from Roles`)

		//if statement makes to sure that query was a success; if successful then each row in the Roles scheme is read
		if err != nil {
			log.Println(err)
		} else {
			defer rows.Close()
			for rows.Next() {
				//variable to save the role
				var role Role

				if err := rows.Scan(&role.Name); err != nil {
					log.Println(err)
					return
				}
				
				//a new role is added to the array
				roles = append(roles, role)
			}
			if err := rows.Err(); err != nil {
				log.Println(err)
			}
		}

		//converts array into a JSON and sends it to requestor
		render.JSON(res, req, roles)
	})

	router.Get("/listLocks", func(res http.ResponseWriter, req *http.Request) {
		//array that will save each door/lock from the database
		var doors []Door

		//variable will save the querry command for locks
		rows, err := db.Query(`select Doors.name from Doors`)

		//if statement makes sure that the query was a success; if successful then each row in the Doors scheme is read
		if err != nil{
			log.Println(err)
		}else{
			defer rows.Close()
			for rows.Next(){
				//variable to save the door/lock
				var door Door

				if err := rows.Scan(&door.Name); err != nil{
					log.Println(err)
					return
				}

				//a new door/lock is added to the array
				doors = append(doors, door)
			}
			if err := rows.Err(); err != nil{
				log.Println(err)
			}
		}

		//converts array into a JSON and sends it to the requester
		render.JSON(res, req, doors)
	})

	http.ListenAndServe(":8080", router)
}

