package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/elgs/gosqljson"
	"golang.org/x/crypto/bcrypt"
)

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "html/singup.html")
		return
	}

	username := req.FormValue("username")
	password := req.FormValue("password")
	name := req.FormValue("name")
	email := req.FormValue("email")
	phone := req.FormValue("phone")

	var user string

	err := dbMaster.QueryRow(
		`SELECT 
			username 
		FROM 
			newt.users 
		WHERE 
			username = ?
	`, username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, "Server error, unable to create your account.", http.StatusInternalServerError)
			return
		}

		_, err = dbMaster.Exec(
			`INSERT INTO 
				newt.users(
					username, 
					password,
					name,
					email,
					phone
				) 
			VALUES (?,?,?,?,?)
		`, username, hashedPassword, name, email, phone)

		if err != nil {
			http.Error(res, "Server error, unable to create your account.", http.StatusInternalServerError)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		fmt.Println(err)
		http.Error(res, "Server error, unable to create your account.", http.StatusInternalServerError)
		return
	default:
		http.Redirect(res, req, "/", http.StatusMovedPermanently)
	}
}

func loginPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "html/login.html")
		return
	}

	const theCase = "original"

	username := req.FormValue("username")
	password := req.FormValue("password")

	var databasePassword string

	err := dbMaster.QueryRow(
		`SELECT
			password
		FROM 
			newt.users 
		WHERE username = ?
		`, username).Scan(&databasePassword)

	if err != nil {
		http.Redirect(res, req, "/login", http.StatusMovedPermanently)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/login", http.StatusMovedPermanently)
		return
	}

	const query = `
	SELECT
		username, 
		name,
		email,
		phone
	FROM 
		newt.users 
	WHERE username = "%s"
`

	userQuery := fmt.Sprintf(query, username)
	user, err := gosqljson.QueryDbToMap(dbMaster, theCase, userQuery)
	if err != nil {
		log.Fatal(err)
	}

	generateJSON := map[string]interface{}{
		"user": user,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}
