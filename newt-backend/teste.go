package main

import (
	"encoding/json"
	"net/http"
)

func testeBarros(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {

		generateJSON := map[string]interface{}{
			"GaloDoido?": "GaloDoido!",
		}

		jsonData, err := json.Marshal(generateJSON)

		if err != nil {
			renderError(res, err.Error(), http.StatusBadRequest)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(jsonData)
	}
}
