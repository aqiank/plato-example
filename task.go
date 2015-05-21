package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"plato/db"
	"plato/db/dateutil"
	"plato/entity"
	"plato/debug"
	"plato/server/api"
	"plato/server/session"
)

const (
	createTaskTableSQL = `id INTEGER PRIMARY KEY,
			      parent_id INTEGER NOT NULL,
			      user_id INTEGER NOT NULL,
			      source_id INTEGER NOT NULL,
			      title TEXT NOT NULL,
			      description TEXT NOT NULL,
			      start_date DATETIME NOT NULL,
			      end_date DATETIME NOT NULL,
			      done BOOLEAN NOT NULL,
			      is_milestone BOOLEAN NOT NULL,
			      restricted BOOLEAN NOT NULL,
			      created_at DATETIME NOT NULL`

	insertTaskSQL = `INSERT INTO task (parent_id, user_id, source_id, title, description, start_date, end_date, done, is_milestone, restricted, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, (SELECT datetime("now")))`

	getTasksSQL = `SELECT * FROM task WHERE source_id = ? AND parent_id = 0`
	getTaskChildrenSQL = `SELECT * FROM task WHERE parent_id = ?`

	getRemainingTasksSQL = `SELECT * FROM task WHERE source_id = ? AND parent_id = 0 AND done = 0`
	getDoneTasksSQL = `SELECT * FROM task WHERE source_id = ? AND parent_id = 0 AND done = 1`
	getMilestoneTasksSQL = `SELECT * FROM task WHERE source_id = ? AND parent_id = 0 AND is_milestone = 1`
)

type Task struct {
	ID int64
	ParentID int64
	UserID int64
	SourceID int64
	Title string
	Description string
	StartDate time.Time
	EndDate time.Time
	Done bool
	IsMilestone bool
	Restricted bool
	CreatedAt time.Time

	Children []Task
	User entity.User
}

func init() {
	if err := db.CreateTable("task", createTaskTableSQL); err != nil {
		log.Fatal(err)
	}
}

func insertTask(parentID, userID, sourceID int64, title, description string, startDate, endDate time.Time, done, isMilestone, restricted bool) error {
	if _, err := db.Exec(insertTaskSQL, parentID, userID, sourceID, title, description, startDate, endDate, done, isMilestone, restricted); err != nil {
		return debug.Error(err)
	}
	return nil
}

func getTasks(sourceID int64) []Task {
	ts := QueryTasks(getTasksSQL, sourceID)
	for i := range ts {
		ts[i].Children = getTaskChildren(ts[i].ParentID)
	}
	return ts
}

func getDoneTasks(sourceID int64) []Task {
	ts := QueryTasks(getDoneTasksSQL, sourceID)
	for i := range ts {
		ts[i].Children = getTaskChildren(ts[i].ParentID)
	}
	return ts
}

func getRemainingTasks(sourceID int64) []Task {
	ts := QueryTasks(getRemainingTasksSQL, sourceID)
	for i := range ts {
		ts[i].Children = getTaskChildren(ts[i].ParentID)
	}
	return ts
}

func getMilestoneTasks(sourceID int64) []Task {
	ts := QueryTasks(getMilestoneTasksSQL, sourceID)
	for i := range ts {
		ts[i].Children = getTaskChildren(ts[i].ParentID)
	}
	return ts
}

func getTaskChildren(parentID int64) []Task {
	return QueryTasks(getTaskChildrenSQL, parentID)
}

func QueryTasks(q string, data ...interface{}) []Task {
	var ts []Task

	rows, err := db.Query(q, data...)
	if err != nil {
		debug.Warn(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		var createdAt []byte

		if err = rows.Scan(
			&t.ID,
			&t.ParentID,
			&t.UserID,
			&t.SourceID,
			&t.Title,
			&t.Description,
			&t.StartDate,
			&t.EndDate,
			&t.Done,
			&t.IsMilestone,
			&t.Restricted,
			&createdAt,
		); err != nil {
			debug.Warn(err)
			return nil
		}

		if t.CreatedAt, err = dateutil.ParseDatetime(createdAt); err != nil {
			debug.Warn(err)
			return nil
		}

		t.User = db.GetUser(t.UserID)

		ts = append(ts, t)
	}

	return ts
}

func insertAssignee(userID, taskID int64) {
	taskIDStr := strconv.FormatInt(taskID, 10)
	if err := db.UpdateMeta("user", userID, "assignee", taskIDStr); err != nil {
		debug.Warn(err)
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request, bundle interface{}) (interface{}, error) {
	user := session.User(r)
	if !session.IsLoggedIn(user) {
		http.Redirect(w, r, "/", 302)
		return nil, api.ErrNotLoggedIn
	}

	switch r.Method {
	case "POST":
		tp := dateutil.TimeParser{}
		startDate := tp.ParseDate(r.FormValue("startDate"))
		endDate := tp.ParseDate(r.FormValue("endDate"))
		if tp.Err != nil {
			return nil, debug.Error(tp.Err)
		}

		postID, err := strconv.ParseInt(r.FormValue("postID"), 10, 0)
		if err != nil {
			return nil, debug.Error(err)
		}

		title := r.FormValue("title")
		description := r.FormValue("description")
		isMilestone := r.FormValue("isMilestone") == "true"
		restricted := r.FormValue("restricted") == "true"

		if err := insertTask(0, user.ID(), postID, title, description, startDate, endDate, false, isMilestone, restricted); err != nil {
			return nil, debug.Error(err)
		}
	default:
		http.Redirect(w, r, "/", 302)
		return nil, api.ErrInvalidMethod
	}

	http.Redirect(w, r, "/project/" + r.FormValue("postID"), 302)
	return user, nil
}

func handleTask() {
	api.Handle("/task", taskHandler)
}
