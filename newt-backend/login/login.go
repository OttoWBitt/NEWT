package login

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/OttoWBitt/NEWT/common"
	"github.com/OttoWBitt/NEWT/crypto"
	"github.com/OttoWBitt/NEWT/db"
	mail "github.com/OttoWBitt/NEWT/email"

	"github.com/OttoWBitt/NEWT/jwt"
	"golang.org/x/crypto/bcrypt"
)

type inputData struct {
	UserName     string `json:"userName"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	RecoverToken string `json:"recoverToken"`
}

func Signup(res http.ResponseWriter, req *http.Request) {

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

	var user string

	err = db.DB.QueryRow(`
		SELECT 
			username 
		FROM 
			newt.users 
		WHERE 
			username = ? AND 
			email = ?
	`, userName).Scan(&user, &email)

	switch {
	case err == sql.ErrNoRows:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			erro := "Server error, unable to create your account."
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
			VALUES (?,?,?,?,?)
		`, userName, hashedPassword, name, email)

		if err != nil {
			erro := "Server error, unable to create your account."
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		erro := "Server error, unable to create your account."
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	default:
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusMovedPermanently)
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

	email := info.Email

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
	res.Write([]byte("SUCCESS"))

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

	res.Write([]byte("Password Updated!"))

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
		"id":       user.Id,
		"userName": user.UserName,
		"name":     user.Name,
		"token":    signedToken,
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
