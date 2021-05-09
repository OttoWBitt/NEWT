package colourblind

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/OttoWBitt/NEWT/common"
	"github.com/OttoWBitt/NEWT/db"
)

type InputData struct {
	Questions []Questions `json:"questions"`
	Token     string      `json:"token"`
}
type Questions struct {
	Useranswer    int `json:"UserAnswer"`
	Correctanswer int `json:"CorrectAnswer"`
}

const queryInsertColourBlindResult = `
INSERT INTO
	newt.colour_blind_test(
		user_id,
		points
	)
VALUES (%d,%d)
`

func CheckColourBlindness(res http.ResponseWriter, req *http.Request) {

	//data json
	data, erro := ioutil.ReadAll(req.Body)
	if erro != nil {
		common.RenderResponse(res, erro.Error(), http.StatusInternalServerError)
		return
	}

	var info InputData

	if erro = json.Unmarshal(data, &info); erro != nil {
		common.RenderResponse(res, erro.Error(), http.StatusInternalServerError)
		return
	}

	//token := info.Token
	questions := info.Questions

	// user, httpResponse, err := common.ValidateAndReturnLoggedUser(token)
	// if err != nil {
	// 	common.RenderResponse(res, err.Error(), httpResponse)
	// 	return
	// }

	var points = 0

	for _, ques := range questions {
		if ques.Useranswer == ques.Correctanswer {
			points = points + 1
		}
	}

	query := fmt.Sprintf(queryInsertColourBlindResult, 1, points)

	_, err := db.DB.Exec(query)
	if err != nil {
		common.RenderResponse(res, err.Error(), http.StatusInternalServerError)
		return
	}

	common.RenderResponse(res, "Sucess", http.StatusOK)
}
