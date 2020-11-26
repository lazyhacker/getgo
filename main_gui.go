// +build gui

package main // import "lazyhacker.dev/getgo"

import (
	"flag"
	"fmt"
	"log"

	"lazyhacker.dev/getgo/internal/guimain"
	"lazyhacker.dev/getgo/internal/lib"
)

var (
	dl      = flag.String("dir", "", "Directory path to download to.")
	version = *flag.String("version", "", "Specific version to download (e.g. 1.14.7)")
	show    = flag.Bool("show", true, "If true, print out the file downloaded.")
	kind    = flag.String("kind", "archive", "What kind of file to download (archive, installer).")

	win = flag.Bool("gui", false, "Run with a GUI.")
)

func main() {

	flag.Parse()

	stable, checksum, err := lib.LatestVersion(*kind)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if *win {
		guimain.LoadGUI(stable, checksum)
	} else {
		err := lib.DownloadAndVerify(*dl, stable, checksum)
		if err != nil {
			log.Fatalf("%v", err)
		}

		if *show {
			fmt.Printf("%v\n", stable)
		}
	}
}
