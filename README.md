getgo checks https://golang.org/dl/?mode=json to determine the latest
stable version of [Go](https://golang.org) to download and verifies its
checksum.  The tool won't download the binary if a verified one already exists
locally.

It is command line utility with an experimental GUI mode that is only
half-baked.

![animemated screenshot](getgo.gif)

## Install

### Pre-built Binaries

Pre-built binaries are available in the releases section or it can be built if
you already have Go installed:

`go get lazyhacker.dev/getgo`

### Compile From Source

If you've downloaded the source then the standard Go tool for building can be
used:

```
go build main.go
```

## Usage

Just run `getgo` to download the most recent stable archive for the platform it
is running from.  To also extract the archive run it with the '-x' flag:

```
getgo -x <dir to extract to>
```
On Windows, getgo can be told to download the installer with

```
getgo -kind installer
```

To download to a specific directory:

`getgo -dir ~/Downloads`

To download another OS and/or arch version, use the '-os' and the '-arch' flags:

```
getgo -os windows -arch amd64
```

To get help info:

`getgo -help`

### Full Example

I use getgo make it easier for me to upgrade Go when a new release comes out.
The steps are generally:

1. Download Go from golang.org.
1. Verify the download.
1. Delete /usr/local/go.
1. Untar the downloaded .tar.gz file to /usr/local

With getgo, this will do the above:

```
sudo rm /usr/local/go
sudo getgo -x /usr/local
```

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

