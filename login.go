package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func recoverPassword(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "html/recover.html")
		return
	}

	email := req.FormValue("email")

	id, err := getUserIDByEmail(email)
	if err != nil {
		return
	}

	encID := encodeToString(id)

	erro := sendEmail(encID, email)
	if erro != nil {
		renderError(res, "Error sending email", http.StatusInternalServerError)
	}
	res.Write([]byte("SUCCESS"))
}

func resetPassword(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "html/reset.html")
		return
	}

	id := req.FormValue("id")
	if id == "" {
		renderError(res, "No recovery code", http.StatusInternalServerError)
	}

	decID, err := decodeString(id)
	if err != nil {
		renderError(res, "Could not decrypt", http.StatusInternalServerError)
	}

	password := req.FormValue("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		renderError(res, "Server error, unable to create your account.", http.StatusInternalServerError)
		return
	}

	_, err = dbMaster.Exec(
		`UPDATE
			newt.users
		SET 
			password = ?
		WHERE
			id = ?
	`, hashedPassword, decID)

	if err != nil {
		http.Error(res, "Server error, unable to update your password.", http.StatusInternalServerError)
		return
	}

	res.Write([]byte("Password Updated!"))
}

type inputData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func jsonLoginPage(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {

		const theCase = "original"

		//data json
		data, erro := ioutil.ReadAll(req.Body)
		if erro != nil {
			renderError(res, erro.Error(), http.StatusInternalServerError)
		}

		var info inputData

		if erro = json.Unmarshal(data, &info); erro != nil {
			renderError(res, erro.Error(), http.StatusInternalServerError)
		}

		email := info.Email
		password := info.Password

		var databasePassword string

		err := dbMaster.QueryRow(
			`SELECT
				password
			FROM 
				newt.users 
			WHERE email = ?
		`, email).Scan(&databasePassword)

		if err != nil {
			renderError(res, err.Error(), http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
		if err != nil {
			renderError(res, err.Error(), http.StatusInternalServerError)
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
		WHERE email = "%s"
`

		userQuery := fmt.Sprintf(query, email)
		user, err := gosqljson.QueryDbToMap(dbMaster, theCase, userQuery)
		if err != nil {
			renderError(res, err.Error(), http.StatusInternalServerError)
		}

		generateJSON := map[string]interface{}{
			"user": user,
		}

		jsonData, err := json.Marshal(generateJSON)

		if err != nil {
			renderError(res, err.Error(), http.StatusInternalServerError)
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(jsonData)
	}
}
