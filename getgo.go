// getgo is a command line tool to download the latest stable version of Go
// (http://golang.org) that matches the OS and architecture that it is executed
// from.  It will check the sha256 checksum to make sure the downloaded file is
// verified or delete it if it doesn't.
package main // import "lazyhacker.dev/getgo"

import (
	"crypto/sha256"
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
	STABLE_VERSION  = "https://golang.org/VERSION?m=text"
	GO_DOWNLOAD_URL = "https://dl.google.com/go"
)

var (
	sha_extension = ".sha256"
	dl            = flag.String("dir", "", "Directory path to download to.")
	version       = *flag.String("version", "", "Specific version to download (e.g. 1.14.7)")
	show          = flag.Bool("show", true, "If true, print out the file downloaded.")
)

func main() {

	flag.Parse()
	// Get the OS and architecture
	goos := runtime.GOOS
	arch := runtime.GOARCH

	resp, err := http.Get(STABLE_VERSION)
	if version == "" {
		if err != nil {
			log.Fatalf("Unable to get the latest version number. %v", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Unable to read the version number. %v", err)
		}
		version = string(body)
	}

	if arch == "arm" {
		arch = "armv6l"
	}
	var gofile string

	// Construct the file name for the stable binary.
	switch goos {
	case "linux":
		gofile = fmt.Sprintf("%v.%v-%v.tar.gz", version, goos, arch)
	case "freebsd":
		gofile = fmt.Sprintf("%v.%v-%v.tar.gz", version, goos, arch)
	case "windows":
		gofile = fmt.Sprintf("%v.%v-%v.msi", version, goos, arch)
	case "darwin":
		gofile = fmt.Sprintf("%v.%v-%v.pkg", version, goos, arch)
	default:
		log.Fatalln("Unknown OS... can't download.")

	}

	// Get the checksum value from the checksum file.
	resp, err = http.Get(fmt.Sprintf("%v/%v", GO_DOWNLOAD_URL, gofile+sha_extension))
	if err != nil {
		log.Fatalf("Unable to download sha256 file for %v. %v", gofile, err)
	}
	c, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read the checksum from file. %v", err)
	}
	sha256content := string(c)

	// Check if the binary has already been downloaded.
	var filepath string
	if *dl != "" {
		filepath = *dl + string(os.PathSeparator) + gofile
	} else {
		filepath = gofile
	}
	if _, err := os.Stat(filepath); err == nil {
		if m, _ := checksumMatch(filepath, sha256content); m {
			log.Println("Existing file is the latest stable and checksum verified.  Skipping download.")
			return
		}
	}

	// Download the golang binary
	download := fmt.Sprintf("%v/%v", GO_DOWNLOAD_URL, gofile)
	out, err := os.Create(filepath)
	if err != nil {
		log.Fatalf("Unable to create %v locally. %v", gofile, err)
	}
	defer out.Close()

	resp, err = http.Get(download)
	if err != nil {
		log.Fatalf("Unable to get the latest version number. %v", err)
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("Unable to download %v. %v", gofile, err)
	}

	// Compute the checksum
	if m, v := checksumMatch(filepath, sha256content); !m {
		log.Printf("Calcuated checksum %v != %v. Removing download.\n", v, sha256content)
		os.Remove(filepath)
	}

	if *show {
		fmt.Printf("%v", filepath)
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
