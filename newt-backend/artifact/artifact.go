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
	user := common.ValidateAndReturnLoggedUser(jwt, res)

	_, _, haveFile := req.FormFile("artifactFile")

	if haveFile != http.ErrMissingFile && len(req.FormValue("artifactLink")) == 0 {

		fileOps.UploadFileHandler(res, req)
		fmt.Println(user)

	} else if haveFile == http.ErrMissingFile && len(req.FormValue("artifactLink")) > 0 {

		strSubId := req.FormValue("artifactSubjectID")
		subId, err := strconv.Atoi(strSubId)
		if err != nil {
			common.RenderError(res, err.Error(), http.StatusInternalServerError)
		}

		err = saveLink(*user, subId, req.FormValue("artifactLink"), req.FormValue("artifactName"), req.FormValue("artifactDescription"))
		if err != nil {
			common.RenderError(res, "could not save link", http.StatusInternalServerError)
		}

		res.WriteHeader(http.StatusOK)
		res.Write([]byte("SUCCESS"))

	} else if haveFile != http.ErrMissingFile && len(req.FormValue("artifactLink")) > 0 {

		fileOps.UploadFileHandler(res, req)
		fmt.Println("asdasd")
	}
}

func saveLink(user common.UserInfo, subjectId int, link, name, description string) error {

	_, err := db.DB.Exec(
		`INSERT INTO
			newt.artifact(
				user_id,
				name,
				description,
				subject_id,
				link
			)
		VALUES (?,?,?,?,?)
	`, user.Id, name, description, subjectId, link)

	if err != nil {
		return err
	}
	return nil
}
