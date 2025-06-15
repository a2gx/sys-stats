package main

import "fmt"

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

func getVersion() string {
	data := fmt.Sprintf("BuildDate: %s, GitHash: %s", buildDate, gitHash)
	return fmt.Sprintf("%s\n%s", release, data)
}
