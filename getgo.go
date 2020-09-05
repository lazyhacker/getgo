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
	GO_DOWNLOAD_URL = "https://golang.org/dl" // redicts to https://dl.google.com/go
)

var (
	sha_extension = ".sha256"
	dl            = flag.String("dir", "", "Directory path to download to.")
	version       = *flag.String("version", "", "Specific version to download (e.g. 1.14.7)")
	show          = flag.Bool("show", true, "If true, print out the file downloaded.")
	kind          = flag.String("kind", "archive", "What kind of file to download (archive, installer).")
)

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

func main() {

	flag.Parse()
	// Get the OS and architecture
	goos := runtime.GOOS
	arch := runtime.GOARCH

	var gfs []goFilesStruct
	var max *goFilesStruct

	if version == "" {

		resp, err := http.Get(STABLE_VERSION)
		if err != nil {
			log.Fatalf("Unable to get the latest version number. %v", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Unable to read the body of the response. %v", err)
		}
		jserr := json.Unmarshal(body, &gfs)
		if jserr != nil {
			log.Fatalf("Unable to unmarshal response body. %v", jserr)
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

	if max == nil {
		log.Fatal("Unable to find a stable version!")
	}

	if arch == "arm" {
		arch = "armv6l"
	}

	for i, v := range max.Files {
		if v.Os == goos && v.Arch == arch && v.Kind == *kind {
			if *show {
				fmt.Printf("%v\n", v.Filename)
			}
			// Check if the binary has already been downloaded.
			var filepath string
			if *dl != "" {
				filepath = *dl + string(os.PathSeparator) + v.Filename
			} else {
				filepath = v.Filename
			}
			if _, err := os.Stat(filepath); err == nil {
				if m, _ := checksumMatch(filepath, v.Sha256); m {
					log.Println("Existing file is the latest stable and checksum verified.  Skipping download.")
					return
				}
			}

			// Download the golang binary
			download := fmt.Sprintf("%v/%v", GO_DOWNLOAD_URL, v.Filename)
			out, err := os.Create(filepath)
			if err != nil {
				log.Fatalf("Unable to create %v locally. %v", filepath, err)
			}
			defer out.Close()

			resp, err := http.Get(download)
			if err != nil {
				log.Fatalf("Unable to fetch the binary %v. %v", download, err)
			}
			_, err = io.Copy(out, resp.Body)
			if err != nil {
				log.Fatalf("Unable to write %v. %v", filepath, err)
			}

			// Compute the checksum
			if m, sum := checksumMatch(filepath, v.Sha256); !m {
				log.Printf("Calculated checksum %v != %v. Removing download.\n", sum, v.Sha256)
				os.Remove(filepath)
			}
			break
		}
		if i == len(max.Files)-1 {
			log.Fatal("Unable to find any files to download.")
		}
	}
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
