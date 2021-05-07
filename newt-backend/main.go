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

	mux.HandleFunc("/api/signup", signupPage)
	mux.HandleFunc("/api/login", login)
	mux.HandleFunc("/api/recover", recoverPassword)
	mux.HandleFunc("/api/reset", resetPassword)

	mux.HandleFunc("/api/teste", testeBarros)

	mux.HandleFunc("/link", insertLinks)
	mux.HandleFunc("/link/recover", retrieveLinks)

	mux.HandleFunc("/fetch/files", fetchFiles)
	mux.HandleFunc("/fetch/files/id", fetchFilesByID)
	mux.HandleFunc("/fetch/links", fetchLinks)
	mux.HandleFunc("/fetch/links/id", fetchLinkByID)
	mux.HandleFunc("/fetch/getAllUser", fetchAllByUser)

	mux.HandleFunc("/", homePage)

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*", "Access-Control-Allow-Origin"},
	})

	//default functions
	handler := corsWrapper.Handler(mux)
	http.ListenAndServe(":3001", handler)

	defer dbMaster.Close()

}

func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "html/index.html")
}
