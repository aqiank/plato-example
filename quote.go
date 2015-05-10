package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"
)

type Quote struct {
	Author string `json:"author"`
	Text   string `json:"text"`
}

var quotes []Quote

func loadQuotes() {
	var qs [][2]string

	file, err := os.Open("quotes.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.NewDecoder(file).Decode(&qs)
	if err != nil {
		log.Fatal(err)
	}

	random := func(length int) []int {
		var indices []int

		in := func(numbers []int, n int) bool {
			for _, v := range numbers {
				if v == n {
					return true
				}
			}
			return false
		}

		rand.Seed(int64(time.Now().Second()))
		for i := 0; i < length; i++ {
			r := rand.Intn(length)
			if in(indices, r) {
				i--
				continue
			}
			indices = append(indices, r)
		}

		return indices
	}

	rns := random(len(qs))

	for _, v := range rns {
		quote := Quote{Author: qs[v][0], Text: qs[v][1]}
		quotes = append(quotes, quote)
	}
}
