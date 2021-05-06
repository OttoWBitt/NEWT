package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/OttoWBitt/NEWT/jwt"
	"golang.org/x/crypto/bcrypt"
)

type userInfo struct {
	id       int
	UserName string
	Name     string
	email    string
}

type inputData struct {
	UserName     string `json:"userName"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	RecoverToken string `json:"recoverToken"`
}

func signupPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {

		//data json
		data, erro := ioutil.ReadAll(req.Body)
		if erro != nil {
			renderError(res, erro.Error(), http.StatusInternalServerError)
		}

		var info inputData

		if erro = json.Unmarshal(data, &info); erro != nil {
			renderError(res, erro.Error(), http.StatusInternalServerError)
		}

		userName := info.UserName
		password := info.Password
		name := info.Name
		email := info.Email

		var user string

		err := dbMaster.QueryRow(`
		SELECT 
			username 
		FROM 
			newt.users 
		WHERE 
			username = ?
	`, userName).Scan(&user)

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
					email
				) 
			VALUES (?,?,?,?,?)
		`, userName, hashedPassword, name, email)

			if err != nil {
				http.Error(res, "Server error, unable to create your account.", http.StatusInternalServerError)
				return
			}

			res.Write([]byte("User created!"))
			return
		case err != nil:
			http.Error(res, "Server error, unable to create your account.", http.StatusInternalServerError)
			return
		default:
			renderError(res, erro.Error(), http.StatusMovedPermanently)
		}
	}
}

func login(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {

		//data json
		data, erro := ioutil.ReadAll(req.Body)
		if erro != nil {
			renderError(res, erro.Error(), http.StatusInternalServerError)
		}

		var info inputData

		if erro = json.Unmarshal(data, &info); erro != nil {
			renderError(res, erro.Error(), http.StatusInternalServerError)
		}

		userName := info.UserName
		password := info.Password

		var databasePassword string

		err := dbMaster.QueryRow(`
		SELECT
			password
		FROM 
			newt.users 
		WHERE username = ?
		`, userName).Scan(&databasePassword)

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
			id,
			username, 
			name,
			email
		FROM 
			newt.users 
		WHERE username = "%s"
	`

		userQuery := fmt.Sprintf(query, userName)

		var user userInfo

		err = dbMaster.QueryRow(userQuery).Scan(&user.id, &user.UserName, &user.Name, &user.email)
		if err != nil {
			renderError(res, err.Error(), http.StatusBadRequest)
			return
		}

		jwtWrapper := jwt.JwtWrapper{
			SecretKey:       "ChaveSecretaDoNEWTas65d@#$@#423jkl2j3423@#$2354ds5f4sd5f4sdf())@!sd6f5s6d4f54234",
			Issuer:          "AuthService",
			ExpirationHours: 1,
		}

		signedToken, err := jwtWrapper.GenerateToken(user.email, user.UserName, user.Name, user.id)
		if err != nil {
			renderError(res, err.Error(), http.StatusBadRequest)
			return
		}

		generateJSON := map[string]interface{}{
			"userName": user.UserName,
			"name":     user.Name,
			"token":    signedToken,
		}

		jsonData, err := json.Marshal(generateJSON)

		if err != nil {
			renderError(res, err.Error(), http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(jsonData)
	}
}

func recoverPassword(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {

		//data json
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			renderError(res, err.Error(), http.StatusInternalServerError)
		}

		var info inputData

		if err := json.Unmarshal(data, &info); err != nil {
			renderError(res, err.Error(), http.StatusInternalServerError)
		}

		email := info.Email

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
}

func resetPassword(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {

		//data json
		data, err := ioutil.ReadAll(req.Body)
		if err != nil {
			renderError(res, err.Error(), http.StatusInternalServerError)
		}

		var info inputData

		if err := json.Unmarshal(data, &info); err != nil {
			renderError(res, err.Error(), http.StatusInternalServerError)
		}

		recoverToken := info.RecoverToken
		password := info.Password

		if recoverToken == "" {
			renderError(res, "No recovery code", http.StatusInternalServerError)
		}

		decID, err := decodeString(recoverToken)
		if err != nil {
			renderError(res, "Could not decrypt", http.StatusInternalServerError)
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			renderError(res, "Server error, unable to create your account.", http.StatusInternalServerError)
			return
		}

		_, err = dbMaster.Exec(`
		UPDATE
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
}
