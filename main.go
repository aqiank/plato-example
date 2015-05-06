package main

import (
        "log"
        "net/http"
        "os"
        "os/signal"

        "plato/debug"
        "plato/server"
        "plato/server/service"
)

func indexPageHandler(w http.ResponseWriter, r *http.Request) error {
        return server.ServePage(w, r, "index")
}

func projectPageHandler(w http.ResponseWriter, r *http.Request) error {
        return server.ServePage(w, r, "project")
}

func onSignUpSuccessHandler(w http.ResponseWriter, r *http.Request, data interface{}) error {
        switch data.(type) {
        case int64:
                debug.Log("Successfully signed up!")
        }

        return nil
}

func onPostNewSuccessHandler(w http.ResponseWriter, r *http.Request, data interface{}) error {
        http.Redirect(w, r, "/", 302)
        return nil
}

func main() {
        loadQuotes()
        loadCountries()

        // Demonstrate page handler
        server.HandlePage("/", indexPageHandler)
        server.HandlePage("/project", projectPageHandler)

        // Demonstrate API callback
        server.SetSuccessCallback("/signup", onSignUpSuccessHandler)
        server.SetSuccessCallback("/post/new", onPostNewSuccessHandler)

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
