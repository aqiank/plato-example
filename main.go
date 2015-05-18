package main

import (
        "log"
        "net/http"
        "os"
        "os/signal"

        "plato/server"
        "plato/server/service"
)

func indexPageHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
        return nil, server.ServePage(w, r, "index", nil)
}

func main() {
        loadQuotes()
        loadCountries()
        loadProfessions()

	handleProject()
	handleDashboard()
	handleProfile()
	handleSearch()

        // Demonstrate page handler
        server.HandlePage("/", indexPageHandler)

        // Demonstrate files handler
        server.HandleFiles("/css/", "/font/", "/img/", "/js/", "/lib/", "/pt-data/")

        // Demonstrate service
        service.AttachAll(service.Service{
		"Quotes": quotes,
		"Countries": countries,
		"Professions": professions,
		"RecommendedProjects": recommendedProjects,
		"LatestRelatedProjects": latestRelatedProjects,
		"GetProject": getProject,
		"GetApplicants": getApplicants,
	})

        log.Println("starting plato-example..")

        go func() {
                server.Run()
        }()

        sigc := make(chan os.Signal, 1)
        signal.Notify(sigc, os.Interrupt, os.Kill)
        <-sigc

        log.Println("stopping plato-example..")
}
