package main

import (
	"errors"
	"net/http"
	"path"
	"strconv"

	"plato/debug"
	"plato/db"
	"plato/server"
	"plato/server/service"
	"plato/server/session"
)

var (
	ErrUpdateProfile = errors.New("failed to update profile")
)

func profilePageHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	user := session.User(r)
	if !session.IsLoggedIn(user) {
		http.Redirect(w, r, "/", 302)
		return nil, nil
	}

	if r.FormValue("what") == "" {
		return nil, server.ServePage(w, r, "profile", nil)
	}

	if err := user.Update(w, r); err != nil {
		return nil, server.ErrorCallback(w, r, ErrUpdateProfile)
	}

	return nil, nil
}

func viewProfilePageHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	user := session.User(r)
	if !session.IsLoggedIn(user) {
		http.Redirect(w, r, "/", 302)
		return nil, nil
	}

	base := path.Base(r.URL.Path[1:])
	userID, err := strconv.ParseInt(base, 10, 0)
	if err != nil {
		return nil, debug.Error(err)
	}

	otherUser := db.GetUser(userID)
	if otherUser == nil {
		http.Redirect(w, r, "/", 302)
		return nil, nil
	}

	return nil, server.ServePage(w, r, "profile-view", service.Service{"OtherUser": otherUser})
}

func handleProfile() {
        server.HandlePage("/profile", profilePageHandler)
        server.HandlePage("/profile/", viewProfilePageHandler)
}
