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
		common.RenderError(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)

}
