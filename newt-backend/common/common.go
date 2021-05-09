package common

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/OttoWBitt/NEWT/db"
	"github.com/OttoWBitt/NEWT/jwt"
)

type UserInfo struct {
	Id       int
	UserName string
	Name     string
	Email    string
}

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

func GetUserByID(id int) (*UserInfo, error) {

	var resId, username, name, email string

	err := db.DB.QueryRow(
		`SELECT 
			id ,
			username,
			name,
			email
		FROM 
			newt.users 
		WHERE 
			id = ?
	`, id).Scan(&resId, &username, &name, &email)

	if err != nil {
		return nil, err
	}

	intId, err := strconv.Atoi(resId)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		UserName: username,
		Name:     name,
		Id:       intId,
		Email:    email,
	}, nil
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

func RenderResponse(res http.ResponseWriter, message string, statusCode int) {
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

func ValidateAndReturnLoggedUser(jwt string) (*UserInfo, int, error) {

	if len(jwt) == 0 {
		return nil, http.StatusForbidden, errors.New("UserNotLoggedIn")
	}

	jwtUser, err := DecodeJwt(jwt)
	if err != nil {
		erro := fmt.Sprintf("UserNotLoggedIn - %s", err)
		return nil, http.StatusForbidden, errors.New(erro)
	}

	user, err := GetUserByID(jwtUser.ID)
	if err != nil {
		erro := fmt.Sprintf("erro: %s", err)
		return nil, http.StatusInternalServerError, errors.New(erro)
	}

	if user.UserName != jwtUser.UserName {
		erro := fmt.Sprintf("InvalidLoggedUser - %s", err)
		return nil, http.StatusForbidden, errors.New(erro)
	}

	return user, http.StatusOK, nil
}
