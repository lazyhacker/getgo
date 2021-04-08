package lib

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	STABLE_VERSION  = "https://golang.org/dl/?mode=json"
	GO_DOWNLOAD_URL = "https://golang.org/dl" // redirects to https://dl.google.com/go
)

var (
	OperatingSystems = []string{"aix", "darwin", "dragonfly", "freebsd", "illumos", "js", "linux", "netbsd", "openbsd", "plan9", "solaris", "windows"}
	Architectures    = []string{"386", "amd64", "arm", "arm64", "ppc64", "ppc64le", "mips", "mips64", "mips64le", "riscv64", "s390x", "wasm"}
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

// DownloadAndVerify downloads the Go binary that is passed in and verify
// it against its sha256 checksum.  If there is already a local file
// that matches the checksum then it will not download another version.
func DownloadAndVerify(destdir, filename, checksum, extractDir string) error {

	var filepath, calcSum string
	var sumMatch bool

	if destdir != "" {
		filepath = destdir + string(os.PathSeparator) + filename
	} else {
		filepath = filename
	}

	// Check if the binary has already been downloaded.
	if _, err := os.Stat(filepath); err == nil {
		sumMatch, calcSum = checksumMatch(filepath, checksum)
		if sumMatch {
			log.Println("Existing file is the latest stable and checksum verified.  Skipping download.")
			extractArchive(extractDir, filepath)
			return nil
		}
	}

	// Download the golang binary
	log.Println("Beginning download.")
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

	log.Println("Download complete!  Verifying checksum...")
	// Compute and verify the checksum of the downloaded file.
	sumMatch, calcSum = checksumMatch(filepath, checksum)
	if !sumMatch {
		os.Remove(filepath)
		return fmt.Errorf("Calculated checksum %v != %v. Removing download.\n", calcSum, checksum)
	}

	extractArchive(extractDir, filepath)
	return nil

}

func extractArchive(dest, archive string) error {

	if len(dest) > 0 {

		switch ext := filepath.Ext(archive); ext {
		case ".gz":
			log.Printf("Untar %v.", archive)
			err := untarArchive(archive, dest)
			if err != nil {
				return fmt.Errorf("unable to extract the archive. %v", err)
			}
		case ".zip":
			log.Printf("Unzipping %v", archive)
			err := unzipArchive(archive, dest)
			if err != nil {
				return fmt.Errorf("unable to extract the archive. %v", err)
			}
		default:
			log.Println("Only .zip and .gz supported for extracting archives.")
		}
	}

	return nil
}

// LatestVersion returns the highest stable versions for the platform.
func LatestVersion(goos, goarch, kind string) (filename string, sha256sum string, err error) {

	var gfs []goFilesStruct
	var max *goFilesStruct

	log.Printf("Checking for the latest stable version for %v-%v.", goos, goarch)
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

	log.Printf("Latest stable version is %v.\n", max.Version)

	for _, v := range max.Files {
		if v.Os == goos && v.Arch == goarch && v.Kind == kind {
			return v.Filename, v.Sha256, nil
		}
	}
	return "", "", fmt.Errorf("No download found for OS=%v ARCH=%v KIND=%v.", goos, goarch, kind)
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

// untarArchive will extract the downloaded archive to the path specified.
func untarArchive(archive, path string) error {

	// Open the archive and pass it to the gzip library to uncompress.
	f, err := os.Open(archive)
	if err != nil {
		return fmt.Errorf("unable to open archive. %v", err)
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return fmt.Errorf("unable to read archive. %v", err)
	}
	defer gz.Close()

	// Read the tarbar from the gzip archive.
	t := tar.NewReader(gz)
	for {

		hdr, err := t.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error when trying to untar the archive. %v", err)
		}
		dest := filepath.Join(path, hdr.Name)

		info := hdr.FileInfo()

		if info.IsDir() {
			if err := os.MkdirAll(dest, 0755); err != nil {
				return err
			}
			continue
		}

		if _, err := os.Stat(dest); !os.IsNotExist(err) {
			log.Println("Existing files at destination.  Aborting extraction...")
			return err
		}
		file, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return fmt.Errorf("unable to create file on disk. %v", err)
		}

		_, err = io.Copy(file, t)
		if err != nil {
			return fmt.Errorf("unable to write file to disk. %v", err)
		}
		file.Close()
		fmt.Printf("%s\n", hdr.Name)
	}

	return nil
}

func unzipArchive(archive, path string) error {

	r, err := zip.OpenReader(archive)
	if err != nil {
		return fmt.Errorf("unable to access archive %v.  %v", archive, err)
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	for _, f := range r.File {

		info := f.FileInfo()
		dest := filepath.Join(path, f.Name)

		if info.IsDir() {
			if err := os.MkdirAll(dest, 0755); err != nil {
				return err
			}
			continue
		}

		c, err := f.Open()
		if err != nil {
			return fmt.Errorf("Unable to open archive %v. %v", archive, err)
		}

		if _, err := os.Stat(dest); !os.IsNotExist(err) {
			log.Println("Existing files at destination.  Aborting extraction...")
			return err
		}

		file, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return fmt.Errorf("unable to create file on disk. %v", err)
		}

		_, err = io.Copy(file, c)
		if err != nil {
			return fmt.Errorf("unable to write file to disk. %v", err)
		}
		file.Close()
		fmt.Printf("%s\n", f.Name)
	}

	return nil

}
