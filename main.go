package main

import (
        "log"
        "net/http"
        "os"
        "os/signal"

        "plato/db/dateutil"
        "plato/debug"
        "plato/server"
        "plato/server/service"
)

func indexPageHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
        return nil, server.ServePage(w, r, "index")
}

func newProjectPageHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
        return nil, server.ServePage(w, r, "project-new")
}

func onSignUpSuccessHandler(w http.ResponseWriter, r *http.Request, data interface{}) error {
        http.Redirect(w, r, "/", 302)
        return nil
}

func newProjectHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
        data, err := server.PostHandler(w, r)
        if err != nil {
                return nil, debug.Error(err)
        }

        imageURL, err := saveProjectImage(r)
        if err != nil {
                return nil, debug.Error(err)
        }

        tp := dateutil.TimeParser{}
        startDateStr := r.FormValue("start-date")
        endDateStr := r.FormValue("end-date")
        startDate := tp.ParseDatetime(startDateStr)
        endDate := tp.ParseDatetime(endDateStr)
        if tp.Err != nil {
                return nil, debug.Error(err)
        }

        postID := data.(int64)
        status := r.FormValue("status")

        id, err := InsertProject(postID, status, imageURL, startDate, endDate)
        if err != nil {
                return nil, debug.Error(err)
        }

        return id, nil
}

func main() {
        loadQuotes()
        loadCountries()

        // Demonstrate page handler
        server.HandlePage("/", indexPageHandler)
        server.HandlePage("/project", newProjectPageHandler)

        // Demonstrate API callback
        server.SetSuccessCallback("/signup", onSignUpSuccessHandler)

        // Demonstrate files handler
        server.HandleFiles("/css/", "/font/", "/img/", "/js/")

        // Demonstrate service
        service.AttachAll(service.Service{"Quotes": quotes})
        service.AttachAll(service.Service{"Countries": countries})

        log.Println("starting plato-example..")

        go func() {
                server.Run()
        }()

        sigc := make(chan os.Signal, 1)
        signal.Notify(sigc, os.Interrupt, os.Kill)
        <-sigc

        log.Println("stopping plato-example..")
}
