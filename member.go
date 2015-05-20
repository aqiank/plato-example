package main

import (
	"log"
	"strconv"

	"plato/db"
	"plato/debug"
	"plato/entity"
)

const (
	createMemberTableSQL = `post_id INTEGER NOT NULL,
				user_id INTEGER NOT NULL,
				role TEXT NOT NULL,
				status TEXT NOT NULL`

	insertMemberSQL = `INSERT INTO member (post_id, user_id, role, status)
			   VALUES (?, ?, ?, ?)`

	updateMemberStatusSQL = `UPDATE member SET status = ?
				 WHERE post_id = ? AND user_id = ?`

	getMembersSQL = `SELECT * FROM member WHERE post_id = ? AND status = ?`

	getMembersOfProjectsBySQL = `SELECT member.* FROM member
				     INNER JOIN pt_post ON pt_post.author_id = ? AND member.post_id = pt_post.id
				     WHERE status = ?`

	deleteMemberSQL = `SELECT * FROM member WHERE post_id = ? AND user_id ?`

	isMemberSQL = `SELECT COUNT(*) FROM member WHERE post_id = ? AND user_id = ?`

	hasMemberSQL = `SELECT COUNT(*) FROM member
			WHERE post_id = ? AND user_id = ? AND status = ?`
)

type Member struct {
	User entity.User

	PostID int64
	UserID int64
	Role string
	Status string
}

func init() {
	if err := db.CreateTable("member", createMemberTableSQL); err != nil {
		log.Fatal(err)
	}
}

func insertMember(postID, userID int64, role string) error {
	if _, err := db.Exec(insertMemberSQL, postID, userID, role); err != nil {
		return debug.Error(err)
	}
	return nil
}

func getMembers(postID int64, status string) []Member {
	return QueryMembers(getMembersSQL, postID, status)
}

func getMembersOfProjectsBy(authorID int64, status string) []Member {
	return QueryMembers(getMembersOfProjectsBySQL, authorID, status)
}

func deleteMember(postID, userID int64) error {
	if _, err := db.Exec(deleteMemberSQL, postID, userID); err != nil {
		return debug.Error(err)
	}
	return nil
}

func isMember(postID, userID int64) bool {
	return db.Exists(isMemberSQL, postID, userID)
}

func QueryMembers(q string, data ...interface{}) []Member {
	var ms []Member

	rows, err := db.Query(q, data...)
	if err != nil {
		debug.Warn(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var m Member

		if err = rows.Scan(
			&m.PostID,
			&m.UserID,
			&m.Role,
			&m.Status,
		); err != nil {
			debug.Warn(err)
			return nil
		}

		m.User = db.GetUser(m.UserID)

		ms = append(ms, m)
	}

	return ms
}

func supportProject(postID, userID int64) {
	if err := db.UpdateMeta("post", postID, "support", strconv.FormatInt(userID, 10)); err != nil {
		debug.Warn(err)
	}
}

func applyProject(postID, userID int64, role string) {
	if _, err := db.Exec(insertMemberSQL, postID, userID, role, "pending"); err != nil {
		debug.Warn(err)
	}
}

func joinProject(postID int64, userID int64) {
	if _, err := db.Exec(updateMemberStatusSQL, "accepted", postID, userID); err != nil {
		debug.Warn(err)
	}
}

func supportedProject(postID, userID int64) bool {
	return db.HasMeta("post", postID, "support", strconv.FormatInt(userID, 10))
}

func appliedProject(postID, userID int64) bool {
	return db.Exists(hasMemberSQL, postID, userID, "pending")
}

func joinedProject(postID, userID int64) bool {
	return db.Exists(hasMemberSQL, postID, userID, "accepted")
}
