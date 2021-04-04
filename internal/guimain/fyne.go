// +build gui,fyne

package guimain

import (
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"

	"lazyhacker.dev/getgo/internal/lib"
)

func LoadGUI(filename, checksum string) {
	a := app.New()
	w := a.NewWindow("GetGo")

	wd, _ := os.Getwd()
	dirLabel := canvas.NewText("Directory", colornames.Gray)
	dirValue := canvas.NewText(wd, color.White)
	fileLabel := canvas.NewText("Latest Stable", colornames.Gray)
	fileValue := canvas.NewText(filename, color.White)
	shaLabel := canvas.NewText("Sha256", colornames.Gray)
	shaValue := canvas.NewText(checksum, color.White)

	formGrid := fyne.NewContainerWithLayout(
		layout.NewFormLayout(),
		dirLabel, dirValue,
		fileLabel, fileValue,
		shaLabel, shaValue)

	grid := fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		formGrid,
		widget.NewButton("Download Latest",
			func() {
				prog := dialog.NewProgressInfinite("Downloading", fileValue.Text, w)
				prog.Show()
				err := lib.DownloadAndVerify(dirValue.Text, filename, checksum, "")
				prog.Hide()
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				dialog.ShowInformation("Complete!", "File downloaded.", w)
			},
		),
	)
	w.SetContent(grid)
	w.ShowAndRun()
}
