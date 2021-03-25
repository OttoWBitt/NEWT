package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/elgs/gosqljson"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

var dbMaster *sql.DB

//dataBase opens the data base connection to master
func dataBaseMaster(dsn string) {

	var err error
	dbMaster, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	dataBaseMaster(os.Getenv("MYSQL"))

	http.HandleFunc("/", homePage)
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/login", loginPage)
	http.ListenAndServe(":3000", nil)

	defer dbMaster.Close()

}

func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "html/index.html")
}

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
			http.Error(res, "Server error, unable to create your account.", 500)
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
			http.Error(res, "Server error, unable to create your account.", 500)
			return
		}

		res.Write([]byte("User created!"))
		return
	case err != nil:
		fmt.Println(err)
		http.Error(res, "Server error, unable to create your account.", 500)
		return
	default:
		http.Redirect(res, req, "/", 301)
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
		http.Redirect(res, req, "/login", 301)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
	if err != nil {
		http.Redirect(res, req, "/login", 301)
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
