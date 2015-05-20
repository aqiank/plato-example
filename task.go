package main

import (
	"log"
	"strconv"
	"time"

	"plato/db"
	"plato/db/dateutil"
	"plato/debug"
)

const (
	createTaskTableSQL = `id INTEGER PRIMARY KEY,
			      parent_id INTEGER NOT NULL,
			      source_id INTEGER NOT NULL,
			      title TEXT NOT NULL,
			      description TEXT NOT NULL,
			      status TEXT NOT NULL,
			      start_date DATETIME NOT NULL,
			      end_date DATETIME NOT NULL,
			      is_milestone BOOLEAN NOT NULL,
			      restricted BOOLEAN NOT NULL,
			      created_at DATETIME NOT NULL`

	insertTaskSQL = `INSERT INTO task (parent_id, source_id, title, description, status, start_date, end_date, is_milestone, restricted, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, (SELECT datetime("now")))`

	getTasksSQL = `SELECT * FROM task WHERE source_id = ? AND parent_id = 0`
	getTaskChildrenSQL = `SELECT * FROM task WHERE parent_id = ?`
)

type Task struct {
	ID int64
	ParentID int64
	SourceID int64
	Title string
	Description string
	Status string
	StartDate time.Time
	EndDate time.Time
	IsMilestone bool
	Restricted bool
	Children []Task
}

func init() {
	if err := db.CreateTable("task", createTaskTableSQL); err != nil {
		log.Fatal(err)
	}
}

func insertTask(parentID, sourceID int64, title, description, status string, startDate, endDate time.Time, isMilestone, restricted bool) error {
	if _, err := db.Exec(insertTaskSQL, parentID, sourceID, title, description, status, startDate, endDate, isMilestone, restricted); err != nil {
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
		var startDate, endDate []byte

		if err = rows.Scan(
			&t.ID,
			&t.ParentID,
			&t.SourceID,
			&t.Title,
			&t.Description,
			&t.Status,
			startDate,
			endDate,
			&t.IsMilestone,
			&t.Restricted,
		); err != nil {
			debug.Warn(err)
			return nil
		}

		tp := dateutil.TimeParser{}
		t.StartDate = tp.ParseDatetime(startDate)
		t.EndDate = tp.ParseDatetime(endDate)
		if tp.Err != nil {
			debug.Warn(err)
			return nil
		}

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

	// TODO
}

func handleTask() {
	api.Handle("/project/task", taskHandler)
}
