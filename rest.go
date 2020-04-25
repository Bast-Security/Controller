package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-chi/jwtauth"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fabiocolacio/hotp"

	"strconv"
	"context"
	"math/big"
	"crypto/tls"
	"crypto/elliptic"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"database/sql"
	"io/ioutil"
	"net/http"
	"time"
	"fmt"
	"log"
)

var (
	tokenAuth *jwtauth.JWTAuth
	signKey []byte
)


func init() {
	signKey = make([]byte, 16)
	if _, err := rand.Read(signKey); err != nil {
		log.Fatal("Unable to generate JWT signing key.")
	}
	tokenAuth = jwtauth.New("HS256", signKey, nil)
}

func server() http.Server {
	return http.Server{
		Addr: addr,
		Handler: router(),
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			},
		},
	}
}

func router() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	unimp := func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(500)
		fmt.Fprintln(res, "This route has not been implemented yet")
	}

	router.Post("/register", handleRegister)
	router.Post("/challenge", getChallenge)
	router.Post("/login", handleLogin)

	// Must be logged in as admin to access these routes
	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(jwtauth.Authenticator)

		router.Route("/systems", func(router chi.Router) {
			router.Post("/", addSystem)
			router.Get("/", listSystems)

			router.Route("/{systemId}", func(router chi.Router) {
				router.Use(systemContext)

				router.Get("/totp", totp)

				router.Route("/users", func(router chi.Router) {
					router.Post("/", addUser)
					router.Get("/", listUsers)

					router.Route("/{userId}", func(router chi.Router) {
						router.Get("/", getUser)
						router.Put("/", editUser)
						router.Delete("/", delUser)
						router.Get("/log", unimp)
					})
				})

				router.Route("/locks", func(router chi.Router) {
					router.Get("/", listLocks)
				})

				router.Route("/roles", func(router chi.Router) {
					router.Post("/", addRole)
					router.Get("/", listRoles)

					router.Route("/{role}", func(router chi.Router) {
						router.Get("/", getRole)
						router.Put("/", editRole)
						router.Delete("/", delRole)
						router.Get("/log", unimp)
					})
				})

				router.Get("/log", unimp)
			})
		})
	})

	router.Route("/locks", func(router chi.Router) {
		router.Post("/register", addLock)

		router.Route("/{lockId}", func(router chi.Router) {
			router.Get("/", lockChallenge)
			router.Post("/login", lockLogin)
			router.Post("/access", accessRequest)
		})
	})

	return router
}

func systemContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var (
			userId int64
			systemId int64
			err error
		)

		if userId, err = getId(req); err == nil {
			if systemId, err = strconv.ParseInt(chi.URLParam(req, "systemId"), 10, 64); err == nil {
				if !hasAccess(userId, systemId) {
					err = fmt.Errorf("No association between user %d and system %d.\n", userId, systemId)
				}
			}
		}

		if err != nil {
			log.Println(err)
			res.WriteHeader(400)
		} else {
			ctx := context.WithValue(context.WithValue(req.Context(), "adminId", userId), "systemId", systemId)
			next.ServeHTTP(res, req.WithContext(ctx))
		}
	})
}

func getId(req *http.Request) (id int64, err error) {
	var claims jwt.MapClaims
	if _, claims, err = jwtauth.FromContext(req.Context()); err == nil {
		var (
			i interface{}
			ok bool
		)

		if i, ok = claims["id"]; !ok {
			err = fmt.Errorf("No id in claims")
		} else {
			id = int64(i.(float64))
		}
	}
	return
}

func hasAccess(adminId, systemId int64) bool {
	var id int64
	row := db.QueryRow(`SELECT admin FROM AdminSystem WHERE admin=? AND system=?;`, adminId, systemId)
	if err := row.Scan(&id); err != nil {
		log.Println("Admin doesn't have access to the system. ", err)
		return false
	}
	if id == adminId {
		return true
	}
	return false
}

func addSystem(res http.ResponseWriter, req *http.Request) {
	var (
		uid int64
		system System
		result sql.Result
		err error
	)

	if uid, err = getId(req); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	if err := render.DecodeJSON(req.Body, &system); err != nil && system.Name != "" {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	system.TotpKey = make([]byte, 32)
	if _, err := rand.Read(system.TotpKey); err != nil {
		log.Println("Failed to create TotpKey ", err)
		res.WriteHeader(500)
		return
	}

	if result, err = db.Exec(`INSERT INTO Systems (name, totpKey) VALUES (?, ?);`, system.Name, system.TotpKey); err == nil {
		if system.Id, err = result.LastInsertId(); err == nil {
			_, err = db.Exec(`INSERT INTO AdminSystem (admin, system) VALUES (?, ?);`, uid, system.Id)
		}
	}
	if err != nil {
		log.Println(err)
		res.WriteHeader(500)
	}

	res.WriteHeader(200)
}

func getChallenge(res http.ResponseWriter, req *http.Request) {
	challengeData := make([]byte, 16)
	if _, err := rand.Read(challengeData); err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}

	var user Admin
	if err := render.DecodeJSON(req.Body, &user); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	if _, err := db.Exec(`UPDATE Admins SET challenge = ? WHERE id = ?;`, challengeData, user.Id); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	user.Challenge = challengeData
	render.JSON(res, req, user)
}

func handleLogin(res http.ResponseWriter, req *http.Request) {
	var user Admin
	if err := render.DecodeJSON(req.Body, &user); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	var (
		challenge []byte
		pubKey ecdsa.PublicKey = ecdsa.PublicKey{ Curve: elliptic.P384(), X: new(big.Int), Y: new(big.Int) }
		x []byte
		y []byte
	)

	row := db.QueryRow(`SELECT challenge, keyX, keyY FROM Admins WHERE id = ?;`, user.Id)

	if err := row.Scan(&challenge, &x, &y); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}
	if err := pubKey.X.UnmarshalText(x); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}
	if err := pubKey.Y.UnmarshalText(y); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	sig := &struct{ R, S *big.Int }{}
	if _, err := asn1.Unmarshal(user.Response, sig); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	sha := sha256.New()
	sha.Write(challenge)
	hashed := sha.Sum(nil)

	if ecdsa.Verify(&pubKey, hashed, sig.R, sig.S) {
		if _, token, err := tokenAuth.Encode(jwt.MapClaims{ "id": user.Id }); err != nil {
			log.Println(err)
			res.WriteHeader(500)
		} else {
			res.Write([]byte(token))
		}
	} else {
		log.Println("Verify failed")
		res.WriteHeader(400)
	}
}

func handleRegister(res http.ResponseWriter, req *http.Request) {
	var (
		pubKey ecdsa.PublicKey
		id int64
		err error
	)

	if err := render.DecodeJSON(req.Body, &pubKey); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	result, err := db.Exec(`INSERT INTO Admins (keyX, keyY) VALUES (?, ?);`, pubKey.X.String(), pubKey.Y.String());
	if err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}

	if id, err = result.LastInsertId(); err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}

	render.JSON(res, req, map[string]int64{ "id": id })
}

func addLock(res http.ResponseWriter, req *http.Request) {
	var door Door

	render.DecodeJSON(req.Body, &door)

	result, err := db.Exec(`INSERT INTO Doors (name, system, keyX, keyY) VALUES (?, ?, ?, ?);`, door.Name, door.System, door.KeyX.String(), door.KeyY.String())

	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
	}

	if result != nil {
		if id, err := result.LastInsertId(); err != nil {
			log.Println(err)
			res.WriteHeader(500)
		} else {
			render.JSON(res, req, map[string]int64{ "id": id })
		}
	} else {
		log.Println("Result is null!")
		res.WriteHeader(500)
	}
}

func accessRequest(res http.ResponseWriter, req *http.Request) {
	var creds struct{
			Card string `json:"card,omitempty"`
			Pin string `json:"pin,omitempty"`
		}
	
	if err := render.DecodeJSON(req.Body, &creds); err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}

	lockId := chi.URLParam(req, "lockId")

	var row *sql.Row

	if len(creds.Card) > 0 {
		row = db.QueryRow(`SELECT id FROM Users
				   INNER JOIN UserRole ON UserRole.userid = Users.id
				   INNER JOIN Permissions ON UserRole.role = Permissions.role
				   INNER JOIN Doors ON Permissions.system = Doors.system
				   WHERE Permissions.door = ?
				   AND cardno = ?`, lockId, creds.Card)
	} else if len(creds.Pin) > 0 {
		row = db.QueryRow(`SELECT id FROM Users
				   INNER JOIN UserRole ON UserRole.userid = Users.id
				   INNER JOIN Permissions ON UserRole.role = Permissions.role
				   INNER JOIN Doors ON Permissions.system = Doors.system
				   WHERE Permissions.door = ?
				   AND pin = ?`, lockId, creds.Pin)
	} else {
		res.WriteHeader(403)
		return
	}
	

	var userId string
	if err := row.Scan(&userId); err == sql.ErrNoRows {
		res.WriteHeader(403)
		return
	} else if err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	} else {
		log.Println("ACCESS GRANTED TO ", userId, " at Door ", lockId)
		res.WriteHeader(200)
	}
}

func lockChallenge(res http.ResponseWriter, req *http.Request) {
	challengeData := make([]byte, 16)
	if _, err := rand.Read(challengeData); err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}

	lockId := chi.URLParam(req, "lockId")

	if _, err := db.Exec(`UPDATE Doors SET challenge = ? WHERE id = ?;`, challengeData, lockId); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	var door Door
	door.Challenge = challengeData
	render.JSON(res, req, door)
}

func lockLogin(res http.ResponseWriter, req *http.Request) {
	var (
		lockId int64
		err error
	)

	lockId, err = strconv.ParseInt(chi.URLParam(req, "lockId"), 10, 64)

	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	var (
		challenge []byte
		response []byte
		pubKey ecdsa.PublicKey = ecdsa.PublicKey{ Curve: elliptic.P384(), X: new(big.Int), Y: new(big.Int) }
		x []byte
		y []byte
	)

	row := db.QueryRow(`SELECT challenge, keyX, keyY FROM Doors WHERE id = ?;`, lockId)

	if err := row.Scan(&challenge, &x, &y); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}
	if err := pubKey.X.UnmarshalText(x); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}
	if err := pubKey.Y.UnmarshalText(y); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	response, err = ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}

	sig := &struct{ R, S *big.Int }{}
	if _, err := asn1.Unmarshal(response, sig); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	sha := sha256.New()
	sha.Write(challenge)
	hashed := sha.Sum(nil)

	if ecdsa.Verify(&pubKey, hashed, sig.R, sig.S) {
		if _, token, err := tokenAuth.Encode(jwt.MapClaims{ "doorid": lockId }); err != nil {
			log.Println(err)
			res.WriteHeader(500)
		} else {
			res.Write([]byte(token))
		}
	} else {
		log.Println("Verify failed")
		res.WriteHeader(400)
	}
}

func totp(res http.ResponseWriter, req *http.Request) {
	var (
		totpKey  []byte
		now      int64
		duration int64
		epoch    int64
		digits   int
		err      error
	)

	systemId := req.Context().Value("systemId").(int64)

	row := db.QueryRow(`SELECT FROM Systems totpKey WHERE id = ?;`, systemId)

	if err = row.Scan(&totpKey); err != nil {
		log.Println(err)
		res.WriteHeader(400)
	}

	duration = 60 * 5 // 5 mins
	now = time.Now().Unix()
	epoch = 0
	digits = 6

	code := hotp.Totp(totpKey, now, epoch, duration, digits)

	render.JSON(res, req, map[string]int{ "code": code })
}

func addRole(res http.ResponseWriter, req *http.Request) {
	var (
		role Role
		err error
	)

	render.DecodeJSON(req.Body, &role)

	role.System = req.Context().Value("systemId").(int64)

	_, err = db.Exec(`INSERT INTO Roles (name, system) VALUES (?, ?);`, role.Name, role.System)

	if err != nil{
		log.Println(err)
		res.WriteHeader(400)
	} else {
		res.WriteHeader(200)
	}
}

func getRole(res http.ResponseWriter, req *http.Request) {
	var (
		role Role
		row *sql.Row
		rows *sql.Rows
		err error
	)

	system := req.Context().Value("systemId").(int64)
	name := chi.URLParam(req, "role")

	row = db.QueryRow(`SELECT id, FROM Roles WHERE system=? AND name=?;`, system, name)

	if err = row.Scan(&role.Id); err == nil {
		role.Name = name
		role.System = system

		if rows, err = db.Query(`SELECT door FROM Permissions WHERE system=? AND role=?;`, system, role.Id); err == nil {
			defer rows.Close()

			for rows.Next() {
				var door Door

				rows.Scan(&door.Id)

				row = db.QueryRow(`SELECT name FROM Doors WHERE id=?`, door.Id)
				row.Scan(&door.Name)

				role.Doors = append(role.Doors, door)
			}
		}
	}

	if err != nil {
		log.Println("Failed to fetch role: ", err)
		res.WriteHeader(404)
		return
	}

	render.JSON(res, req, role)
}

func editRole(res http.ResponseWriter, req *http.Request) {
	var role Role

	system := req.Context().Value("systemId").(int64)
	name := chi.URLParam(req, "role")

	render.DecodeJSON(req.Body, &role)

	if len(role.Name) > 0 {
		if _, err := db.Exec(`UPDATE Roles SET name=? WHERE name=? AND system=?;`, role.Name, name, system); err != nil {
			log.Println("Failed to update role name: ", err)
			res.WriteHeader(500)
			return
		}
	}

	if _, err := db.Exec(`DELETE FROM Permissions WHERE role=?;`, role.Id); err != nil {
		log.Println("Failed to remove old permission set: ", err)
		res.WriteHeader(500)
		return
	}

	for _, door := range role.Doors {
		if _, err := db.Exec(`INSERT INTO Permissions (system, role, door) VALUES (?, ?, ?);`, system, role.Id, door.Id); err != nil {
			log.Println("Failed to add permission: ", err)
			res.WriteHeader(500)
			return
		}
	}

	res.WriteHeader(200)
}

func delRole(res http.ResponseWriter, req *http.Request) {
	var err error

	system := req.Context().Value("systemId").(int64)
	name := chi.URLParam(req, "role")

	if _, err = db.Exec(`DELETE FROM UserRole WHERE role=? AND system=?;`, name, system); err == nil {
		if _, err = db.Exec(`DELETE FROM Permissions WHERE role=? AND system=?;`, name, system); err == nil {
			_, err = db.Exec(`DELETE FROM Roles WHERE name=? AND system=?;`, name, system)
		}
	}

	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	res.WriteHeader(200)
}

func addUser(res http.ResponseWriter, req *http.Request) {
	var user User

	system := req.Context().Value("systemId").(int64)
	render.DecodeJSON(req.Body, &user)

	_, err := db.Exec(`INSERT INTO Users (system, name, email, pin, cardno, phone) VALUES (?, ?, ?, ?, ?, ?);`, system, user.Name, user.Email, user.Pin, user.CardNo, user.Phone)

	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
	} else {
		res.WriteHeader(200)
	}
}

func getUser(res http.ResponseWriter, req *http.Request) {
	var (
		user User
		row *sql.Row
		rows *sql.Rows
		err error
	)

	system := req.Context().Value("systemId").(int64)
	userId := chi.URLParam(req, "userId")

	row = db.QueryRow(`SELECT name, email, phone, pin, cardno FROM Users WHERE id=? AND system=?`, userId, system)
	if err = row.Scan(&user.Name, &user.Email, &user.Phone, &user.Pin, &user.CardNo); err == nil {
		if rows, err = db.Query(`SELECT role FROM Roles WHERE system=? userid=?;`, system, userId); err == nil {
			defer rows.Close()

			for rows.Next() {
				var role Role

				rows.Scan(&role.Id)

				row = db.QueryRow(`SELECT name FROM Roles WHERE id=?;`, role.Id)

				if err = rows.Scan(&role.Name); err == nil {
					user.Roles = append(user.Roles, role)
				}
			}
		}
	}

	if err != nil {
		log.Println(err)
		res.WriteHeader(404)
	}

	render.JSON(res, req, user)
}

func editUser(res http.ResponseWriter, req *http.Request) {
	system := req.Context().Value("systemId").(int64)
	userId := chi.URLParam(req, "userId")

	var user User
	render.DecodeJSON(req.Body, &user)

	if len(user.Name) > 0 {
		if _, err := db.Exec(`UPDATE Users SET name=? WHERE id=?;`, user.Name, userId); err != nil {
			log.Println(err)
			res.WriteHeader(500)
			return
		}
	}

	if len(user.Email) > 0 {
		if _, err := db.Exec(`UPDATE Users SET email=? WHERE id=?;`, user.Email, userId); err != nil {
			log.Println(err)
			res.WriteHeader(500)
			return
		}
	}

	if len(user.Phone) > 0 {
		if _, err := db.Exec(`UPDATE Users SET phone=? WHERE id=?;`, user.Phone, userId); err != nil {
			log.Println(err)
			res.WriteHeader(500)
			return
		}
	}

	if len(user.Pin) > 0 {
		if _, err := db.Exec(`UPDATE Users SET pin=? WHERE id=?;`, user.Pin, userId); err != nil {
			log.Println(err)
			res.WriteHeader(500)
			return
		}
	}

	if len(user.CardNo) > 0 {
		if _, err := db.Exec(`UPDATE Users SET cardno=? WHERE id=?;`, user.CardNo, userId); err != nil {
			log.Println(err)
			res.WriteHeader(500)
			return
		}
	}

	if _, err := db.Exec(`DELETE FROM Roles WHERE userid=?;`, userId); err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}

	for _, role := range user.Roles {
		if _, err := db.Exec(`INSERT INTO UserRole (system, userid, role) VALUES (?, ?, ?);`, system, userId, role.Id); err != nil {
			log.Println(err)
			res.WriteHeader(500)
			return
		}
	}

	res.WriteHeader(200)
}

func delUser(res http.ResponseWriter, req *http.Request) {
	var err error

	system := req.Context().Value("systemId").(int64)
	userId := chi.URLParam(req, "userId")

	if _, err = db.Exec(`DELETE FROM UserRole WHERE userid=? AND system=?;`, userId, system); err == nil {
		_, err = db.Exec(`DELETE FROM Users WHERE id=? AND system=?;`, userId, system)
	}

	if err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	res.WriteHeader(200)
}

func listLocks(res http.ResponseWriter, req *http.Request) {
	//array that will save each door/lock from the database
	var doors []Door

	system := req.Context().Value("systemId").(int64)

	//variable will save the querry command for locks
	rows, err := db.Query(`SELECT Doors.name FROM Doors WHERE system=?`, system)
	defer rows.Close()

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

	system := req.Context().Value("systemId").(int64)

	//variable will save the querry command
	rows, err := db.Query(`SELECT Roles.name FROM Roles WHERE system=?;`, system)
	defer rows.Close()

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

	system := req.Context().Value("systemId").(int64)

	rows, err := db.Query(`SELECT Users.id, Users.name, Users.email, Users.phone, Users.pin, Users.cardno FROM Users WHERE system=?;`, system)
	defer rows.Close()

	if err != nil {
		log.Println(err)
	} else {
		defer rows.Close()
		for rows.Next() {
			var user User

			if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.Phone, &user.Pin, &user.CardNo); err != nil {
				log.Println(err)
				return
			}

			users = append(users, user)
		}
		if err := rows.Err(); err != nil {
			log.Println(err)
		}
	}

	//converts array into a JSON and sends it to requestor
	render.JSON(res, req, users)
}

func listSystems(res http.ResponseWriter, req *http.Request) {
	var (
		systems []System
		id int64
		rows *sql.Rows
		err error
	)

	if id, err = getId(req); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	if rows, err = db.Query(`SELECT id, name FROM Systems
		INNER JOIN AdminSystem ON Systems.id = AdminSystem.system
		WHERE AdminSystem.admin = ?;`, id); err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var system System
		rows.Scan(&system.Id, &system.Name)
		systems = append(systems, system)
	}
	if err = rows.Err(); err != nil {
		log.Println(err)
	}

	render.JSON(res, req, systems)
}

