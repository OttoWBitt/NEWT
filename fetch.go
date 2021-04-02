package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/elgs/gosqljson"
)

func fetchFiles(res http.ResponseWriter, req *http.Request) {
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

	const queryFileComments = `
	SELECT
		id,
		user_id,
		comment,
		type_id as file_id
	FROM 
		newt.comments
	WHERE
		type = "file"
	`

	const queryFileReactions = `
	SELECT
		likes,
		deslikes,
		type_id as file_id
	FROM 
		newt.reactions
	WHERE
		type = "file"
	`

	files, err := gosqljson.QueryDbToMap(dbMaster, theCase, queryFiles)
	if err != nil {
		log.Fatal(err)
	}

	comments, err := gosqljson.QueryDbToMap(dbMaster, theCase, queryFileComments)
	if err != nil {
		log.Fatal(err)
	}

	reactions, err := gosqljson.QueryDbToMap(dbMaster, theCase, queryFileReactions)
	if err != nil {
		log.Fatal(err)
	}

	generateJSON := map[string]map[string]interface{}{
		"files": {
			"files":     files,
			"comments":  comments,
			"reactions": reactions,
		},
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}

func fetchFilesByID(res http.ResponseWriter, req *http.Request) {
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

	files, err := gosqljson.QueryDbToMap(dbMaster, theCase, qf)
	if err != nil {
		log.Fatal(err)
	}

	comments, err := gosqljson.QueryDbToMap(dbMaster, theCase, qfc)
	if err != nil {
		log.Fatal(err)
	}

	reactions, err := gosqljson.QueryDbToMap(dbMaster, theCase, qfr)
	if err != nil {
		log.Fatal(err)
	}

	generateJSON := map[string]map[string]interface{}{
		"files": {
			"files":     files,
			"comments":  comments,
			"reactions": reactions,
		},
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}
