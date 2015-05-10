package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"plato/db"
	"plato/db/dateutil"
	"plato/debug"
	"plato/entity"
	"plato/server"
	"plato/server/service"
)

const (
	// project
	createProjectTableSQL = `post_id INTEGER NOT NULL,
				 tagline TEXT NOT NULL,
				 status TEXT NOT NULL,
				 image_url TEXT NOT NULL,
				 start_date DATETIME NOT NULL,
				 end_date DATETIME NOT NULL,
				 recommended BOOLEAN NOT NULL`

	insertProjectSQL = `INSERT INTO project (post_id, tagline, status, image_url, start_date, end_date, recommended)
			    VALUES (?, ?, ?, ?, ?, ?, ?)`

	updateProjectSQL = `UPDATE project SET tagline = ?, status = ?, image_url = ?, start_date = ?, end_date = ? WHERE post_id = ?`

	updateProjectWithoutImageSQL = `UPDATE project SET tagline = ?, status = ?, start_date = ?, end_date = ? WHERE post_id = ?`

	getProjectSQL = `SELECT * FROM project WHERE post_id = ?`

	recommendedProjectsSQL = `SELECT * FROM project WHERE recommended = 1`

	// profession
	createProfessionTableSQL = `post_id INTEGER NOT NULL,
				    name TEXT NOT NULL,
				    count INTEGER NOT NULL`

	insertProfessionSQL = `INSERT INTO profession (post_id, name, count)
			       VALUES (?, ?, ?)`

	updateProfessionSQL = `UPDATE profession SET count = ? WHERE post_id = ? AND name = ?`

	getProfessionSQL = `SELECT * FROM profession WHERE post_id = ?`
)

type Project struct {
	PostID int64
	Tagline string
	Status string
	ImageURL string
	Recommended bool
	StartDate time.Time
	EndDate time.Time
}

type Profession struct {
	PostID int64
	Name string
	Count int64
}

func (p Project) Post() entity.Post {
	post, _ := db.GetPost(p.PostID)
	return post
}

func (p Project) Title() string {
	return p.Post().Title()
}

func (p Project) Content() string {
	return p.Post().Content()
}

func (p Project) ShortContent(n int) string {
	return p.Post().ShortContent(n)
}

func (p Project) DaysLeft() int {
	return int(p.EndDate.Sub(time.Now()) / time.Hour / 24)
}

func (p Project) Started() bool {
	return time.Since(p.StartDate) >= 0
}

func (p Project) Ended() bool {
	return time.Since(p.EndDate) >= 0
}

func (p Project) Professions() []Profession {
	var ps []Profession

	rows, err := db.Query(getProfessionSQL, p.PostID);
	if err != nil {
		debug.Warn(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var p Profession
		if rows.Scan(
			&p.PostID,
			&p.Name,
			&p.Count,
		); err != nil {
			debug.Warn(err)
			return nil
		}

		ps = append(ps, p)
	}

	return ps
}

func insertProject(postID int64, tagline, status, imageURL string, startTime, endTime time.Time) (int64, error) {
	res, err := db.Exec(insertProjectSQL, postID, tagline, status, imageURL, startTime, endTime, false)
	if err != nil {
		return 0, debug.Error(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, debug.Error(err)
	}

	return id, nil
}

func updateProject(postID int64, tagline, status, imageURL string, startTime, endTime time.Time) error {
	if imageURL != "" {
		if _, err := db.Exec(updateProjectSQL, tagline, status, imageURL, startTime, endTime, postID); err != nil {
			return debug.Error(err)
		}
	} else {
		if _, err := db.Exec(updateProjectWithoutImageSQL, tagline, status, startTime, endTime, postID); err != nil {
			return debug.Error(err)
		}
	}
	return nil
}

func getProject(id int64) (Project, error) {
	var p Project

	if err := db.QueryRow(getProjectSQL, id).Scan(
		&p.PostID,
		&p.Tagline,
		&p.Status,
		&p.ImageURL,
		&p.StartDate,
		&p.EndDate,
		&p.Recommended,
	); err != nil {
		return p, debug.Error(err)
	}

	return p, nil
}

func init() {
	db := db.Instance()

	db.CreateTable("project", createProjectTableSQL)
	db.CreateTable("profession", createProfessionTableSQL)
	if db.Err != nil {
		os.Exit(1)
	}
}

func recommendedProjects() []Project {
	return queryProject(recommendedProjectsSQL)
}

func queryProject(q string, data ...interface{}) []Project {
	var ps []Project
	var rows *sql.Rows
	var err error

	if data != nil {
		rows, err = db.Query(q, data)
	} else {
		rows, err = db.Query(q)
	}
	if err != nil {
		debug.Warn(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var p Project

		if err := rows.Scan(
			&p.PostID,
			&p.Tagline,
			&p.Status,
			&p.ImageURL,
			&p.StartDate,
			&p.EndDate,
			&p.Recommended,
		); err != nil {
			debug.Warn(err)
			return nil
		}

		ps = append(ps, p)
	}

	return ps
}

func saveProjectImage(id int64, r *http.Request) (string, error) {
	folderPath := fmt.Sprintf("%s/%s/%d/", db.DataDir, "project/img", id)
	imageURL, err := db.SaveImage(folderPath, "image", r)
	if err != nil {
		return "", debug.Error(err)
	}
	return "/" + imageURL, nil
}

func newProjectPageHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
        return nil, server.ServePage(w, r, "project-new", nil)
}

func projectPageHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var p Project

	base := path.Base(r.URL.Path[1:])
	id, err := strconv.ParseInt(base, 10, 0)
	if err != nil {
		goto out
	}

	p, err = getProject(id)
	if err != nil {
		goto out
	}

	return nil, server.ServePage(w, r, "project", service.Service{"Project": p})
out:
	debug.Warn(err)
	http.Redirect(w, r, "/", 302)
	return nil, nil
}

func projectHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	data, err := server.PostHandler(w, r)
	if err != nil {
		return nil, debug.Error(err)
	}

	postID := data.(int64)
	imageURL, _ := saveProjectImage(postID, r)

	tp := dateutil.TimeParser{}
	startDate := tp.ParseDate(r.FormValue("startDate"))
	endDate := tp.ParseDate(r.FormValue("endDate"))
	if tp.Err != nil {
		return nil, debug.Error(tp.Err)
	}

	var id int64
	tagline := r.FormValue("tagline")
	status := r.FormValue("status")
	switch r.FormValue("method") {
	case "POST":
		if id, err = insertProject(postID, tagline, status, imageURL, startDate, endDate); err != nil {
			return nil, debug.Error(err)
		}
		if err = insertProfession(id, r); err != nil {
			return nil, debug.Error(err)
		}
		http.Redirect(w, r, fmt.Sprintf("%s%d", "/project/", postID), 302)
	case "PUT":
		if err = updateProject(postID, tagline, status, imageURL, startDate, endDate); err != nil {
			return nil, debug.Error(err)
		}
		if err = updateProfession(postID, r); err != nil {
			return nil, debug.Error(err)
		}
		http.Redirect(w, r, fmt.Sprintf("%s%d", "/project/edit/", postID), 302)
	case "GET":
		// TODO
	}

	return id, nil
}

func editProjectPageHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	var p Project

	base := path.Base(r.URL.Path[1:])
	id, err := strconv.ParseInt(base, 10, 0)
	if err != nil {
		goto out
	}

	p, err = getProject(id)
	if err != nil {
		goto out
	}

	return nil, server.ServePage(w, r, "project-edit", service.Service{"Project": p})
out:
	debug.Warn(err)
	http.Redirect(w, r, "/", 302)
	return nil, nil
}

func insertProfession(postID int64, r *http.Request) error {
	for k, v := range r.Form {
		if len(v) == 0 || !strings.Contains(k, "profession") {
			continue
		}

		// check if there's space character
		idx := strings.IndexRune(k, ' ')
		if idx == -1 || idx + 1 >= len(k) {
			continue
		}
		idx++

		cnt, err := strconv.ParseInt(v[0], 10, 0)
		if err != nil {
			return debug.Error(err)
		}

		if _, err = db.Exec(insertProfessionSQL, postID, k[idx:], cnt); err != nil {
			return debug.Error(err)
		}
	}

	return nil
}

func updateProfession(postID int64, r *http.Request) error {
	for k, v := range r.Form {
		if len(v) == 0 || !strings.Contains(k, "profession") {
			continue
		}

		// check if there's space character
		idx := strings.IndexRune(k, ' ')
		if idx == -1 || idx + 1 >= len(k) {
			continue
		}
		idx++

		cnt, err := strconv.ParseInt(v[0], 10, 0)
		if err != nil {
			return debug.Error(err)
		}

		if _, err = db.Exec(updateProfessionSQL, cnt, postID, k[idx:]); err != nil {
			return debug.Error(err)
		}
	}

	return nil
}

func handleProject() {
	server.SetSuccessCallback("/post/comment", commentSuccess)

        server.HandlePage("/project", projectHandler)
        server.HandlePage("/project/", projectPageHandler)
        server.HandlePage("/project/new", newProjectPageHandler)
        server.HandlePage("/project/edit/", editProjectPageHandler)
}

func commentSuccess(w http.ResponseWriter, r *http.Request, data interface{}) error {
	id, ok := data.(int64)
	if !ok {
		http.Redirect(w, r, "/", 302)
		return nil
	}

	url := fmt.Sprintf("/project/%d", id)
	http.Redirect(w, r, url, 302)
	return nil
}
