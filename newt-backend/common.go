package main

import "net/http"

func userExists(id int) error {

	var retID int

	err := dbMaster.QueryRow(
		`SELECT 
			id 
		FROM 
			newt.users 
		WHERE 
			id = ?
	`, id).Scan(&retID)

	if err != nil {
		return err
	}
	return nil
}

func getUserIDByEmail(email string) (string, error) {

	var id string

	err := dbMaster.QueryRow(
		`SELECT 
			id 
		FROM 
			newt.users 
		WHERE 
			email = ?
	`, email).Scan(&id)

	if err != nil {
		return "", err
	}
	return id, nil
}

func renderError(res http.ResponseWriter, message string, statusCode int) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte(message))
}
