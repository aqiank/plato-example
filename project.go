package main

import (
        "io"
        "net/http"
        "os"
	"time"

        "plato/db"
        "plato/debug"
        "plato/entity"
)

const (
        dataDir = "pt-data"
)

const (
        createProjectTableSql = `id INTEGER PRIMARY KEY,
                                 post_id INTEGER NOT NULL,
                                 status TEXT NOT NULL,
                                 image_url TEXT NOT NULL,
                                 start_date DATETIME NOT NULL,
                                 end_date DATETIME NOT NULL`

        getProject = `SELECT * FROM project WHERE id = ?`

        insertProjectSql = `INSERT INTO project (post_id, status, image_url, start_date, end_date)
                            VALUES (?, ?, ?, ?, ?)`
)

type project struct {
        id int64
        postID int64
        status string
        imageURL string
        startDate time.Time
        endDate time.Time
}

func init() {
        if err := db.CreateTable("project", createProjectTableSql); err != nil {
                os.Exit(1)
        }

        //InsertProject("Example Title", "Lorem ipsum dolor sit amet")
}

func InsertProject(postID int64, status, imageURL string, startTime, endTime time.Time) (int64, error) {
        res, err := db.Exec(insertProjectSql, postID, status, imageURL, startTime, endTime)
        if err != nil {
                return 0, debug.Error(err)
        }
        id, err := res.LastInsertId()
        if err != nil {
                return 0, debug.Error(err)
        }
        return id, nil
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

func saveProjectImage(r *http.Request) (string, error) {
        file, header, err := r.FormFile("image")
        if err != nil {
                return "", debug.Error(err)
        }
        defer file.Close()

        folderPath := dataDir + "/project/img/"
        imageURL := folderPath + header.Filename
        if err = os.MkdirAll(folderPath, os.ModeDir | 0700); err != nil {
                return "", debug.Error(err)
        }

        output, err := os.OpenFile(imageURL, os.O_CREATE|os.O_WRONLY, os.ModeDir|0755)
	if err != nil {
		return "", debug.Error(err)
	}

	_, err = io.Copy(output, file)
	if err != nil {
		return "", debug.Error(err)
	}

        return "", nil
}
