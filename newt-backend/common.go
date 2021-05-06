package main

import (
	"net/http"

	"github.com/OttoWBitt/NEWT/jwt"
)

func userExists(id int) error {

	var retID int

	err := dbMaster.QueryRow(
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

func getUserIDByEmail(email string) (string, error) {

	var id string

	err := dbMaster.QueryRow(
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

func renderError(res http.ResponseWriter, message string, statusCode int) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(message))
}

func decodeJwt(token string) (*jwt.JwtClaim, error) {
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
