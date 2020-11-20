package main

import (
	"flag"
	"runtime"
)

const (
	STABLE_VERSION  = "https://golang.org/dl/?mode=json"
	GO_DOWNLOAD_URL = "https://golang.org/dl" // redirects to https://dl.google.com/go
)

var (
	dl      = flag.String("dir", "", "Directory path to download to.")
	version = *flag.String("version", "", "Specific version to download (e.g. 1.14.7)")
	show    = flag.Bool("show", true, "If true, print out the file downloaded.")
	kind    = flag.String("kind", "archive", "What kind of file to download (archive, installer).")
	goos    = runtime.GOOS
	arch    = runtime.GOARCH
)

func init() {

	// For ARM architecture, use v6l for Raspberry Pi.
	if arch == "arm" {
		arch = "armv6l"
	}

}
