package main

import (
        "net/http"

        "plato/server"
)

func indexPageHandler(w http.ResponseWriter, r *http.Request) error {
        return server.ServePage(w, r, "index")
}

func main() {
        // Demonstrate page handler
        server.HandlePage("/", indexPageHandler)

        // Demonstrate files handler
        server.HandleFiles("/img/")

        server.Run()
}
