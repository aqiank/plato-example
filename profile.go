package main

import (
	"errors"
	"net/http"
	"path"
	"strconv"

	"plato/db"
	"plato/debug"
	"plato/server/page"
	"plato/server/service"
	"plato/server/session"
)

var (
	ErrUpdateProfile = errors.New("failed to update profile")
)

func profilePageHandler(w http.ResponseWriter, r *http.Request) error {
	user := session.User(r)
	if !session.IsLoggedIn(user) {
		http.Redirect(w, r, "/", 302)
		return nil
	}

	if r.FormValue("what") == "" {
		return page.Serve(w, r, "profile", nil)
	}

	if err := user.Update(w, r); err != nil {
		return ErrUpdateProfile
	}

	return nil
}

func viewProfilePageHandler(w http.ResponseWriter, r *http.Request) error {
	user := session.User(r)
	if !session.IsLoggedIn(user) {
		http.Redirect(w, r, "/", 302)
		return nil
	}

	base := path.Base(r.URL.Path[1:])
	userID, err := strconv.ParseInt(base, 10, 0)
	if err != nil {
		return debug.Error(err)
	}

	otherUser := db.GetUser(userID)
	if otherUser == nil {
		http.Redirect(w, r, "/", 302)
		return nil
	}

	return page.Serve(w, r, "profile-view", service.Service{"OtherUser": otherUser})
}

func handleProfile() {
	page.Handle("/profile", profilePageHandler)
	page.Handle("/profile/", viewProfilePageHandler)
}
