package fetch

import (
	"encoding/json"
	"fmt"
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

	var retArtifacts []common.ReturnArtifacts

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

		artifact := &common.ReturnArtifacts{
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
