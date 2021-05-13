package colourblind

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/OttoWBitt/NEWT/common"
	"github.com/OttoWBitt/NEWT/db"
	"github.com/gorilla/mux"
)

type InputData struct {
	Questions []Questions `json:"questions"`
	UserId    string      `json:"user_id"`
}
type Questions struct {
	Useranswer    int `json:"UserAnswer"`
	Correctanswer int `json:"CorrectAnswer"`
}

const queryInsertColourBlindResult = `
INSERT INTO
	newt.colour_blind_test(
		user_id,
		points,
		is_colour_blind
	)
VALUES (%d,%d,%t)
`

func CheckColourBlindness(res http.ResponseWriter, req *http.Request) {

	//data json
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	var info InputData

	if err = json.Unmarshal(data, &info); err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	questions := info.Questions
	userId := info.UserId

	userid, err := strconv.Atoi(userId)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	var points = 0

	for _, ques := range questions {
		if ques.Useranswer == ques.Correctanswer {
			points = points + 1
		}
	}

	var probablyColourBlind = false
	if points < 12 {
		probablyColourBlind = true
	}

	query := fmt.Sprintf(queryInsertColourBlindResult, userid, points, probablyColourBlind)

	_, err = db.DB.Exec(query)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusInternalServerError)
		return
	}

	generateJSON := map[string]interface{}{
		"points":               points,
		"probablyColourBlind?": probablyColourBlind,
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

type test struct {
	id            string
	userId        string
	points        string
	date          string
	isColourBlind bool
}

type ReturnTests struct {
	Id                    int             `json:"id"`
	User                  common.UserInfo `json:"user"`
	Points                int             `json:"points"`
	Date                  string          `json:"date"`
	IsProbablyColourBlind bool            `json:"isProbablyColourBlind"`
}

const colourBlindQuery = `
	SELECT
		id,
		user_id,
		points,
		date,
		is_colour_blind
	FROM 
		newt.colour_blind_test
`

func FetchAllTests(res http.ResponseWriter, req *http.Request) {

	var tests []test

	rows, err := db.DB.Query(colourBlindQuery)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	for rows.Next() {
		test := new(test)
		if err := rows.Scan(&test.id, &test.userId, &test.points, &test.date, &test.isColourBlind); err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusBadRequest)
			return
		}
		tests = append(tests, *test)
	}

	var returnTes []ReturnTests

	for _, tes := range tests {
		userId, err := strconv.Atoi(tes.userId)
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

		testId, err := strconv.Atoi(tes.id)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		testPoints, err := strconv.Atoi(tes.points)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		const layout = "2006-01-02 15:04:05"
		tesDate, err := time.Parse(layout, tes.date)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		test := &ReturnTests{
			Id:                    testId,
			User:                  *user,
			Points:                testPoints,
			Date:                  common.ConvertUtcToCustomTimeDate(tesDate),
			IsProbablyColourBlind: tes.isColourBlind,
		}

		returnTes = append(returnTes, *test)
	}

	generateJSON := map[string]interface{}{
		"data": returnTes,
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

const colourBlindQueryByID = `
	SELECT
		id,
		user_id,
		points,
		date,
		is_colour_blind
	FROM 
		newt.colour_blind_test
	WHERE
		user_id = %s
`

func FetchTestsByUserID(res http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id, ok := vars["id"]
	if !ok {
		erro := "id is missing in parameters"
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	var tests []test

	query := fmt.Sprintf(colourBlindQueryByID, id)

	rows, err := db.DB.Query(query)
	if err != nil {
		erro := err.Error()
		common.RenderResponse(res, &erro, http.StatusBadRequest)
		return
	}

	for rows.Next() {
		test := new(test)
		if err := rows.Scan(&test.id, &test.userId, &test.points, &test.date, &test.isColourBlind); err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusBadRequest)
			return
		}
		tests = append(tests, *test)
	}

	var returnTes []ReturnTests

	for _, tes := range tests {
		userId, err := strconv.Atoi(tes.userId)
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

		testId, err := strconv.Atoi(tes.id)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		testPoints, err := strconv.Atoi(tes.points)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		const layout = "2006-01-02 15:04:05"
		tesDate, err := time.Parse(layout, tes.date)
		if err != nil {
			erro := err.Error()
			common.RenderResponse(res, &erro, http.StatusInternalServerError)
			return
		}

		test := &ReturnTests{
			Id:                    testId,
			User:                  *user,
			Points:                testPoints,
			Date:                  common.ConvertUtcToCustomTimeDate(tesDate),
			IsProbablyColourBlind: tes.isColourBlind,
		}

		returnTes = append(returnTes, *test)
	}

	generateJSON := map[string]interface{}{
		"data": returnTes,
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
