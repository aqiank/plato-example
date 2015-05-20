package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"plato/db"
	"plato/debug"
)

const (
	createRequirementTableSQL = `post_id INTEGER NOT NULL,
				     name TEXT NOT NULL,
				     count INTEGER NOT NULL`

	insertRequirementSQL = `INSERT INTO requirement (post_id, name, count)
				VALUES (?, ?, ?)`

	updateRequirementSQL = `UPDATE requirement SET count = ? WHERE post_id = ? AND name = ?`

	getRequirementSQL = `SELECT * FROM requirement WHERE post_id = ?`

	neededRequirementSQL = `SELECT count FROM requirement WHERE post_id = ? AND name = ?`
)

type Requirement struct {
	PostID int64
	Name   string
	Count  int64
}

func init() {
	if err := db.CreateTable("requirement", createRequirementTableSQL); err != nil {
		log.Fatal(err)
	}
}

func insertRequirement(postID int64, r *http.Request) error {
	for k, v := range r.Form {
		if len(v) == 0 || !strings.Contains(k, "profession") {
			continue
		}

		// check if there's space character
		idx := strings.IndexRune(k, ' ')
		if idx == -1 || idx+1 >= len(k) {
			continue
		}
		idx++

		cnt, err := strconv.ParseInt(v[0], 10, 0)
		if err != nil {
			return debug.Error(err)
		}

		if _, err = db.Exec(insertRequirementSQL, postID, k[idx:], cnt); err != nil {
			return debug.Error(err)
		}
	}

	return nil
}

func updateRequirement(postID int64, r *http.Request) error {
	for k, v := range r.Form {
		if len(v) == 0 || !strings.Contains(k, "profession") {
			continue
		}

		// check if there's space character
		idx := strings.IndexRune(k, ' ')
		if idx == -1 || idx+1 >= len(k) {
			continue
		}
		idx++

		cnt, err := strconv.ParseInt(v[0], 10, 0)
		if err != nil {
			return debug.Error(err)
		}

		if _, err = db.Exec(updateRequirementSQL, cnt, postID, k[idx:]); err != nil {
			return debug.Error(err)
		}
	}

	return nil
}
