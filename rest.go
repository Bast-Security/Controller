package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-chi/jwtauth"
	jwt "github.com/dgrijalva/jwt-go"

	"math/big"
	"crypto/tls"
	"crypto/elliptic"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"net/http"
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

	router.Post("/register", handleRegister)
	router.Get("/login", getChallenge)
	router.Post("/login", handleLogin)

	router.Group(func(router chi.Router) {
		router.Use(jwtauth.Verifier(tokenAuth))
		router.Use(jwtauth.Authenticator)

		router.Post("/addSystem", addSystem)
		router.Post("/addUser", addUser)
		router.Post("/addRole", addRole)
		router.Post("/addLock", addLock)
		router.Get("/listUsers", listUsers)
		router.Get("/listRoles", listRoles)
		router.Get("/listLocks", listLocks)
	})

	return router
}

func addSystem(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(500)
}

func getChallenge(res http.ResponseWriter, req *http.Request) {
	challengeData := make([]byte, 16)
	if _, err := rand.Read(challengeData); err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}

	if _, err := db.Exec(`UPDATE Admins SET challenge = ? WHERE id = 1;`, challengeData); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	challenge := map[string][]byte{ "challenge": challengeData }
	render.JSON(res, req, challenge)
}

func handleLogin(res http.ResponseWriter, req *http.Request) {
	response := struct{ Response []byte }{ }
	if err := render.DecodeJSON(req.Body, &response); err != nil {
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

	row := db.QueryRow(`SELECT challenge, keyX, keyY FROM Admins WHERE id = 1;`)

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
	if _, err := asn1.Unmarshal(response.Response, sig); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	sha := sha256.New()
	sha.Write(challenge)
	hashed := sha.Sum(nil)

	if ecdsa.Verify(&pubKey, hashed, sig.R, sig.S) {
		if _, token, err := tokenAuth.Encode(jwt.MapClaims{ "uid": 1 }); err != nil {
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
	var pubKey ecdsa.PublicKey

	if err := render.DecodeJSON(req.Body, &pubKey); err != nil {
		log.Println(err)
		res.WriteHeader(400)
		return
	}

	_, err := db.Exec(`INSERT INTO Admins (keyX, keyY) VALUES (?, ?);`, pubKey.X.String(), pubKey.Y.String());
	if err != nil {
		log.Println(err)
		res.WriteHeader(500)
		return
	}

	res.WriteHeader(200)
}

func addLock(res http.ResponseWriter, req *http.Request) {
	var door string

	render.DecodeJSON(req.Body, &door)

	_, err := db.Exec(`INSERT INTO Door (name) VALUES (?);`, door)

	if err != nil{
		log.Println(err)
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
		log.Println(err)
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
		log.Println(err)
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

	//variable will save the querry command
	rows, err := db.Query(`select Roles.name from Roles`)
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

	//variable will save the querry command
	rows, err := db.Query(`select Users.id, Users.name, Users.email, Users.pin, Users.cardno from Users`)
	defer rows.Close()

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

