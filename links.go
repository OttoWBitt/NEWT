package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/elgs/gosqljson"
)

func insertLinks(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "html/link.html")
		return
	}

	stringId := req.FormValue("id")
	link := req.FormValue("link")

	id, err := strconv.Atoi(stringId)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	auth := userExists(id)

	fmt.Println(auth) //esse cara tem que vir do cookie do front, falta implementar

	if auth != nil {
		http.Error(res, "User not logged in", http.StatusInternalServerError)
		return
	}

	_, err = dbMaster.Exec(
		`INSERT INTO 
			newt.links(
				link, 
				user_id
			) 
		VALUES (?,?)
	`, link, id)

	if err != nil {
		http.Error(res, "Server error, unable to save your link.", http.StatusInternalServerError)
		return
	}

	res.Write([]byte("link saved!"))
}

func retrieveLinks(res http.ResponseWriter, req *http.Request) {

	const theCase = "original"

	const query = `
	SELECT
		*
	FROM 
		newt.links 
	`

	links, err := gosqljson.QueryDbToMap(dbMaster, theCase, query)
	if err != nil {
		log.Fatal(err)
	}

	generateJSON := map[string]interface{}{
		"links": links,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}
