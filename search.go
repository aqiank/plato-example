package main

import (
	"net/http"

	"plato/db"
	"plato/entity"
	"plato/server/page"
	"plato/server/service"
)

const (
	searchProjectsSQL = `SELECT project.* FROM project
			    INNER JOIN pt_post ON project.post_id = pt_post.id
			    WHERE pt_post.title LIKE ?`

	searchUsersSQL = `SELECT * FROM pt_user WHERE nicename LIKE ? OR email LIKE ?`
)

func searchProjects(s string) []Project {
	return queryProjects(searchProjectsSQL, "%"+s+"%")
}

func searchUsers(s string) []entity.User {
	ss := "%"+s+"%"
	return db.QueryUsers(searchUsersSQL, ss, ss)
}

func searchPageHandler(w http.ResponseWriter, r *http.Request) error {
	s := r.FormValue("s")
	return page.Serve(w, r, "search", service.Service{"Projects": searchProjects(s), "Users": searchUsers(s)})
}

func handleSearch() {
	page.Handle("/search", searchPageHandler)
}
