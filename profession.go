package main

import (
	"bufio"
	"os"

	"plato/debug"
)

var professions []string

func loadProfessions() {
	file, err := os.Open("professions.txt")
	if err != nil {
		debug.Warn(err)
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		professions = append(professions, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		debug.Warn(err)
		return
	}

	debug.Log("loaded", len(professions), "professions")
}
