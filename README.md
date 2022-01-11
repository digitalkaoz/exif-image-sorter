# Exif Image Sorter

> read a folder containing images and move all files into a new folder structure created by the exif informations.

- it falls back to some filename parsing if no exif date could be found
- it also modifies the file utimes to match the exif date

There are 2 implementations that are doing the same (its more like a toy project with a concrete need for myself)

## NodeJS

### Installation

```shell
$ cd node
$ yarn
$ yarn build
$ yarn package
```

> this will create a binary called `sort-images`

## Golang

```shell
$ cd go
$ go get
$ go build
```

## Usage

> both versions spit out the same file with the same interface

```shell
$ ./exif-image-sorter SRC_FOLDER DEST_FOLDER
```

## Benchmarking

here are some stats with 18 files (containing exif data and parsing from filename)

```shell
# node
real	0m0.755s
user	0m0.176s
sys     0m0.405s
```

```shell
#golang
real	0m0.021s
user	0m0.017s
sys     0m0.026s
```

> real parallelism (multi proc) always wins against a simple concurrency approach (event-loop)
