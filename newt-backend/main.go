package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
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
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", uploadFileHandler())
	mux.HandleFunc("/download", download)

	fs := http.FileServer(http.Dir(uploadPath))
	mux.Handle("/files/", http.StripPrefix("/files", fs))

	mux.HandleFunc("/signup", signupPage)
	mux.HandleFunc("/login", loginPage)
	mux.HandleFunc("/recover", recoverPassword)
	mux.HandleFunc("/reset", resetPassword)

	mux.HandleFunc("/link", insertLinks)
	mux.HandleFunc("/link/recover", retrieveLinks)

	mux.HandleFunc("/fetch/files", fetchFiles)
	mux.HandleFunc("/fetch/files/id", fetchFilesByID)
	mux.HandleFunc("/fetch/links", fetchLinks)
	mux.HandleFunc("/fetch/links/id", fetchLinkByID)
	mux.HandleFunc("/fetch/getAllUser", fetchAllByUser)

	mux.HandleFunc("/", homePage)

	//kaka
	mux.HandleFunc("/api/login", jsonLoginPage)

	//default functions
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(":3001", handler)

	defer dbMaster.Close()

}

func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "html/index.html")
}
