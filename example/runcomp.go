package main

import (
	"log"
	"os"

	"github.com/chchench/gotd"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("2 parameters expected, one path for target test run JSON file and another for comp test run.")
	}
	gotd.CompareRuns(os.Args[1], os.Args[2])
	os.Exit(0)
}
