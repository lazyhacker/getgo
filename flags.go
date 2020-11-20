package main

import (
	"flag"
	"runtime"
)

var (
	dl      = flag.String("dir", "", "Directory path to download to.")
	version = *flag.String("version", "", "Specific version to download (e.g. 1.14.7)")
	show    = flag.Bool("show", true, "If true, print out the file downloaded.")
	kind    = flag.String("kind", "archive", "What kind of file to download (archive, installer).")
	goos    = runtime.GOOS
	arch    = runtime.GOARCH
)
