package controllers

import (
	"net/http"
	"strconv"
)

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	// Connect to database
	db := Connect()
	defer db.Close()

	// Get value from form
	err := r.ParseForm()
	if err != nil {
		return
	}
	fullname := r.Form.Get("fullname")
	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	// encrypt password
	address := r.Form.Get("address")

	// Query
	query := "UPDATE users SET"

	if fullname != "" {
		query += " fullname = '" + fullname + "',"
	}
	if username != "" {
		query += " username = '" + username + "',"
	}
	if email != "" {
		query += " email = '" + email + "',"
	}
	if password != "" {
		query += " password = '" + password + "',"
	}
	if address != "" {
		query += " address = '" + address + "',"
	}
	queryNew := query[:len(query)-1] // Delete last coma
	userId := GetIdFromCookie(r)
	queryNew += " WHERE user_id = " + strconv.Itoa(userId)

	_, errQuery := db.Exec(queryNew)

	if errQuery != nil {
		SendErrorResponse(w, 400)
	} else {
		SendSuccessResponse(w)
	}
}
