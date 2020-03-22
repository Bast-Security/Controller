package main

func getSetting(name string) (value string, err error) {
	row := db.QueryRow("SELECT value FROM Settings WHERE name = ?;", name)
	err = row.Scan(&value)
	return
}

