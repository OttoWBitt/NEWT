package comments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/OttoWBitt/NEWT/common"
	"github.com/OttoWBitt/NEWT/db"
)

const queryInsertComment = `
INSERT INTO
	newt.comments(
		user_id,
		artifact_id,
		comment
	)
VALUES (%d,%d,"%s")
`

func InsertComments(res http.ResponseWriter, req *http.Request) {

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

	//data json
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	var info common.Comment

	if err = json.Unmarshal(data, &info); err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	if user.Id != info.User.Id {
		erro := "logged user is not the same as received user"
		common.RenderResponse(res, &erro, http.StatusForbidden)
		return
	}

	query := fmt.Sprintf(queryInsertComment, info.User.Id, info.Artifact.Id, info.Comment)

	rows, err := db.DB.Exec(query)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC()

	lastId, _ := rows.LastInsertId()

	returnObj := &common.Comment{
		Id:       int(lastId),
		Date:     common.ConvertUtcToCustomTime(now),
		User:     info.User,
		Artifact: info.Artifact,
		Comment:  info.Comment,
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
