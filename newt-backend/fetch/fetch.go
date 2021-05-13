package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/OttoWBitt/NEWT/common"
	"github.com/OttoWBitt/NEWT/db"
	"github.com/elgs/gosqljson"
	"github.com/gorilla/mux"
)

func FetchAllSubjects(res http.ResponseWriter, req *http.Request) {
	const theCase = "original"

	const querySubjects = `
	SELECT
		id,
		subject as name
	FROM 
		newt.subjects
	`

	subjects, err := gosqljson.QueryDbToMap(db.DB, theCase, querySubjects)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	generateJSON := map[string]interface{}{
		"data":   subjects,
		"errors": nil,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
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
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	for rows.Next() {
		artifact := new(Artifact)
		if err := rows.Scan(&artifact.Id, &artifact.Name, &artifact.UserId, &artifact.Description, &artifact.SubjectId,
			&artifact.Link, &artifact.DonwloadLink); err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusBadRequest)
			return
		}
		artifacts = append(artifacts, *artifact)
	}

	var retArtifacts []common.Artifact

	for _, art := range artifacts {
		userId, err := strconv.Atoi(art.UserId)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		user, err := common.GetUserByID(userId)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		subQuery := fmt.Sprintf(subjectQuery, art.SubjectId)

		var subject common.Subject
		var subIdStr string

		err = db.DB.QueryRow(subQuery).Scan(&subIdStr, &subject.Name)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusBadRequest)
			return
		}

		subId, err := strconv.Atoi(subIdStr)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		subject.Id = subId

		artId, err := strconv.Atoi(art.Id)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		subject.Id = subId

		artifact := &common.Artifact{
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
		"data":   retArtifacts,
		"errors": nil,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}

func FetchArtifactsByID(res http.ResponseWriter, req *http.Request) {

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
	WHERE 
		id = %s
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
	vars := mux.Vars(req)
	id, ok := vars["id"]
	if !ok {
		erro := "id is missing in parameters"
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	var artifacts []Artifact

	artQuery := fmt.Sprintf(artifactQuery, id)

	rows, err := db.DB.Query(artQuery)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	for rows.Next() {
		artifact := new(Artifact)
		if err := rows.Scan(&artifact.Id, &artifact.Name, &artifact.UserId, &artifact.Description, &artifact.SubjectId,
			&artifact.Link, &artifact.DonwloadLink); err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusBadRequest)
			return
		}
		artifacts = append(artifacts, *artifact)
	}

	var retArtifacts common.Artifact

	for _, art := range artifacts {
		userId, err := strconv.Atoi(art.UserId)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		user, err := common.GetUserByID(userId)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		subQuery := fmt.Sprintf(subjectQuery, art.SubjectId)

		var subject common.Subject
		var subIdStr string

		err = db.DB.QueryRow(subQuery).Scan(&subIdStr, &subject.Name)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusBadRequest)
			return
		}

		subId, err := strconv.Atoi(subIdStr)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		subject.Id = subId

		artId, err := strconv.Atoi(art.Id)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		subject.Id = subId

		artifact := &common.Artifact{
			Id:           artId,
			Name:         art.Name,
			User:         *user,
			Description:  art.Description,
			Subject:      subject,
			Link:         art.Link,
			DonwloadLink: art.DonwloadLink,
		}

		retArtifacts = *artifact
	}

	generateJSON := map[string]interface{}{
		"data":   retArtifacts,
		"errors": nil,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}

type comment struct {
	id         int
	userId     int
	artifactId int
	date       string
	comment    string
}

func FetchCommentsByArtifactID(res http.ResponseWriter, req *http.Request) {

	const queryComments = `
	SELECT
		id,
		user_id,
		artifact_id,
		date,
		comment
	FROM 
		newt.comments
	WHERE
		artifact_id = %s
	`

	vars := mux.Vars(req)
	id, ok := vars["id"]
	if !ok {
		erro := "id is missing in parameters"
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	comQuery := fmt.Sprintf(queryComments, id)

	rows, err := db.DB.Query(comQuery)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	var comments []comment

	for rows.Next() {
		com := new(comment)
		if err := rows.Scan(&com.id, &com.userId, &com.artifactId, &com.date, &com.comment); err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusBadRequest)
			return
		}
		comments = append(comments, *com)
	}

	var retComments []common.Comment

	for _, com := range comments {

		user, err := common.GetUserByID(com.userId)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		const layout = "2006-01-02 15:04:05"
		comDate, err := time.Parse(layout, com.date)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		comments := &common.Comment{
			Id:      com.id,
			Date:    common.ConvertUtcToCustomTime(comDate),
			User:    *user,
			Comment: com.comment,
			Artifact: common.Artifact{
				Id: com.artifactId,
			},
		}

		retComments = append(retComments, *comments)
	}

	generateJSON := map[string]interface{}{
		"data":   retComments,
		"errors": nil,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}
