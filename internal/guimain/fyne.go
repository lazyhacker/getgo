// +build gui

package guimain

import (
	"os"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"lazyhacker.dev/getgo/internal/lib"
)

/* GetGo describes the elements of the main UI and the data needed by the UI.*/
type GetGo struct {
	parent          fyne.Window // reference the parent window for the app.
	operatingSystem string      // GOOS version
	architecture    string      // GOARCH (e.g. amd64, arm, etc.)
	kind            string      // archive or installer
	savepath        string      // directory to save to
	filename        string      // Go binary to download
	shasum          string      // SHA256 of the binary

	formInput *widget.Form // Input form for OS, Arch and Kind

	// The read-only file info.  Using widget.Form as a cheat for the layout.
	formFileInfo *widget.Form

	lblChecksumValue *widget.Label // Label to display the SHA256 value.
	lblSavePath      *widget.Label // Label to show the directory to save to.
	lblFileName      *widget.Label // Label to show the binary filename to download

	iconDir *widget.Icon

	btnDownload *widget.Button // button to start the download
	btnSaveFile *widget.Button // button to pick the save directory

	saveDirWidget fyne.CanvasObject // container for the save button and label.
}

/* downloadGo starts the downloading and the corresponding dialogs. */
func (a *GetGo) downloadGo() {
	a.filename, a.shasum, _ = lib.LatestVersion(a.operatingSystem, a.architecture, a.kind)
	a.lblFileName.SetText(a.filename)
	a.lblChecksumValue.SetText(a.shasum)
	prog := dialog.NewProgressInfinite("Downloading", "", a.parent)
	prog.Show()
	err := lib.DownloadAndVerify(a.savepath, a.filename, a.shasum, "")
	prog.Hide()
	if err != nil {
		dialog.ShowError(err, a.parent)
		return
	}
	dialog.ShowInformation("Complete!", "File downloaded.", a.parent)

}

/*
Init should be called first to initialize the UI elements and default values.
It can be called again to reset the UI elements back to the defaults.
*/
func (a *GetGo) Init(w fyne.Window) {

	defaultkind := "archive"

	pwd, _ := os.Getwd()

	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		defaultkind = "installer"
	}

	a.parent = w
	a.operatingSystem = runtime.GOOS
	a.architecture = runtime.GOARCH
	a.savepath = pwd
	a.kind = defaultkind
	a.lblChecksumValue = widget.NewLabel("")
	a.lblFileName = widget.NewLabel("")

	// Select boxes for OS, arch and kind
	selOS := widget.NewSelect(lib.OperatingSystems,
		func(s string) {
			a.operatingSystem = s
		},
	)
	selArch := widget.NewSelect(lib.Architectures, func(s string) {
		a.architecture = s
	},
	)
	selKind := widget.NewSelect(
		[]string{"archive", "installer"},
		func(s string) {
			a.kind = s
		},
	)

	// Add the selection dropdowns into a single form.
	a.formInput = widget.NewForm(
		widget.NewFormItem("OS", selOS),
		widget.NewFormItem("Arch", selArch),
		widget.NewFormItem("Kind", selKind),
	)

	// Set the default selection.
	selOS.SetSelected(a.operatingSystem)
	selArch.SetSelected(a.architecture)
	selKind.SetSelected(a.kind)

	a.btnSaveFile = widget.NewButtonWithIcon("", theme.FolderIcon(),
		func() {
			dialog.ShowFolderOpen(
				func(uri fyne.ListableURI, err error) {
					if uri == nil || err != nil {
						return
					}
					// TODO: Should look into data binding for these fields.
					a.savepath = uri.Path()
					a.lblSavePath.SetText(uri.Path())
				}, w)
		},
	)

	// Default to the directory the app was launched.
	a.lblSavePath = widget.NewLabel(pwd)

	a.formFileInfo = widget.NewForm(
		widget.NewFormItem("File", a.lblFileName),
		widget.NewFormItem("Sha256", a.lblChecksumValue),
	)

	a.btnDownload = widget.NewButton("Download",
		func() {
			a.downloadGo()
		})

	// Want two column side-by-side layout similar to a form but widget.Form
	// takes a string and not an icon for the input label.
	a.saveDirWidget = container.New(layout.NewHBoxLayout(),
		a.btnSaveFile,
		a.lblSavePath,
	)
}

/*
LoadGUI is the main entry point to start an Fyne app.  It will initialize an
app, create a window, set the container and start the app.
*/
func LoadGUI() {
	a := app.New()
	w := a.NewWindow("GetGo")
	w.Resize(fyne.Size{Width: 680, Height: 280})

	ui := GetGo{}
	ui.Init(w)
	c := container.New(layout.NewVBoxLayout(),
		ui.formInput,
		ui.saveDirWidget,
		ui.btnDownload,
		ui.formFileInfo,
	)

	w.SetContent(c)
	w.ShowAndRun()
}
