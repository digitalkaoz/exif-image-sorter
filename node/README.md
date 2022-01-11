# Exif Image Sorter

sort your images by moving them into folders based on exif dates (or fallback to filename patterns)

## Installation

```shell
$ yarn
```

## Usage

```shell
$ yarn build
$ node ./dist/index.js

$ node ./dist/index.js SRC_FOLDER DEST_FOLDER
```

> it will created folders for each year and month found in your exif data.
> modifies the created/modified file timestamps based on your exif data as well (so stupid file browsers can sort them properly)

## Standalone Binary

```shell
$ yarn build

# "pkg -t node12-macos-x64 -o sort-files dist/index.js"
```

> creates a OSX build

for other OSes see [pkg](https://www.npmjs.com/package/pkg)
