package main

import (
	"net/http"
	"strconv"

	"plato/db"
	"plato/debug"
	"plato/entity"
	"plato/server/page"
	"plato/server/session"
)

const (
	getApplicantsSQL = `SELECT project.post_id, pt_post_meta.value FROM project
			    INNER JOIN pt_post
			    ON project.post_id = pt_post.id
			    INNER JOIN pt_post_meta
			    ON pt_post_meta.post_id = pt_post.id
			    WHERE pt_post.author_id = ? AND pt_post_meta.key = "apply"`
)

type Applicant struct {
	PostID int64
	UserID int64
}

func (a Applicant) User() entity.User {
	return db.GetUser(a.UserID)
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

func dashboardPageHandler(w http.ResponseWriter, r *http.Request) error {
	user := session.User(r)
	if !session.IsLoggedIn(user) {
		http.Redirect(w, r, "/", 302)
		return nil
	}

	generateProjectTimelineJSON(user)
	if r.FormValue("what") == "" {
		return page.Serve(w, r, "dashboard", nil)
	}

	return nil
}

func handleDashboard() {
	page.Handle("/dashboard", dashboardPageHandler)
}
