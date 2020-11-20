# tfr

`tfr` is a lightweight command-line `TFRecords` processor that 
reads serialized `.tfrecord` files and outputs to stdout in JSON format.

## Install

Binaries are available from the [releases](https://github.com/emla2805/tfr/releases/latest) page.

If you have Go installed, just run `go get`.

    go get github.com/emla2805/tfr

## Usage

Parse a single file on the terminal

```bash
tfr data_tfrecord-00000-of-00001
```

or, read from `stdin`

```bash
cat data_tfrecord-00000-of-00001 | tfr -n 1
```