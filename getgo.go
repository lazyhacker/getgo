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
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
)

const (
	STABLE_VERSION  = "https://golang.org/dl/?mode=json"
	GO_DOWNLOAD_URL = "https://golang.org/dl" // redirects to https://dl.google.com/go
)

var (
	sha_extension = ".sha256"
	dl            = flag.String("dir", "", "Directory path to download to.")
	version       = *flag.String("version", "", "Specific version to download (e.g. 1.14.7)")
	show          = flag.Bool("show", true, "If true, print out the file downloaded.")
	kind          = flag.String("kind", "archive", "What kind of file to download (archive, installer).")
	goos          = runtime.GOOS
	arch          = runtime.GOARCH
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

	stable, err := latestVersion()
	if err != nil {
		log.Fatalf("%v", err)
	}

	f, _, err := downloadAndVerify(stable)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if *show {
		fmt.Printf("%v\n", f)
	}
}

// downloadAndVerify downloads the Go binary that is passed in and verify
// it against its sha256 checksum.  If there is already a local file
// that matches the checksum then it will not download another version.
func downloadAndVerify(gfs *goFilesStruct) (filename string, sum string, err error) {

	var filepath string
	var sumMatch bool

	for i, v := range gfs.Files {
		if v.Os == goos && v.Arch == arch && v.Kind == *kind {
			// Check if the binary has already been downloaded.
			if *dl != "" {
				filepath = *dl + string(os.PathSeparator) + v.Filename
			} else {
				filepath = v.Filename
			}

			if _, err := os.Stat(filepath); err == nil {
				sumMatch, sum = checksumMatch(filepath, v.Sha256)
				if sumMatch {
					log.Println("Existing file is the latest stable and checksum verified.  Skipping download.")
					return v.Filename, sum, nil
				}
			}

			// Download the golang binary
			download := fmt.Sprintf("%v/%v", GO_DOWNLOAD_URL, v.Filename)
			out, err := os.Create(filepath)
			if err != nil {
				return "", "", fmt.Errorf("unable to create %v locally. %v", filepath, err)
			}
			defer out.Close()

			resp, err := http.Get(download)
			if err != nil {
				return "", "", fmt.Errorf("unable to fetch the binary %v. %v", download, err)
			}
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				return "", "", fmt.Errorf("unable to write %v. %v", filepath, err)
			}

			// Compute the checksum
			if !sumMatch {
				log.Printf("Calculated checksum %v != %v. Removing download.\n", sum, v.Sha256)
				os.Remove(filepath)
			}
			return v.Filename, sum, nil
		}
		if i == len(gfs.Files)-1 {
			return "", "", fmt.Errorf("Unable to find any files to download.")
		}
	}

	return "", "", nil
}

// latestVersion returns the highest named version that is marked stable.
func latestVersion() (*goFilesStruct, error) {

	var gfs []goFilesStruct
	var max *goFilesStruct

	if version == "" {
		resp, err := http.Get(STABLE_VERSION)
		if err != nil {
			return nil, fmt.Errorf("unable to get the latest version number. %v", err)
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to read the body of the response. %v", err)
		}

		jserr := json.Unmarshal(body, &gfs)
		if jserr != nil {
			return nil, fmt.Errorf("unable to unmarshal response body. %v", jserr)
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

	return max, nil
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
