package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/OttoWBitt/NEWT/db"
	"github.com/OttoWBitt/NEWT/jwt"
	"github.com/mitchellh/mapstructure"
)

type UserInfo struct {
	Id       int    `json:"id"`
	UserName string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type Subject struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Artifact struct {
	Id           int      `json:"id"`
	Name         string   `json:"name"`
	User         UserInfo `json:"user"`
	Description  string   `json:"description"`
	Subject      Subject  `json:"subject"`
	Link         *string  `json:"link"`
	DonwloadLink *string  `json:"downloadLink"`
}

type Comment struct {
	Id       int      `json:"id"`
	Date     string   `json:"date"`
	User     UserInfo `json:"user"`
	Artifact Artifact `json:"artifact"`
	Comment  string   `json:"comment"`
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

func RenderResponse(res http.ResponseWriter, message *string, statusCode int) {

	if message == nil && statusCode == http.StatusOK {
		res.WriteHeader(statusCode)
		res.Write([]byte("sucess"))
		return
	}
	res.WriteHeader(statusCode)
	fmt.Println(statusCode)
	fmt.Println(*message)
	generateJSON := map[string]interface{}{
		"data":   nil,
		"errors": message,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}

func DecodeJwt(token string) (*jwt.JwtClaim, error) {
	jwtWrapper := jwt.JwtWrapper{
		SecretKey: jwt.SecretKey,
		Issuer:    "AuthService",
	}

	claims, err := jwtWrapper.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func ValidateAndReturnLoggedUser(jwtUser UserInfo) (*UserInfo, int, error) {

	user, err := GetUserByID(jwtUser.Id)
	if err != nil {
		erro := fmt.Sprintf("erro: %s", err)
		return nil, http.StatusInternalServerError, errors.New(erro)
	}

	if user.UserName != jwtUser.UserName || user.Email != jwtUser.Email || user.Id != jwtUser.Id {
		erro := fmt.Sprintf("InvalidLoggedUser - %s", err)
		return nil, http.StatusForbidden, errors.New(erro)
	}

	return user, http.StatusOK, nil
}

func GetUserFromContext(ctx context.Context) (*UserInfo, error) {

	teste := ctx.Value("user")
	var user UserInfo

	err := mapstructure.Decode(teste, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetSubjectByID(id int) (*Subject, error) {

	var resId, name string

	err := db.DB.QueryRow(
		`SELECT 
			id ,
			subject
		FROM 
			newt.subjects 
		WHERE 
			id = ?
	`, id).Scan(&resId, &name)

	if err != nil {
		return nil, err
	}

	intId, err := strconv.Atoi(resId)
	if err != nil {
		return nil, err
	}

	return &Subject{
		Name: name,
		Id:   intId,
	}, nil
}

func ConvertUtcToCustomTime(timestamp time.Time) string {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}

	now := timestamp.In(loc)

	return fmt.Sprintf("%d:%d:%d %d/%d/%d", now.Hour(), now.Minute(), now.Second(), now.Day(), now.Month(), now.Year())
}

func ConvertUtcToCustomTimeDate(timestamp time.Time) string {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		panic(err)
	}

	now := timestamp.In(loc)

	return fmt.Sprintf("%d/%d/%d", now.Day(), now.Month(), now.Year())
}
