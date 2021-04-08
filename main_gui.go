// +build gui

package main // import "lazyhacker.dev/getgo"

import (
	"flag"

	"lazyhacker.dev/getgo/internal/guimain"
)

func main() {

	flag.Parse()

	guimain.LoadGUI()
}
