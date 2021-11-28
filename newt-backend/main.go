package main

import (
	"net/http"
	"os"

	"github.com/OttoWBitt/NEWT/artifact"
	colourblind "github.com/OttoWBitt/NEWT/colourBlind"
	"github.com/OttoWBitt/NEWT/comments"
	"github.com/OttoWBitt/NEWT/db"
	"github.com/OttoWBitt/NEWT/fetch"
	"github.com/OttoWBitt/NEWT/fileOps"
	"github.com/OttoWBitt/NEWT/login"
	"github.com/OttoWBitt/NEWT/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	db.OpenDataBase(os.Getenv("MYSQL"))
	router := mux.NewRouter()

	// Rotas de login
	router.HandleFunc("/api/signup", login.Signup).Methods("POST")
	router.HandleFunc("/api/login", login.Login).Methods("POST")
	router.HandleFunc("/api/recover/{email}", login.RecoverPassword).Methods("POST")
	router.HandleFunc("/api/reset", login.ResetPassword).Methods("POST")

	// Rotas do caso de uso newt
	router.Handle("/api/artifact/new", middleware.AuthMiddleware(http.HandlerFunc(artifact.InsertArtifacts))).Methods("POST")
	router.Handle("/api/comment/new", middleware.AuthMiddleware(http.HandlerFunc(comments.InsertComments))).Methods("POST")
	router.Handle("/api/artifact/all", middleware.AuthMiddleware(http.HandlerFunc(fetch.FetchAllArtifacts))).Methods("GET")
	router.Handle("/api/artifact/{id}", middleware.AuthMiddleware(http.HandlerFunc(fetch.FetchArtifactsByID))).Methods("GET")
	router.Handle("/api/subject/all", middleware.AuthMiddleware(http.HandlerFunc(fetch.FetchAllSubjects))).Methods("GET")
	router.Handle("/api/artifact/{id}/comments", middleware.AuthMiddleware(http.HandlerFunc(fetch.FetchCommentsByArtifactID))).Methods("GET")

	fs := http.FileServer(http.Dir(fileOps.UploadPath))
	router.PathPrefix("/api/files/").Handler(http.StripPrefix("/api/files", fs))

	// Rotas pro trabalho do kaka (n usa o middleware de JWT)
	router.HandleFunc("/api/loginwithemail", login.LoginWithEmail).Methods("POST")
	router.HandleFunc("/api/signupwithemail", login.SignupWithEmail).Methods("POST")
	router.HandleFunc("/api/colourblind/new", colourblind.CheckColourBlindness).Methods("POST")
	router.HandleFunc("/api/colourblind", colourblind.FetchAllTests).Methods("GET")
	router.HandleFunc("/api/colourblind/{id}", colourblind.FetchTestsByUserID).Methods("GET")

	// router.HandleFunc("/api/teste", testeBarros).Methods("GET")
	// router.Handle("/api/testeKaka", middleware.AuthMiddleware(http.HandlerFunc(testeKaka))).Methods("GET")

	//router.HandleFunc("/", homePage)

	corsWrapper := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Origin", "Accept", "*", "Access-Control-Allow-Origin"},
	})

	//default functions
	handler := corsWrapper.Handler(router)
	http.ListenAndServe(":3001", handler)

	defer db.DB.Close()

}

// func homePage(res http.ResponseWriter, req *http.Request) {
// 	http.ServeFile(res, req, "html/index.html")
// }
