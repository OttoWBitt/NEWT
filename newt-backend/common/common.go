package common

import (
	"net/http"

	"github.com/OttoWBitt/NEWT/db"
	"github.com/OttoWBitt/NEWT/jwt"
)

func UserExists(id int) error {

	var retID int

	err := db.DB.QueryRow(
		`SELECT 
			id 
		FROM 
			newt.users 
		WHERE 
			id = ?
	`, id).Scan(&retID)

	if err != nil {
		return err
	}
	return nil
}

func GetUserIDByEmail(email string) (string, error) {

	var id string

	err := db.DB.QueryRow(
		`SELECT 
			id 
		FROM 
			newt.users 
		WHERE 
			email = ?
	`, email).Scan(&id)

	if err != nil {
		return "", err
	}
	return id, nil
}

func RenderError(res http.ResponseWriter, message string, statusCode int) {
	res.WriteHeader(statusCode)
	res.Write([]byte(message))
}

func DecodeJwt(token string) (*jwt.JwtClaim, error) {
	jwtWrapper := jwt.JwtWrapper{
		SecretKey: "ChaveSecretaDoNEWTas65d@#$@#423jkl2j3423@#$2354ds5f4sd5f4sdf())@!sd6f5s6d4f54234",
		Issuer:    "AuthService",
	}

	claims, err := jwtWrapper.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
