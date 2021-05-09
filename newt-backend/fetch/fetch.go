package fetch

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/OttoWBitt/NEWT/common"
	"github.com/OttoWBitt/NEWT/db"
	"github.com/elgs/gosqljson"
)

func FetchAllSubjects(res http.ResponseWriter, req *http.Request) {
	const theCase = "original"

	const querySubjects = `
	SELECT
		id,
		subject
	FROM 
		newt.subjects
	`

	subjects, err := gosqljson.QueryDbToMap(db.DB, theCase, querySubjects)
	if err != nil {
		log.Fatal(err)
	}

	generateJSON := map[string]interface{}{
		"subjects": subjects,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}

type Artifact struct {
	Id           string
	Name         string
	UserId       string
	Description  string
	SubjectId    string
	Link         *string
	DonwloadLink *string
}

type Subject struct {
	Id   int
	Name string
}

type returnArtifacts struct {
	Id           int
	Name         string
	User         common.UserInfo
	Description  string
	Subject      Subject
	Link         *string
	DonwloadLink *string
}

func FetchAllArtifacts(res http.ResponseWriter, req *http.Request) {

	const artifactQuery = `
	SELECT
		id,
		name,
		user_id,
		description,
		subject_id,
		link,
		document_download_link
	FROM 
		newt.artifact
	`

	const subjectQuery = `
	SELECT
		id,
		subject
	FROM
		newt.subjects
	WHERE
		id = %s
	`

	var artifacts []Artifact

	rows, err := db.DB.Query(artifactQuery)
	if err != nil {
		common.RenderError(res, err.Error(), http.StatusBadRequest)
		return
	}

	for rows.Next() {
		artifact := new(Artifact)
		if err := rows.Scan(&artifact.Id, &artifact.Name, &artifact.UserId, &artifact.Description, &artifact.SubjectId,
			&artifact.Link, &artifact.DonwloadLink); err != nil {
			common.RenderError(res, err.Error(), http.StatusBadRequest)
		}
		artifacts = append(artifacts, *artifact)
	}

	var retArtifacts []returnArtifacts

	for _, art := range artifacts {
		userId, err := strconv.Atoi(art.UserId)
		if err != nil {
			common.RenderError(res, err.Error(), http.StatusInternalServerError)
		}

		user, err := common.GetUserByID(userId)
		if err != nil {
			common.RenderError(res, err.Error(), http.StatusInternalServerError)
		}

		subQuery := fmt.Sprintf(subjectQuery, art.SubjectId)

		var subject Subject
		var subIdStr string

		err = db.DB.QueryRow(subQuery).Scan(&subIdStr, &subject.Name)
		if err != nil {
			common.RenderError(res, err.Error(), http.StatusBadRequest)
			return
		}

		subId, err := strconv.Atoi(subIdStr)
		if err != nil {
			common.RenderError(res, err.Error(), http.StatusInternalServerError)
		}

		subject.Id = subId

		artId, err := strconv.Atoi(art.Id)
		if err != nil {
			common.RenderError(res, err.Error(), http.StatusInternalServerError)
		}

		subject.Id = subId

		artifact := &returnArtifacts{
			Id:           artId,
			Name:         art.Name,
			User:         *user,
			Description:  art.Description,
			Subject:      subject,
			Link:         art.Link,
			DonwloadLink: art.DonwloadLink,
		}

		retArtifacts = append(retArtifacts, *artifact)
	}

	generateJSON := map[string]interface{}{
		"artifacts": retArtifacts,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}

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
