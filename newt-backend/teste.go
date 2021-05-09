package main

import (
	"encoding/json"
	"net/http"

	"github.com/OttoWBitt/NEWT/common"
)

func testeBarros(res http.ResponseWriter, req *http.Request) {

	generateJSON := map[string]interface{}{
		"GaloDoido?": "GaloDoido!",
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		common.RenderResponse(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)

}

type kaka struct {
	UserAnswer    int
	CorrectAnswer int
}

func testeKaka(res http.ResponseWriter, req *http.Request) {

	var teste []kaka
	for i := 0; i < 14; i++ {
		temp := &kaka{
			CorrectAnswer: i,
			UserAnswer:    i + 1,
		}
		teste = append(teste, *temp)
	}

	generateJSON := map[string]interface{}{
		"token":     "TESTE DO MEU NE",
		"questions": teste,
	}

	jsonData, err := json.Marshal(generateJSON)

	if err != nil {
		common.RenderResponse(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)

}
