// getgo is a command line tool to download the latest stable version of Go
// (http://golang.org) that matches the OS and architecture that it is executed
// from.  It will check the sha256 checksum to make sure the downloaded file is
// verified or delete it if it doesn't.
package main // import "lazyhacker.dev/getgo"

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"image/color"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"golang.org/x/image/colornames"
)

var (
	dl      = flag.String("dir", "", "Directory path to download to.")
	version = *flag.String("version", "", "Specific version to download (e.g. 1.14.7)")
	show    = flag.Bool("show", true, "If true, print out the file downloaded.")
	kind    = flag.String("kind", "archive", "What kind of file to download (archive, installer).")
	win     = flag.Bool("gui", false, "Run with a GUI.")
	goos    = runtime.GOOS
	arch    = runtime.GOARCH
)

const (
	STABLE_VERSION  = "https://golang.org/dl/?mode=json"
	GO_DOWNLOAD_URL = "https://golang.org/dl" // redirects to https://dl.google.com/go
)

// goFilesStruct maps to the JSON format from STABLE_VERSION.
type goFilesStruct struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []struct {
		Filename string `json:"filename"`
		Os       string `json:"os"`
		Arch     string `json:"arch"`
		Version  string `json:"version"`
		Sha256   string `json:"sha256"`
		Size     int    `json:"size"`
		Kind     string `json:"kind"`
	} `json:"files"`
}

func init() {

	// For ARM architecture, use v6l for Raspberry Pi.init
	if arch == "arm" {
		arch = "armv6l"
	}

}

func main() {

	flag.Parse()

	stable, checksum, err := LatestVersion()
	if err != nil {
		log.Fatalf("%v", err)
	}

	if *win {
		LoadGUI(stable, checksum)
	} else {
		err := DownloadAndVerify(stable, checksum)
		if err != nil {
			log.Fatalf("%v", err)
		}

		if *show {
			fmt.Printf("%v\n", stable)
		}
	}
}

// DownloadAndVerify downloads the Go binary that is passed in and verify
// it against its sha256 checksum.  If there is already a local file
// that matches the checksum then it will not download another version.
func DownloadAndVerify(filename, checksum string) error {

	var filepath, calcSum string
	var sumMatch bool

	if *dl != "" {
		filepath = *dl + string(os.PathSeparator) + filename
	} else {
		filepath = filename
	}

	// Check if the binary has already been downloaded.
	if _, err := os.Stat(filepath); err == nil {
		sumMatch, calcSum = checksumMatch(filepath, checksum)
		if sumMatch {
			log.Println("Existing file is the latest stable and checksum verified.  Skipping download.")
			return nil
		}
	}

	// Download the golang binary
	download := fmt.Sprintf("%v/%v", GO_DOWNLOAD_URL, filename)
	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("unable to create %v locally. %v", filepath, err)
	}
	defer out.Close()
	resp, err := http.Get(download)
	if err != nil {
		return fmt.Errorf("unable to fetch the binary %v. %v", download, err)
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("unable to write %v. %v", filepath, err)
	}

	// Compute and verify the checksum of the downloaded file.
	sumMatch, calcSum = checksumMatch(filepath, checksum)
	if !sumMatch {
		os.Remove(filepath)
		return fmt.Errorf("Calculated checksum %v != %v. Removing download.\n", calcSum, checksum)
	}

	return nil

}

// LatestVersion returns the highest stable versions for the platform.
func LatestVersion() (filename string, sha256sum string, err error) {

	var gfs []goFilesStruct
	var max *goFilesStruct

	if version == "" {
		resp, err := http.Get(STABLE_VERSION)
		if err != nil {
			return "", "", fmt.Errorf("unable to get the latest version number. %v", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", "", fmt.Errorf("unable to read the body of the response. %v", err)
		}

		jserr := json.Unmarshal(body, &gfs)
		if jserr != nil {
			return "", "", fmt.Errorf("unable to unmarshal response body. %v", jserr)
		}

		for i, v := range gfs {
			if !v.Stable {
				continue
			}
			if max == nil {
				max = &gfs[i]
			}
			if v.Version > max.Version {
				max = &gfs[i]
			}
		}
	}

	log.Printf("Latest stable version is %v.\n", max.Version)

	for _, v := range max.Files {
		if v.Os == goos && v.Arch == arch && v.Kind == *kind {
			return v.Filename, v.Sha256, nil
		}
	}
	return "", "", fmt.Errorf("No download found for OS=%v ARCH=%v KIND=%v.", goos, arch, *kind)
}

func checksumMatch(f, v string) (bool, string) {

	hash := sha256.New()
	gf, err := os.Open(f)
	defer gf.Close()
	if _, err = io.Copy(hash, gf); err != nil {
		log.Printf("Unable to compute sha256 checksum. %v", err)
	}
	sha256sum := fmt.Sprintf("%x", hash.Sum(nil))

	if sha256sum != v {
		return false, sha256sum
	}

	return true, sha256sum
}

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
				err := DownloadAndVerify(filename, checksum)
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
