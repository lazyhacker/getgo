// +build !gui

// getgo is a command line tool to download the latest stable version of Go
// (http://golang.org) that matches the OS and architecture that it is executed
// from.  It will check the sha256 checksum to make sure the downloaded file is
// verified or delete it if it doesn't.
package main // import "lazyhacker.dev/getgo"

import (
	"flag"
	"fmt"
	"log"

	"lazyhacker.dev/getgo/internal/lib"
)

var (
	dl      = flag.String("dir", "", "Directory path to download to.")
	version = *flag.String("version", "", "Specific version to download (e.g. 1.14.7)")
	show    = flag.Bool("show", true, "If true, print out the file downloaded.")
	kind    = flag.String("kind", "archive", "What kind of file to download (archive, installer).")
)

func main() {

	flag.Parse()

	stable, checksum, err := lib.LatestVersion(*kind)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = lib.DownloadAndVerify(*dl, stable, checksum)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if *show {
		fmt.Printf("%v\n", stable)
	}
}
