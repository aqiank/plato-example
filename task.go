package main

import (
	"log"
	"time"

	"plato/db"
	"plato/db/dateutil"
	"plato/debug"
)

const (
	createTaskTableSQL = `parent_id INTEGER NOT NULL,
			      post_id INTEGER NOT NULL,
			      text TEXT NOT NULL,
			      status TEXT NOT NULL,
			      start_date DATETIME NOT NULL,
			      end_date DATETIME NOT NULL,
			      is_milestone BOOLEAN NOT NULL,
			      created_at DATETIME NOT NULL`

	insertTaskSQL = `INSERT INTO task (parent_id, post_id, text, status, start_date, end_date, is_milestone, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, (SELECT datetime("now")))`

	getTasksSQL = `SELECT * FROM task WHERE post_id = ? AND parent_id = 0`
	getTaskChildrenSQL = `SELECT * FROM task WHERE parent_id = ?`
)

type Task struct {
	ParentID int64
	PostID int64
	Text string
	Status string
	StartDate time.Time
	EndDate time.Time
	IsMilestone bool
	Children []Task
}

func init() {
	if err := db.CreateTable("task", createTaskTableSQL); err != nil {
		log.Fatal(err)
	}
}

func insertTask(parentID, postID int64, status, text string, startDate, endDate time.Time, isMilestone bool) error {
	if _, err := db.Exec(insertTaskSQL, parentID, postID, text, startDate, endDate, isMilestone); err != nil {
		return debug.Error(err)
	}
	return nil
}

func getTasks(postID int64) []Task {
	ts := QueryTasks(getTasksSQL, postID)
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
			&t.ParentID,
			&t.PostID,
			&t.Text,
			&t.Status,
			startDate,
			endDate,
			&t.IsMilestone,
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
