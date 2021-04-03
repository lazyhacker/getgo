// +build gui

package main // import "lazyhacker.dev/getgo"

import (
	"flag"
	"log"

	"lazyhacker.dev/getgo/internal/guimain"
	"lazyhacker.dev/getgo/internal/lib"
)

func main() {

	flag.Parse()

	stable, checksum, err := lib.LatestVersion(*kind)
	if err != nil {
		log.Fatalf("%v", err)
	}

	guimain.LoadGUI(stable, checksum)
}
