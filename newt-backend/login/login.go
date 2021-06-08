package login

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/OttoWBitt/NEWT/common"
	"github.com/OttoWBitt/NEWT/crypto"
	"github.com/OttoWBitt/NEWT/db"
	mail "github.com/OttoWBitt/NEWT/email"
	"github.com/gorilla/mux"

	"github.com/OttoWBitt/NEWT/jwt"
	"golang.org/x/crypto/bcrypt"
)

type inputData struct {
	UserName     string `json:"username"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	RecoverToken string `json:"token"`
}

func Signup(res http.ResponseWriter, req *http.Request) {

	//data json
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	var info inputData

	if err = json.Unmarshal(data, &info); err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	userName := info.UserName
	password := info.Password
	name := info.Name
	email := info.Email

	var recUsername, recEmail string

	err = db.DB.QueryRow(`
		SELECT 
			username,
			email 
		FROM 
			newt.users 
		WHERE 
			username = ? OR 
			email = ?
	`, userName, email).Scan(&recUsername, &recEmail)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			erro := "Server error, unable to create your account. - " + err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		row, err := db.DB.Exec(
			`INSERT INTO 
				newt.users(
					username, 
					password,
					name,
					email
				) 
			VALUES (?,?,?,?)
		`, userName, hashedPassword, name, email)

		if err != nil {
			erro := "Server error, unable to create your account. - " + err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		lastId, _ := row.LastInsertId()

		retUser := common.UserInfo{
			Id:       int(lastId),
			UserName: userName,
			Name:     name,
			Email:    email,
		}

		jwtWrapper := jwt.JwtWrapper{
			SecretKey:       jwt.SecretKey,
			Issuer:          "AuthService",
			ExpirationHours: 1,
		}

		signedToken, err := jwtWrapper.GenerateToken(retUser.Email, retUser.UserName, retUser.Name, retUser.Id)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusBadRequest)
			return
		}

		generateJSON := map[string]interface{}{
			"user":  retUser,
			"token": signedToken,
		}

		jsonData, err := json.Marshal(generateJSON)

		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(jsonData)
		return
	case err != nil:
		erro := "Server error, unable to create your account. - " + err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	default:

		fmt.Println(recUsername)
		fmt.Println(recEmail)

		erro := "Error creating user! - "

		if recEmail == email && recUsername == userName {
			erro = erro + "username and email not available"
		} else if recEmail == email {
			erro = erro + "email not available"
		} else if recUsername == userName {
			erro = erro + "username not available"
		}

		common.RenderResponse(res, &erro, http.StatusInternalServerError)
	}

}

func Login(res http.ResponseWriter, req *http.Request) {

	//data json
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
	}

	var info inputData

	if err = json.Unmarshal(data, &info); err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
	}

	userName := info.UserName
	password := info.Password

	var databasePassword string

	err = db.DB.QueryRow(`
		SELECT
			password
		FROM 
			newt.users 
		WHERE username = ?
		`, userName).Scan(&databasePassword)

	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
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

	var user common.UserInfo

	err = db.DB.QueryRow(userQuery).Scan(&user.Id, &user.UserName, &user.Name, &user.Email)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	jwtWrapper := jwt.JwtWrapper{
		SecretKey:       jwt.SecretKey,
		Issuer:          "AuthService",
		ExpirationHours: 1,
	}

	signedToken, err := jwtWrapper.GenerateToken(user.Email, user.UserName, user.Name, user.Id)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	generateJSON := map[string]interface{}{
		"user":  user,
		"token": signedToken,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)

}

func RecoverPassword(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	email, ok := vars["email"]
	if !ok {
		erro := "email is missing in parameters"
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	id, err := common.GetUserIDByEmail(email)
	if err != nil {
		return
	}

	encID := crypto.EncodeToString(id)

	erro := mail.SendEmail(encID, email)
	if erro != nil {
		erro := "Error sending email"
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
	}

	generateJSON := map[string]interface{}{
		"data":   email,
		"errors": nil,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)

}

func ResetPassword(res http.ResponseWriter, req *http.Request) {

	//data json
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
	}

	var info inputData

	if err := json.Unmarshal(data, &info); err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
	}

	recoverToken := info.RecoverToken
	password := info.Password

	if recoverToken == "" {
		erro := "No recovery code"
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
	}

	decID, err := crypto.DecodeString(recoverToken)
	if err != nil {
		erro := "Could not decrypt"
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		erro := "Server error, unable to create your account."
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec(`
		UPDATE
			newt.users
		SET 
			password = ?
		WHERE
			id = ?
	`, hashedPassword, decID)

	if err != nil {
		erro := "Server error, unable to update your password."
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	userId, err := strconv.Atoi(decID)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	user, err := common.GetUserByID(userId)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	generateJSON := map[string]interface{}{
		"data":   user.Email,
		"errors": nil,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)

}

func SignupWithEmail(res http.ResponseWriter, req *http.Request) {

	//data json
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
	}

	var info inputData

	if err = json.Unmarshal(data, &info); err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
	}

	userName := info.UserName
	password := info.Password
	name := info.Name
	email := info.Email

	var ema string

	err = db.DB.QueryRow(`
		SELECT 
			email 
		FROM 
			newt.users 
		WHERE 
			email = ?
	`, userName).Scan(&ema)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			erro := "Server error, unable to create your account."
			fmt.Println("AQUI")
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		_, err = db.DB.Exec(
			`INSERT INTO 
				newt.users( 
					username,
					password,
					name,
					email
				) 
			VALUES (?,?,?,?)
		`, name, hashedPassword, name, email)

		if err != nil {
			erro := "Server error, unable to create your account."
			fmt.Println("AQUI2")
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		erro := "Server error, unable to create your account."
		fmt.Println("AQUI3")
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	default:
		erro := "User already exist."
		common.RenderResponse(res, &erro, http.StatusMovedPermanently)
	}

}

func LoginWithEmail(res http.ResponseWriter, req *http.Request) {

	//data json
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
	}

	var info inputData

	if err = json.Unmarshal(data, &info); err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
	}

	email := info.Email
	password := info.Password

	var databasePassword string

	err = db.DB.QueryRow(`
		SELECT
			password
		FROM 
			newt.users 
		WHERE email = ?
		`, email).Scan(&databasePassword)

	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
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
		WHERE email = "%s"
	`

	userQuery := fmt.Sprintf(query, email)

	var user common.UserInfo

	err = db.DB.QueryRow(userQuery).Scan(&user.Id, &user.UserName, &user.Name, &user.Email)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	jwtWrapper := jwt.JwtWrapper{
		SecretKey:       jwt.SecretKey,
		Issuer:          "AuthService",
		ExpirationHours: 1,
	}

	signedToken, err := jwtWrapper.GenerateToken(user.Email, user.UserName, user.Name, user.Id)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	generateJSON := map[string]interface{}{
		"id":    user.Id,
		"name":  user.Name,
		"token": signedToken,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)

}
