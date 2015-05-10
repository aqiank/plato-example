package main

import (
	"errors"
	"net/http"

	"plato/server"
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

func handleProfile() {
        server.HandlePage("/profile", profilePageHandler)
}
