package main

import (
        "fmt"
        "os"

        "plato/db"
        "plato/entity"
)

const (
        createProjectTableSql = `id INTEGER PRIMARY KEY,
                                 post_id INTEGER NOT NULL,
                                 status TEXT NOT NULL`

        getProject = `SELECT * FROM project WHERE id = ?`
)

type project struct {
        id int64
        postID int64
        status string
}

func init() {
        db := db.Instance()

        db.CreateTable("project", createProjectTableSql)
        if db.Err != nil {
                os.Exit(1)
        }

        //InsertProject("Example Title", "Lorem ipsum dolor sit amet")
}

func InsertProject(title, content string) {
        if _, err := db.InsertPost(1, "project", "comma,separated,categories", title, content); err != nil {
                fmt.Println(err)
        }
}

func (p project) Post() entity.Post {
	post, _ := db.GetPost(p.postID)
	return post
}

func (p project) Title() string {
        return p.Post().Title()
}

func recommendedProjects() []string {
        return []string{"foo", "bar", "baz"}
}
