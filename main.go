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

	http.HandleFunc("/upload", uploadFileHandler())
	http.HandleFunc("/download", download)

	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/files/", http.StripPrefix("/files", fs))

	http.HandleFunc("/signup", signupPage)
	http.HandleFunc("/login", loginPage)

	http.HandleFunc("/link", insertLinks)
	http.HandleFunc("/link/recover", retrieveLinks)

	http.HandleFunc("/send", sendEmail)

	http.HandleFunc("/fetch/files", fetchFiles)
	http.HandleFunc("/fetch/files/id", fetchFilesByID)

	http.HandleFunc("/", homePage)

	http.ListenAndServe(":3000", nil)

	defer dbMaster.Close()

}

func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "html/index.html")
}
