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
	"plato/server/session"
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

	recommendedProjectsSQL = `SELECT project.* FROM project
				  INNER JOIN pt_post ON project.post_id = pt_post.id
				  WHERE recommended = 1 ORDER BY datetime(created_at) DESC LIMIT ?`

	latestRelatedProjectsSQL = `SELECT project.* FROM project
				    INNER JOIN profession ON project.post_id = profession.post_id
				    INNER JOIN pt_post ON project.post_id = pt_post.id
				    WHERE profession.name = ?
				    ORDER BY datetime(pt_post.created_at) DESC LIMIT ?`

	projectMembersSQL = `SELECT pt_user.* FROM pt_user INNER JOIN pt_post_meta
			     ON pt_user.id = pt_post_meta.value
			     WHERE pt_post_meta.key = "join" AND pt_post_meta.post_id = ?`

	// profession
	createProfessionTableSQL = `post_id INTEGER NOT NULL,
				    name TEXT NOT NULL,
				    count INTEGER NOT NULL`

	insertProfessionSQL = `INSERT INTO profession (post_id, name, count)
			       VALUES (?, ?, ?)`

	updateProfessionSQL = `UPDATE profession SET count = ? WHERE post_id = ? AND name = ?`

	getProfessionSQL = `SELECT * FROM profession WHERE post_id = ?`

	neededProfessionSQL = `SELECT count FROM profession WHERE post_id = ? AND name = ?`
)

type Project interface {
	Post() entity.Post
	Title() string
	Content() string
	ShortContent(int) string
	PostID() int64
	Tagline() string
	Status() string
	Recommended() bool
	StartDate() time.Time
	EndDate() time.Time
	ImageURL() string
	DaysLeft() int
	Started() bool
	Ended() bool
	Members() []entity.User
	FilledProfession(string) int64
	NeededProfession(string) int64
	ProfessionProgress(string) int64
	Professions() []Profession
	SupportedBy(int64) bool
	AppliedBy(int64) bool
	JoinedBy(int64) bool
	Supports() int64
}

type project struct {
	postID int64
	tagline string
	status string
	imageURL string
	recommended bool
	startDate time.Time
	endDate time.Time
	members []entity.User
}

type Profession struct {
	PostID int64
	Name string
	Count int64
}

func (p project) Post() entity.Post {
	return db.GetPost(p.postID)
}

func (p project) Title() string {
	return p.Post().Title()
}

func (p project) Content() string {
	return p.Post().Content()
}

func (p project) ShortContent(n int) string {
	return p.Post().ShortContent(n)
}

func (p project) ImageURL() string {
	return p.imageURL
}

func (p project) PostID() int64 {
	return p.postID
}

func (p project) Tagline() string {
	return p.tagline
}

func (p project) Status() string {
	return p.status
}

func (p project) Recommended() bool {
	return p.recommended
}

func (p project) StartDate() time.Time {
	return p.startDate
}

func (p project) EndDate() time.Time {
	return p.endDate
}

func (p project) DaysLeft() int {
	return int(p.endDate.Sub(time.Now()) / time.Hour / 24)
}

func (p project) Started() bool {
	return time.Since(p.startDate) >= 0
}

func (p project) Ended() bool {
	return time.Since(p.endDate) >= 0
}

func (p project) Members() []entity.User {
	if p.members == nil {
		p.members = db.QueryUsers(projectMembersSQL, p.postID)
	}
	return p.members
}

func (p project) FilledProfession(profession string) int64 {
	var count int64

	if p.members == nil {
		p.Members()
	}

	for _, member := range p.members {
		println(member.Profession() == profession)
		if member.Profession() == profession {
			count++
		}
	}

	return count
}

func (p project) NeededProfession(profession string) int64 {
	var count int64
	if err := db.QueryRow(neededProfessionSQL, p.postID, profession).Scan(&count); err != nil {
		debug.Warn(err)
		return 0
	}
	return count
}

func (p project) ProfessionProgress(profession string) int64 {
	fp := p.FilledProfession(profession)
	np := p.NeededProfession(profession)
	if np <= 0 {
		return 100
	}
	return fp * 100 / np
}

func (p project) Professions() []Profession {
	var ps []Profession

	rows, err := db.Query(getProfessionSQL, p.postID);
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

func (p project) SupportedBy(userID int64) bool {
	return supportedProject(p.postID, userID)
}

func (p project) AppliedBy(userID int64) bool {
	return appliedProject(p.postID, userID)
}

func (p project) JoinedBy(userID int64) bool {
	return joinedProject(p.postID, userID)
}

func (p project) Supports() int64 {
	count, _ := db.MetaCount("post", p.postID, "support")
	return count
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

func getProject(id int64) Project {
	var p project

	if err := db.QueryRow(getProjectSQL, id).Scan(
		&p.postID,
		&p.tagline,
		&p.status,
		&p.imageURL,
		&p.startDate,
		&p.endDate,
		&p.recommended,
	); err != nil {
		debug.Warn(err)
		return nil
	}

	return p
}

func init() {
	db := db.Instance()

	db.CreateTable("project", createProjectTableSQL)
	db.CreateTable("profession", createProfessionTableSQL)
	if db.Err != nil {
		os.Exit(1)
	}
}

func recommendedProjects(n int) []Project {
	return queryProject(recommendedProjectsSQL, n)
}

func latestRelatedProjects(profession string, n int) []Project {
	return queryProject(latestRelatedProjectsSQL, profession, n)
}

func queryProject(q string, data ...interface{}) []Project {
	var ps []Project
	var rows *sql.Rows
	var err error

	if data != nil {
		rows, err = db.Query(q, data...)
	} else {
		rows, err = db.Query(q)
	}
	if err != nil {
		debug.Warn(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var p project

		if err := rows.Scan(
			&p.postID,
			&p.tagline,
			&p.status,
			&p.imageURL,
			&p.startDate,
			&p.endDate,
			&p.recommended,
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
	base := path.Base(r.URL.Path[1:])
	id, err := strconv.ParseInt(base, 10, 0)
	if err != nil {
		debug.Warn(err)
		http.Redirect(w, r, "/", 302)
		return nil, nil
	}

	p := getProject(id)
	return nil, server.ServePage(w, r, "project", service.Service{"Project": p})
}

func projectHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	user := session.User(r)
	if !session.IsLoggedIn(user) {
		http.Redirect(w, r, "/", 302)
		return nil, nil
	}

	postIDStr := r.FormValue("postID")
	postID, _ := strconv.ParseInt(postIDStr, 10, 0)

	method := r.FormValue("method")
	switch method {
	case "support":
		supportProject(postID, user.ID())
		http.Redirect(w, r, "/project/"+postIDStr, 302)
		return postID, nil
	case "apply":
		applyProject(postID, user.ID())
		http.Redirect(w, r, "/project/"+postIDStr, 302)
		return postID, nil
	case "accept":
		if !db.IsAuthor(postID, user.ID()) {
			return postID, nil
		}
		joinProject(postID, r.FormValue("userID"))
		http.Redirect(w, r, "/project/"+postIDStr, 302)
		return postID, nil
	}

	data, err := server.PostHandler(w, r)
	if err != nil {
		return nil, debug.Error(err)
	}

	postID = data.(int64)
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
		if _, err = db.InsertActivity(user.ID(), postID, "create project"); err != nil {
			return nil, debug.Error(err)
		}

		joinProject(postID, strconv.FormatInt(user.ID(), 10))
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
	base := path.Base(r.URL.Path[1:])
	id, err := strconv.ParseInt(base, 10, 0)
	if err != nil {
		debug.Warn(err)
		http.Redirect(w, r, "/", 302)
		return nil, nil
	}

	p := getProject(id)
	return nil, server.ServePage(w, r, "project-edit", service.Service{"Project": p})
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

func supportProject(postID, authorID int64) {
	if err := db.UpdateMeta("post", postID, "support", strconv.FormatInt(authorID, 10)); err != nil {
		debug.Warn(err)
	}
}

func applyProject(postID, authorID int64) {
	if err := db.UpdateMeta("post", postID, "apply", strconv.FormatInt(authorID, 10)); err != nil {
		debug.Warn(err)
	}
}

func joinProject(postID int64, userID string) {
	if err := db.UpdateMeta("post", postID, "join", userID); err != nil {
		debug.Warn(err)
	}
}

func supportedProject(postID, authorID int64) bool {
	return db.HasMeta("post", postID, "support", strconv.FormatInt(authorID, 10))
}

func appliedProject(postID, authorID int64) bool {
	return db.HasMeta("post", postID, "apply", strconv.FormatInt(authorID, 10))
}

func joinedProject(postID, authorID int64) bool {
	return db.HasMeta("post", postID, "join", strconv.FormatInt(authorID, 10))
}

func handleProject() {
	server.SetSuccessCallback("/post/comment", commentSuccess)

        server.HandlePage("/project", projectHandler)
        server.HandlePage("/project/", projectPageHandler)
        server.HandlePage("/project/new", newProjectPageHandler)
        server.HandlePage("/project/edit/", editProjectPageHandler)
}

