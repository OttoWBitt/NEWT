package artifact

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/OttoWBitt/NEWT/common"
	"github.com/OttoWBitt/NEWT/db"
	"github.com/OttoWBitt/NEWT/fileOps"
)

func InsertArtifacts(res http.ResponseWriter, req *http.Request) {

	jwt := req.FormValue("token")
	user, httpResponse, err := common.ValidateAndReturnLoggedUser(jwt)
	if err != nil {
		common.RenderResponse(res, err.Error(), httpResponse)
	}

	_, _, haveFile := req.FormFile("artifactFile")

	if haveFile != http.ErrMissingFile && len(req.FormValue("artifactLink")) == 0 {

		fileName, downloadUrl, err, httpResponse := fileOps.UploadFileHandler(res, req)
		if err != nil {
			common.RenderResponse(res, err.Error(), httpResponse)
		}

		strSubId := req.FormValue("artifactSubjectID")
		subId, err := strconv.Atoi(strSubId)
		if err != nil {
			common.RenderResponse(res, err.Error(), http.StatusInternalServerError)
		}

		err = saveArtifact(*user, subId, req.FormValue("artifactName"), req.FormValue("artifactDescription"), nil, fileName, downloadUrl)
		if err != nil {
			common.RenderResponse(res, "could not save link", http.StatusInternalServerError)
		}

		common.RenderResponse(res, "Sucess", http.StatusOK)

	} else if haveFile == http.ErrMissingFile && len(req.FormValue("artifactLink")) > 0 {

		strSubId := req.FormValue("artifactSubjectID")
		subId, err := strconv.Atoi(strSubId)
		if err != nil {
			common.RenderResponse(res, err.Error(), http.StatusInternalServerError)
		}

		artifactLink := req.FormValue("artifactLink")

		err = saveArtifact(*user, subId, req.FormValue("artifactName"), req.FormValue("artifactDescription"), &artifactLink, nil, nil)
		if err != nil {
			common.RenderResponse(res, "could not save link", http.StatusInternalServerError)
		}

		common.RenderResponse(res, "Sucess", http.StatusOK)

	} else if haveFile != http.ErrMissingFile && len(req.FormValue("artifactLink")) > 0 {

		fileName, downloadUrl, err, httpResponse := fileOps.UploadFileHandler(res, req)
		if err != nil {
			common.RenderResponse(res, err.Error(), httpResponse)
		}

		strSubId := req.FormValue("artifactSubjectID")
		subId, err := strconv.Atoi(strSubId)
		if err != nil {
			common.RenderResponse(res, err.Error(), http.StatusInternalServerError)
		}

		artifactLink := req.FormValue("artifactLink")

		err = saveArtifact(*user, subId, req.FormValue("artifactName"), req.FormValue("artifactDescription"), &artifactLink, fileName, downloadUrl)
		if err != nil {
			common.RenderResponse(res, "could not save link", http.StatusInternalServerError)
		}

		common.RenderResponse(res, "Sucess", http.StatusOK)
	}
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
VALUES (%d,%s,%s,%d,%s)
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
VALUES (%d,%s,%s,%d,%s,%s)
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
VALUES (%d,%s,%s,%d,%s,%s,%s)
`

func saveArtifact(user common.UserInfo, subjectId int, name, description string, link, documentName, documentDownloadLink *string) error {

	if link != nil && documentName == nil {

		query := fmt.Sprintf(queryInsertArtifactLink, user.Id, name, description, subjectId, *link)

		_, err := db.DB.Exec(query)
		if err != nil {
			return err
		}
	}

	if documentName != nil && link == nil {

		query := fmt.Sprintf(queryInsertArtifactFile, user.Id, name, description, subjectId, *documentName, *documentDownloadLink)

		_, err := db.DB.Exec(query)
		if err != nil {
			return err
		}
	}

	if documentName != nil && link != nil {

		query := fmt.Sprintf(queryInsertArtifactFileAndLink, user.Id, name, description, subjectId, *link, *documentName, *documentDownloadLink)

		_, err := db.DB.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}
