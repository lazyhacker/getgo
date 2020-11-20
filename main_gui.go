// +build gui

// getgo is a command line tool to download the latest stable version of Go
// (http://golang.org) that matches the OS and architecture that it is executed
// from.  It will check the sha256 checksum to make sure the downloaded file is
// verified or delete it if it doesn't.
package main // import "lazyhacker.dev/getgo"

import (
	"flag"
	"fmt"
	"log"
)

var (
	win = flag.Bool("gui", false, "Run with a GUI.")
)

func main() {

	flag.Parse()

	stable, checksum, err := LatestVersion()
	if err != nil {
		log.Fatalf("%v", err)
	}

	if *win {
		LoadGUI(stable, checksum)
	} else {
		err := DownloadAndVerify(stable, checksum)
		if err != nil {
			log.Fatalf("%v", err)
		}

		if *show {
			fmt.Printf("%v\n", stable)
		}
	}
}
