package main

import (
	"github.com/go-chi/render"

	"io/ioutil"
	"net/http"
	"log"
)

func newAdmin(res http.ResponseWriter, req *http.Request) {
	key, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(500)
		return
	}

	result, err := db.Exec(`INSERT INTO Admins (pubKey) VALUES (?);`, key);
	if err != nil {
		res.WriteHeader(500)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		res.WriteHeader(500)
		return
	}

	render.JSON(res, req, map[string]int64{ "id": id })
}

func addLock(res http.ResponseWriter, req *http.Request) {
	var door string

	render.DecodeJSON(req.Body, &door)

	_, err := db.Exec(`INSERT INTO Door (name) VALUES (?);`, door)

	if err != nil{
		res.WriteHeader(400)
	}else{
		res.WriteHeader(200)
	}
}

func addRole(res http.ResponseWriter, req *http.Request) {
	var role string

	render.DecodeJSON(req.Body, &role)

	_, err := db.Exec(`INSERT INTO Roles (name) VALUES (?);`, role)

	if err != nil{
		res.WriteHeader(400)
	}else{
		res.WriteHeader(200)
	}
}

func addUser(res http.ResponseWriter, req *http.Request) {
	var user User

	render.DecodeJSON(req.Body, &user)

	_, err := db.Exec(`INSERT INTO Users (name, email, pin, cardno) VALUES (?, ?, ?, ?);`, user.Name, user.Email, user.Pin, user.Card)

	if err != nil {
		res.WriteHeader(400)
	} else {
		res.WriteHeader(200)
	}
}

func listLocks(res http.ResponseWriter, req *http.Request) {
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
}

func listRoles(res http.ResponseWriter, req *http.Request) {
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

}

func listUsers(res http.ResponseWriter, req *http.Request) {
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
}

