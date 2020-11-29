// +build gui,gtk3

package guimain

import (
	"fmt"
	"log"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"lazyhacker.dev/getgo/internal/lib"
)

func LoadGUI(filename, checksum string) {
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("GetGo")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	wd, _ := os.Getwd()
	dirValue, err := gtk.LabelNew(wd)
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	fileValue, err := gtk.LabelNew(filename)
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	shaValue, err := gtk.LabelNew(checksum)
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}
	downloadBtn, err := gtk.ButtonNewWithLabel("Download Latest")
	if err != nil {
		log.Fatal("Unable to create button:", err)
	}
	downloadProgress, err := gtk.ProgressBarNew()
	defer downloadProgress.Destroy()
	if err != nil {
		log.Printf("Unable to load progress bar. %v", err)
	}

	downloadBtn.Connect("clicked", func() {

		// Set a timer function that periodically updates the progress bar.
		h, err := glib.TimeoutAdd(1000, updateProgress, downloadProgress)
		if err != nil {
			log.Printf("Error with timeline: %v", err)
		}

		// Have a go routine handle the download so the UI doesn't get blocked
		// waiting for the download to complete before the showing the progress
		// bar.
		done := make(chan int)
		go func() {
			if err := lib.DownloadAndVerify(wd, filename, checksum); err != nil {
				errmsg := gtk.MessageDialogNew(win, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE, fmt.Sprintf("Error download. %v", err))
				defer errmsg.Destroy()
				errmsg.Run()
			}
			done <- 1
		}()

	FS:
		// Listen for the goroutine to tell us that the download is completed.
		// While waiting, don't block the UI so pass control back to GTK.
		for {
			select {
			case <-done:
				break FS
			default:
				if gtk.EventsPending() {
					gtk.MainIteration()
				}
			}
		}

		glib.SourceRemove(h)
		dialog := gtk.MessageDialogNew(win, gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_CLOSE, "Complete!")
		defer dialog.Destroy()
		dialog.Run()
	})

	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}

	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	grid.Add(dirValue)
	grid.Add(fileValue)
	grid.Add(shaValue)
	grid.Add(downloadProgress)
	grid.Add(downloadBtn)

	win.Add(grid)

	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}

func updateProgress(pb *gtk.ProgressBar) bool {
	pb.Pulse()
	return true
}
