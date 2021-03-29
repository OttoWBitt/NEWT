package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
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

	http.HandleFunc("/upload", uploadFileHandler()) //TODO: verificar user do upload, permissoes, salvar nome na base de dados pra fazer dwld dps
	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/files/", http.StripPrefix("/files", fs))

	http.HandleFunc("/", homePage)
	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/link", insertLinks)
	http.HandleFunc("/link/recover", retrieveLinks)
	http.ListenAndServe(":3000", nil)

	defer dbMaster.Close()

}

func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "html/index.html")
}
