package fetch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/OttoWBitt/NEWT/db"
	"github.com/elgs/gosqljson"
)

func FetchFiles(res http.ResponseWriter, req *http.Request) {
	const theCase = "original"

	const queryFiles = `
	SELECT
		id,
		name,
		url,
		subject,
		user_id
	FROM 
		newt.file_metadata
	`

	files, err := gosqljson.QueryDbToMap(db.DB, theCase, queryFiles)
	if err != nil {
		log.Fatal(err)
	}

	generateJSON := map[string]interface{}{
		"files": files,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}

func FetchLinks(res http.ResponseWriter, req *http.Request) {
	const theCase = "original"

	const queryLinks = `
	SELECT
		id,
		link,
		user_id
	FROM 
		newt.links
	`

	links, err := gosqljson.QueryDbToMap(db.DB, theCase, queryLinks)
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

func FetchFilesByID(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "html/getFile.html")
		return
	}

	const theCase = "original"

	fileID := req.FormValue("fileID")

	const queryFiles = `
	SELECT
		id,
		name,
		url,
		subject,
		user_id
	FROM 
		newt.file_metadata
	WHERE
		id = %s
	`

	const queryFileComments = `
	SELECT
		id,
		user_id,
		comment,
		type_id as file_id
	FROM 
		newt.comments
	WHERE
		type = "file" AND
		type_id = %s
	`

	const queryFileReactions = `
	SELECT
		user_id,
		likes,
		deslikes,
		type_id as file_id
	FROM 
		newt.reactions
	WHERE
		type = "file"AND
		type_id = %s
	`

	qf := fmt.Sprintf(queryFiles, fileID)
	qfc := fmt.Sprintf(queryFileComments, fileID)
	qfr := fmt.Sprintf(queryFileReactions, fileID)

	files, err := gosqljson.QueryDbToMap(db.DB, theCase, qf)
	if err != nil {
		log.Fatal(err)
	}

	comments, err := gosqljson.QueryDbToMap(db.DB, theCase, qfc)
	if err != nil {
		log.Fatal(err)
	}

	reactions, err := gosqljson.QueryDbToMap(db.DB, theCase, qfr)
	if err != nil {
		log.Fatal(err)
	}

	generateJSON := map[string]interface{}{
		"files":     files,
		"comments":  comments,
		"reactions": reactions,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}

func FetchLinkByID(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "html/getLink.html")
		return
	}

	const theCase = "original"

	linkID := req.FormValue("linkID")

	const queryLinks = `
	SELECT
		id,
		link,
		user_id
	FROM 
		newt.links
	WHERE
		id = %s
	`

	const queryLinkComments = `
	SELECT
		id,
		user_id,
		comment,
		type_id as link_id
	FROM 
		newt.comments
	WHERE
		type = "link" AND
		type_id = %s
	`

	const queryLinkReactions = `
	SELECT
		user_id,
		likes,
		deslikes,
		type_id as link_id
	FROM 
		newt.reactions
	WHERE
		type = "link"AND
		type_id = %s
	`

	ql := fmt.Sprintf(queryLinks, linkID)
	qfc := fmt.Sprintf(queryLinkComments, linkID)
	qfr := fmt.Sprintf(queryLinkReactions, linkID)

	links, err := gosqljson.QueryDbToMap(db.DB, theCase, ql)
	if err != nil {
		log.Fatal(err)
	}

	comments, err := gosqljson.QueryDbToMap(db.DB, theCase, qfc)
	if err != nil {
		log.Fatal(err)
	}

	reactions, err := gosqljson.QueryDbToMap(db.DB, theCase, qfr)
	if err != nil {
		log.Fatal(err)
	}

	generateJSON := map[string]interface{}{
		"links":     links,
		"comments":  comments,
		"reactions": reactions,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}

func FetchAllByUser(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.ServeFile(res, req, "html/getAllUser.html")
		return
	}

	const theCase = "original"

	userID := req.FormValue("userID")

	// teste := req.FormValue("Teste")

	// fmt.Println(teste)

	const queryLinks = `
	SELECT
		id,
		link,
		user_id
	FROM 
		newt.links
	WHERE
		id = %s
	`
	const queryFiles = `
	SELECT
		id,
		name,
		url,
		subject,
		user_id
	FROM 
		newt.file_metadata
	WHERE
		user_id = %s
	`

	const queryComments = `
	SELECT
		id,
		user_id,
		comment,
		type_id as file_id
	FROM 
		newt.comments
	WHERE
		user_id = %s
	`

	qf := fmt.Sprintf(queryFiles, userID)
	ql := fmt.Sprintf(queryLinks, userID)
	qfc := fmt.Sprintf(queryComments, userID)

	files, err := gosqljson.QueryDbToMap(db.DB, theCase, qf)
	if err != nil {
		log.Fatal(err)
	}

	links, err := gosqljson.QueryDbToMap(db.DB, theCase, ql)
	if err != nil {
		log.Fatal(err)
	}

	comments, err := gosqljson.QueryDbToMap(db.DB, theCase, qfc)
	if err != nil {
		log.Fatal(err)
	}

	generateJSON := map[string]interface{}{
		"links":    links,
		"files":    files,
		"comments": comments,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}
