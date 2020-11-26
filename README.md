A simple program to download the latest stable version of [Go](https://golang.org)
and verifies that it matches the checksum.

It is primarily a command line tool with an experimental GUI mode that is only
half-baked.

It is written to as a convenience tool so that you don't have to go to the web
site, find and download the correct binary for your platform, download the
checksum and then verify it yourself.

The tool won't re-download the binary if a verified one already exists locally.

The tool will not install the binary.  It only downloads it to the directory of
your choice (or current directory if unspecified).

## Install

### Pre-built Binaries

Pre-built binaries are available in the releases section.

### Compile From Source

To compile from the source, use the normal Go command:

`go get lazyhacker.dev/getgo`



## Usage

To get help info:

`getgo -help`

To download to a specific directory:
`getgo -dir ~/Downloads`

NOTE: getgo checks https://golang.org/dl/?mode=json to determine the latest
stable version of Go.

## Experimental GUI version

I'm experimenting and comparing different front-end frameworks starting with
[Fyne](https://fyne.io) and [GTK3](https://github.com/gotk3/gotk3).  To avoid
requiring users from having to install the different frameworks that they might
not to have on their system to compile a version, I'm using build tags to
control what gets compiled.

To try the Fyne version (require installing Fyne and its dependencies):

```
go get lazyhacker.dev/getgo
go run -tags gui,fyne lazyhacker.dev/getgo --gui

```

To try the GTK3 version (requires installing the GTK3 development libraries):

```
go get lazyhacker.dev/getgo
go run -tags gui,gtk3 lazyhacker.dev/getgo --gui

```

