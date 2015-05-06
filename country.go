package main

import (
       "bufio"
       "os"
       "strings"

       "plato/debug"
)

var countries = make(map[string]string)

func loadCountries() {
        file, err := os.Open("countries.txt")
        if err != nil {
                debug.Warn(err)
                return
        }

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                ss := strings.Split(scanner.Text(), ":")
                countries[ss[0]] = ss[1]
        }

        if err = scanner.Err(); err != nil {
               debug.Warn(err)
               return
        }

        debug.Log("loaded", len(countries), "countries")
}
