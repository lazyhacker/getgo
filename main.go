//go:build !gui
// +build !gui

package main // import "lazyhacker.dev/getgo"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"

	"lazyhacker.dev/getgo/internal/download"
)

var (
	dl         = flag.String("dir", "", "Directory path to download to.")
	version    = *flag.String("version", "", "Specific version to download (e.g. 1.14.7)")
	show       = flag.Bool("show", true, "If true, print out the file downloaded.")
	kind       = flag.String("kind", "archive", "What kind of file to download (archive, installer).")
	extract    = flag.String("x", "", "Extract the archive to location specified.")
	plat       = flag.String("os", runtime.GOOS, "Override the OS type to download (e.g. linux, darwin, windows)")
	arch       = flag.String("arch", runtime.GOARCH, "Override the system architecture (e.g. amd64, arm)")
	build_info = flag.Bool("build", false, "Print the build info.")

	build      = "dev" // default
	build_date = "unknown"
	commit     = "none"
)

func main() {

	flag.Parse()

	if *build_info {
		fmt.Printf("Version:%s\nDate:%s\nCommit:%s\n", build, build_date, commit)
		os.Exit(0)
	}
	// For ARM architecture, use v6l for Raspberry Pi.
	if *arch == "arm" {
		*arch = "armv6l"
	}

	stable, checksum, err := download.LatestVersion(*plat, *arch, *kind)
	if err != nil {
		log.Fatalf("%v", err)
	}

	f, err := download.DownloadAndVerify(*dl, stable, checksum, *extract)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if *show {
		fmt.Printf("%v\n", f)
	}
}
