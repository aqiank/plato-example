package main

import (
	"net/http"

	"plato/server"
	"plato/server/session"
)

func dashboardPageHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	user := session.User(r)
	if !session.IsLoggedIn(user) {
		http.Redirect(w, r, "/", 302)
		return nil, nil
	}

	if r.FormValue("what") == "" {
		return nil, server.ServePage(w, r, "dashboard", nil)
	}

	return nil, nil
}

func handleDashboard() {
        server.HandlePage("/dashboard", dashboardPageHandler)
}
