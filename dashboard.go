package main

import (
	"net/http"
	"strconv"

	"plato/db"
	"plato/debug"
	"plato/server"
	"plato/server/session"
)

const (
	getApplicantsSQL = `SELECT project.id, pt_post_meta.value FROM project
			    INNER JOIN pt_post
			    ON project.post_id = post.id
			    INNER JOIN pt_post_meta
			    ON post_meta.post_id = post.id
			    WHERE post.author_id = ?`

)

type Applicant struct {
	PostID int64
	UserID int64
}

func getApplicants(authorID int64) []Applicant {
	var as []Applicant

	rows, err := db.Query(getApplicantsSQL, authorID)
	if err != nil {
		debug.Warn(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var a Applicant
		var userIDStr string

		if err = rows.Scan(
			&a.PostID,
			&userIDStr,
		); err != nil {
			debug.Warn(err)
			return nil
		}

		if a.UserID, err = strconv.ParseInt(userIDStr, 10, 0); err != nil {
			debug.Warn(err)
			return nil
		}

		as = append(as, a)
	}

	return as
}

func dashboardPageHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	user := session.User(r)
	if !session.IsLoggedIn(user) {
		http.Redirect(w, r, "/", 302)
		return nil, nil
	}

	if r.FormValue("what") == "" {
		return nil, server.ServePage(w, r, "dashboard", nil)
	}

	return nil, nil
}

func handleDashboard() {
        server.HandlePage("/dashboard", dashboardPageHandler)
}
