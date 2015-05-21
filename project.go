package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"time"

	"plato/db"
	"plato/db/dateutil"
	"plato/debug"
	"plato/entity"
	"plato/server/api"
	"plato/server/page"
	"plato/server/service"
	"plato/server/session"
)

var (
	ErrNotAuthor = errors.New("user is not the author")
)

const (
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

	getProjectSQL = `SELECT * FROM project WHERE post_id = ? LIMIT 1`

	getProjectsByAuthorIDSQL = `SELECT project.* FROM project
				    INNER JOIN pt_post
				    ON project.post_id = pt_post.id
				    WHERE pt_post.author_id = ?`

	recommendedProjectsSQL = `SELECT project.* FROM project
				  INNER JOIN pt_post ON project.post_id = pt_post.id
				  WHERE recommended = 1 ORDER BY datetime(created_at) DESC LIMIT ?`

	latestRelatedProjectsSQL = `SELECT project.* FROM project
				    INNER JOIN requirement ON project.post_id = requirement.post_id
				    INNER JOIN pt_post ON project.post_id = pt_post.id
				    WHERE requirement.name = ?
				    ORDER BY datetime(pt_post.created_at) DESC LIMIT ?`
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
	Members() []Member
	FilledRequirement(string) int64
	NeededRequirement(string) int64
	RequirementProgress(string) int64
	Requirements() []Requirement
	Tasks() []Task
	SupportedBy(int64) bool
	AppliedBy(int64) bool
	JoinedBy(int64) bool
	Supports() int64
}

type project struct {
	postID      int64
	tagline     string
	status      string
	imageURL    string
	recommended bool
	startDate   time.Time
	endDate     time.Time
	members     []Member
}

func init() {
	if err := db.CreateTable("project", createProjectTableSQL); err != nil {
		log.Fatal(err)
	}
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

func (p project) Members() []Member {
	return getMembers(p.postID, "accepted")
}

func (p project) FilledRequirement(profession string) int64 {
	var count int64

	for _, member := range p.Members() {
		if member.Role == profession {
			count++
		}
	}

	return count
}

func (p project) NeededRequirement(profession string) int64 {
	var count int64
	if err := db.QueryRow(neededRequirementSQL, p.postID, profession).Scan(&count); err != nil {
		debug.Warn(err)
		return 0
	}
	return count
}

func (p project) RequirementProgress(profession string) int64 {
	fp := p.FilledRequirement(profession)
	np := p.NeededRequirement(profession)
	if np <= 0 {
		return 100
	}
	return fp * 100 / np
}

func (p project) Requirements() []Requirement {
	var ps []Requirement

	rows, err := db.Query(getRequirementSQL, p.postID)
	if err != nil {
		debug.Warn(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var p Requirement
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

func (p project) Tasks() []Task {
	return getTasks(p.postID)
}

func (p project) DoneTasks() []Task {
	return getDoneTasks(p.postID)
}

func (p project) RemainingTasks() []Task {
	return getRemainingTasks(p.postID)
}

func (p project) MilestoneTasks() []Task {
	return getMilestoneTasks(p.postID)
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

func getProjectsByAuthorID(authorID int64) []Project {
	return queryProjects(getProjectsByAuthorIDSQL, authorID)
}

func recommendedProjects(n int) []Project {
	return queryProjects(recommendedProjectsSQL, n)
}

func latestRelatedProjects(profession string, n int) []Project {
	return queryProjects(latestRelatedProjectsSQL, profession, n)
}

func queryProjects(q string, data ...interface{}) []Project {
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
	folderPath := fmt.Sprintf("%s/%s/%d", db.DataDir, "project/img", id)
	imageURL, err := db.SaveImage(folderPath, "image", r)
	if err != nil {
		return "", debug.Error(err)
	}
	return "/" + imageURL, nil
}

func newProjectPageHandler(w http.ResponseWriter, r *http.Request) error {
	return page.Serve(w, r, "project-new", nil)
}

func projectPageHandler(w http.ResponseWriter, r *http.Request) error {
	base := path.Base(r.URL.Path[1:])
	id, err := strconv.ParseInt(base, 10, 0)
	if err != nil {
		debug.Warn(err)
		http.Redirect(w, r, "/", 302)
		return nil
	}

	p := getProject(id)
	return page.Serve(w, r, "project", service.Service{"Project": p})
}

func projectHandler(w http.ResponseWriter, r *http.Request, bundle interface{}) (interface{}, error) {
	user := session.User(r)
	if !session.IsLoggedIn(user) {
		http.Redirect(w, r, "/", 302)
		return nil, api.ErrNotLoggedIn
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
		role := r.FormValue("role")
		applyProject(postID, user.ID(), role)
		http.Redirect(w, r, "/project/"+postIDStr, 302)
		return postID, nil
	case "accept", "decline":
		if !db.IsAuthor(postID, user.ID()) {
			return postID, ErrNotAuthor
		}

		applicantUserID, err := strconv.ParseInt(r.FormValue("applicantUserID"), 10, 0)
		if err != nil {
			return postID, ErrNotAuthor
		}

		if method == "accept" {
			joinProject(postID, applicantUserID)
		} else if method == "decline" {
			deleteMember(postID, applicantUserID)
		}
		http.Redirect(w, r, "/dashboard", 302)
		return postID, nil
	}

	data, err := api.PostHandler(w, r, nil)
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
		if err = insertRequirement(id, r); err != nil {
			return nil, debug.Error(err)
		}
		if _, err = db.InsertActivity(user.ID(), postID, "create project"); err != nil {
			return nil, debug.Error(err)
		}

		// generate timeline
		generateTimeline(user)
		joinProject(postID, user.ID())
		http.Redirect(w, r, fmt.Sprintf("%s%d", "/project/", postID), 302)
	case "PUT":
		if err = updateProject(postID, tagline, status, imageURL, startDate, endDate); err != nil {
			return nil, debug.Error(err)
		}
		if err = updateRequirement(postID, r); err != nil {
			return nil, debug.Error(err)
		}
		http.Redirect(w, r, fmt.Sprintf("%s%d", "/project/edit/", postID), 302)
	case "GET":
		// TODO
	}

	return id, nil
}

func editProjectPageHandler(w http.ResponseWriter, r *http.Request) error {
	base := path.Base(r.URL.Path[1:])
	id, err := strconv.ParseInt(base, 10, 0)
	if err != nil {
		debug.Warn(err)
		http.Redirect(w, r, "/", 302)
		return nil
	}

	p := getProject(id)
	return page.Serve(w, r, "project-edit", service.Service{"Project": p})
}
func onComment(w http.ResponseWriter, r *http.Request, data interface{}) (interface{}, error) {
	id, ok := data.(int64)
	if !ok {
		http.Redirect(w, r, "/", 302)
		return nil, nil
	}

	url := fmt.Sprintf("/project/%d", id)
	http.Redirect(w, r, url, 302)
	return id, nil
}

func handleProject() {
	api.Append("/post/comment", onComment)
	api.Handle("/project", projectHandler)

	page.Handle("/project/", projectPageHandler)
	page.Handle("/project/new", newProjectPageHandler)
	page.Handle("/project/edit/", editProjectPageHandler)
}
