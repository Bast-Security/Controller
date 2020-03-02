package main

import (
	"log"
)

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

