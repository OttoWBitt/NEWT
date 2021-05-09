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
	data, erro := ioutil.ReadAll(req.Body)
	if erro != nil {
		common.RenderResponse(res, erro.Error(), http.StatusInternalServerError)
	}

	var info inputData

	if erro = json.Unmarshal(data, &info); erro != nil {
		common.RenderResponse(res, erro.Error(), http.StatusInternalServerError)
	}

	userName := info.UserName
	password := info.Password
	name := info.Name
	email := info.Email

	var user string

	err := db.DB.QueryRow(`
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
			common.RenderResponse(res, "Server error, unable to create your account.", http.StatusInternalServerError)
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
			common.RenderResponse(res, "Server error, unable to create your account.", http.StatusInternalServerError)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		common.RenderResponse(res, "Server error, unable to create your account.", http.StatusInternalServerError)
		return
	default:
		common.RenderResponse(res, erro.Error(), http.StatusMovedPermanently)
	}

}

func Login(res http.ResponseWriter, req *http.Request) {

	//data json
	data, erro := ioutil.ReadAll(req.Body)
	if erro != nil {
		common.RenderResponse(res, erro.Error(), http.StatusInternalServerError)
	}

	var info inputData

	if erro = json.Unmarshal(data, &info); erro != nil {
		common.RenderResponse(res, erro.Error(), http.StatusInternalServerError)
	}

	userName := info.UserName
	password := info.Password

	var databasePassword string

	err := db.DB.QueryRow(`
		SELECT
			password
		FROM 
			newt.users 
		WHERE username = ?
		`, userName).Scan(&databasePassword)

	if err != nil {
		common.RenderResponse(res, err.Error(), http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		common.RenderResponse(res, err.Error(), http.StatusInternalServerError)
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
		common.RenderResponse(res, err.Error(), http.StatusBadRequest)
		return
	}

	jwtWrapper := jwt.JwtWrapper{
		SecretKey:       "ChaveSecretaDoNEWTas65d@#$@#423jkl2j3423@#$2354ds5f4sd5f4sdf())@!sd6f5s6d4f54234",
		Issuer:          "AuthService",
		ExpirationHours: 1,
	}

	signedToken, err := jwtWrapper.GenerateToken(user.Email, user.UserName, user.Name, user.Id)
	if err != nil {
		common.RenderResponse(res, err.Error(), http.StatusBadRequest)
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
		common.RenderResponse(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)

}

func RecoverPassword(res http.ResponseWriter, req *http.Request) {

	//data json
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		common.RenderResponse(res, err.Error(), http.StatusInternalServerError)
	}

	var info inputData

	if err := json.Unmarshal(data, &info); err != nil {
		common.RenderResponse(res, err.Error(), http.StatusInternalServerError)
	}

	email := info.Email

	id, err := common.GetUserIDByEmail(email)
	if err != nil {
		return
	}

	encID := crypto.EncodeToString(id)

	erro := mail.SendEmail(encID, email)
	if erro != nil {
		common.RenderResponse(res, "Error sending email", http.StatusInternalServerError)
	}
	res.Write([]byte("SUCCESS"))

}

func ResetPassword(res http.ResponseWriter, req *http.Request) {

	//data json
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		common.RenderResponse(res, err.Error(), http.StatusInternalServerError)
	}

	var info inputData

	if err := json.Unmarshal(data, &info); err != nil {
		common.RenderResponse(res, err.Error(), http.StatusInternalServerError)
	}

	recoverToken := info.RecoverToken
	password := info.Password

	if recoverToken == "" {
		common.RenderResponse(res, "No recovery code", http.StatusInternalServerError)
	}

	decID, err := crypto.DecodeString(recoverToken)
	if err != nil {
		common.RenderResponse(res, "Could not decrypt", http.StatusInternalServerError)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		common.RenderResponse(res, "Server error, unable to create your account.", http.StatusInternalServerError)
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
		common.RenderResponse(res, "Server error, unable to update your password.", http.StatusInternalServerError)
		return
	}

	res.Write([]byte("Password Updated!"))

}
