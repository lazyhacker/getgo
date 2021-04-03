// +build !gui

package main // import "lazyhacker.dev/getgo"

import (
	"flag"
	"fmt"
	"log"
	"runtime"

	"lazyhacker.dev/getgo/internal/lib"
)

var (
	dl      = flag.String("dir", "", "Directory path to download to.")
	version = *flag.String("version", "", "Specific version to download (e.g. 1.14.7)")
	show    = flag.Bool("show", true, "If true, print out the file downloaded.")
	kind    = flag.String("kind", "archive", "What kind of file to download (archive, installer).")
	extract = flag.String("x", "", "Extract the archive to location specified.")
	plat    = flag.String("os", runtime.GOOS, "Override the OS type to download (e.g. linux, darwin, windows)")
	arch    = flag.String("arch", runtime.GOARCH, "Override the system architecture (e.g. amd64, arm)")
)

func main() {

	flag.Parse()

	// For ARM architecture, use v6l for Raspberry Pi.
	if *arch == "arm" {
		*arch = "armv6l"
	}

	stable, checksum, err := lib.LatestVersion(*plat, *arch, *kind)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = lib.DownloadAndVerify(*dl, stable, checksum, *extract)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if *show {
		fmt.Printf("%v\n", stable)
	}
}
