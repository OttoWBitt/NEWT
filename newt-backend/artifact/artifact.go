package artifact

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/OttoWBitt/NEWT/common"
	"github.com/OttoWBitt/NEWT/db"
	"github.com/OttoWBitt/NEWT/fileOps"
)

func InsertArtifacts(res http.ResponseWriter, req *http.Request) {

	jwtUser, err := common.GetUserFromContext(req.Context())
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	user, httpResponse, err := common.ValidateAndReturnLoggedUser(*jwtUser)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, httpResponse)
		return
	}

	usrId, err := strconv.Atoi(req.FormValue("artifactUserId"))
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	if usrId != user.Id {
		erro := "logged user id is not the same as artifactUserId"
		common.RenderResponse(res, &erro, http.StatusForbidden)
		return
	}

	var returnObj *common.Artifact
	_, _, haveFile := req.FormFile("artifactFile")

	if haveFile != http.ErrMissingFile && len(req.FormValue("artifactLink")) == 0 {

		fileName, downloadUrl, err, httpResponse := fileOps.UploadFileHandler(res, req)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, httpResponse)
			return
		}

		strSubId := req.FormValue("artifactSubjectId")
		subId, err := strconv.Atoi(strSubId)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		lastRow, err := saveArtifact(*user, subId, req.FormValue("artifactName"), req.FormValue("artifactDescription"), nil, fileName, downloadUrl)
		if err != nil {
			erro := "could not save link - " + err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		returnObj, err = returnObject(*user, subId, req.FormValue("artifactName"), req.FormValue("artifactDescription"), nil, fileName, downloadUrl, *lastRow)
		if err != nil {
			erro := "return error - " + err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

	} else if haveFile == http.ErrMissingFile && len(req.FormValue("artifactLink")) > 0 {

		strSubId := req.FormValue("artifactSubjectId")
		subId, err := strconv.Atoi(strSubId)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		artifactLink := req.FormValue("artifactLink")

		lastRow, err := saveArtifact(*user, subId, req.FormValue("artifactName"), req.FormValue("artifactDescription"), &artifactLink, nil, nil)
		if err != nil {
			erro := "could not save link - " + err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		returnObj, err = returnObject(*user, subId, req.FormValue("artifactName"), req.FormValue("artifactDescription"), &artifactLink, nil, nil, *lastRow)
		if err != nil {
			erro := "return error - " + err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

	} else if haveFile != http.ErrMissingFile && len(req.FormValue("artifactLink")) > 0 {

		fileName, downloadUrl, err, httpResponse := fileOps.UploadFileHandler(res, req)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, httpResponse)
			return
		}

		strSubId := req.FormValue("artifactSubjectId")
		subId, err := strconv.Atoi(strSubId)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		artifactLink := req.FormValue("artifactLink")

		lastRow, err := saveArtifact(*user, subId, req.FormValue("artifactName"), req.FormValue("artifactDescription"), &artifactLink, fileName, downloadUrl)
		if err != nil {
			erro := "could not save link - " + err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		returnObj, err = returnObject(*user, subId, req.FormValue("artifactName"), req.FormValue("artifactDescription"), &artifactLink, fileName,
			downloadUrl, *lastRow)
		if err != nil {
			erro := "return error - " + err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

	}

	generateJSON := map[string]interface{}{
		"data":   returnObj,
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

const queryInsertArtifactLink = `
INSERT INTO
	newt.artifact(
		user_id,
		name,
		description,
		subject_id,
		link
	)
VALUES (%d,"%s","%s",%d,"%s")
`

const queryInsertArtifactFile = `
INSERT INTO
	newt.artifact(
		user_id,
		name,
		description,
		subject_id,
		document_name,
		document_download_link
	)
VALUES (%d,"%s","%s",%d,"%s","%s")
`

const queryInsertArtifactFileAndLink = `
INSERT INTO
	newt.artifact(
		user_id,
		name,
		description,
		subject_id,
		link,
		document_name,
		document_download_link
	)
VALUES (%d,"%s","%s",%d,"%s","%s","%s")
`

func saveArtifact(user common.UserInfo, subjectId int, name, description string, link, documentName, documentDownloadLink *string) (*int64, error) {

	if link != nil && documentName == nil {

		query := fmt.Sprintf(queryInsertArtifactLink, user.Id, name, description, subjectId, *link)

		rows, err := db.DB.Exec(query)
		if err != nil {
			return nil, err
		}

		lastId, _ := rows.LastInsertId()
		return &lastId, nil

	}

	if documentName != nil && link == nil {

		query := fmt.Sprintf(queryInsertArtifactFile, user.Id, name, description, subjectId, *documentName, *documentDownloadLink)

		rows, err := db.DB.Exec(query)
		if err != nil {
			return nil, err
		}

		lastId, _ := rows.LastInsertId()
		return &lastId, nil
	}

	if documentName != nil && link != nil {

		query := fmt.Sprintf(queryInsertArtifactFileAndLink, user.Id, name, description, subjectId, *link, *documentName, *documentDownloadLink)

		rows, err := db.DB.Exec(query)
		if err != nil {
			return nil, err
		}

		lastId, _ := rows.LastInsertId()
		return &lastId, nil
	}
	return nil, nil
}

func returnObject(user common.UserInfo, subjectId int, name, description string, link, documentName, documentDownloadLink *string,
	rowId int64) (*common.Artifact, error) {

	subject, err := common.GetSubjectByID(subjectId)
	if err != nil {
		return nil, err
	}

	return &common.Artifact{
		Id:           int(rowId),
		Name:         name,
		User:         user,
		Description:  description,
		Subject:      *subject,
		Link:         link,
		DonwloadLink: documentDownloadLink,
	}, nil
}
