A simple program to download the latest stable version of Go
(http:golang.org) and verifies that the downloaded file matches the checksum.

It is primarily an command line tool with an experimental GUI mode that is only
half-baked.

It is written to as a convenience tool so that you don't have to go to the web
site, find and download the correct binary for your platform, download the
checksum and then verify it yourself.

The tool won't re-download the binary if a verified one already exists locally.

The tool will not install the binary.  It only downloads it to the directory of
your choice (or current directory if unspecified).

## Install

### If a version of Go is already installed:

`go get lazyhacker.dev/getgo` (command line only)

`go get lazyhacker.dev/getgo -tags gui` (GUI version but requires Fyne to be
installed)

### If no version of Go is available locally:

Pre-built binaries for Linux is available in the releases section in case there
is no version of Go already installed.

## Usage

To get help info:
`getgo -help`

To download to a specific directory:
`getgo -dir ~/Downloads`

NOTE: getgo checks https://golang.org/dl/?mode=json to determine the latest
stable version of Go.
