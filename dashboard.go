package main

import (
	"net/http"

	"plato/fileutil"
	"plato/server/page"
	"plato/server/session"
)

func dashboardPageHandler(w http.ResponseWriter, r *http.Request) error {
	user := session.User(r)
	if !session.IsLoggedIn(user) {
		http.Redirect(w, r, "/", 302)
		return nil
	}

	filepath := timelinePath(user)
	if !fileutil.Exists(filepath) {
		generateTimeline(user)
	}

	if r.FormValue("what") == "" {
		return page.Serve(w, r, "dashboard", nil)
	}

	return nil
}

func handleDashboard() {
	page.Handle("/dashboard", dashboardPageHandler)
}
