package main

import (
	"net/http"
	"os"

	"github.com/OttoWBitt/NEWT/artifact"
	"github.com/OttoWBitt/NEWT/db"
	"github.com/OttoWBitt/NEWT/fetch"
	"github.com/OttoWBitt/NEWT/fileOps"
	"github.com/OttoWBitt/NEWT/login"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	db.OpenDataBase(os.Getenv("MYSQL"))
	router := mux.NewRouter()

	router.HandleFunc("/api/signup", login.Signup).Methods("POST")
	router.HandleFunc("/api/login", login.Login).Methods("POST")
	router.HandleFunc("/api/recover", login.RecoverPassword).Methods("POST")
	router.HandleFunc("/api/reset", login.ResetPassword).Methods("POST")

	router.HandleFunc("/api/teste", testeBarros).Methods("GET")

	router.HandleFunc("/api/artifact/new", artifact.InsertArtifacts).Methods("POST")

	fs := http.FileServer(http.Dir(fileOps.UploadPath))
	router.PathPrefix("/api/files/").Handler(http.StripPrefix("/api/files", fs))

	router.HandleFunc("/link", insertLinks)
	router.HandleFunc("/link/recover", retrieveLinks)

	router.HandleFunc("/download", fileOps.Download)

	router.HandleFunc("/fetch/files", fetch.FetchFiles)
	router.HandleFunc("/fetch/files/id", fetch.FetchFilesByID)
	router.HandleFunc("/fetch/links", fetch.FetchLinks)
	router.HandleFunc("/fetch/links/id", fetch.FetchLinkByID)
	router.HandleFunc("/fetch/getAllUser", fetch.FetchAllByUser)

	router.HandleFunc("/", homePage)

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*", "Access-Control-Allow-Origin"},
	})

	//default functions
	handler := corsWrapper.Handler(router)
	http.ListenAndServe(":3001", handler)

	defer db.DB.Close()

}

func homePage(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, "html/index.html")
}
