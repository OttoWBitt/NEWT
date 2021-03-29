package main

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
