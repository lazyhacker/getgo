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

![gui screenshot](getgo-gui.png)


I'm experimenting with building GUI apps with Go.  [Fyne](https://fyne.io) is a
cross-platform Go GUI toolkit that I started testing.

To try the Fyne version (require [installing Fyne](https://developer.fyne.io/#installing)
and its [dependencies](https://developer.fyne.io/started/#prerequisites)):

```
go get lazyhacker.dev/getgo
go run -tags gui lazyhacker.dev/getgo

```
