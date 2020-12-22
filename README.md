# tfr

`tfr` is a lightweight command-line `TFRecords` processor that 
reads serialized `.tfrecord` files and outputs to stdout in JSON format.

## Install

Binaries are available from the [releases](https://github.com/emla2805/tfr/releases/latest) page.

If you have Go installed, just run `go get`.

    go get github.com/emla2805/tfr

On MacOs, use [Homebrew](https://brew.sh).

    brew tap emla2805/tfr
    brew install tfr

## Usage

Parse a single file on the terminal

```bash
tfr data_tfrecord-00000-of-00001
```

or, read from `stdin`

```bash
cat data_tfrecord-00000-of-00001 | tfr -n 1
```

## Examples

`tfr` is best used with other great tools like [jq](https://github.com/stedolan/jq),
[gsutil](https://cloud.google.com/storage/docs/gsutil) and `gunzip`.

### Compressed tfrecords from Google Cloud Storage
```bash
gsutil cat gs://<bucket>/<path>/data_tfrecord-00000-of-00001.gz | gunzip | tfr -n 1 | jq .
```

### Flatten example structure
```bash
tfr data_tfrecord-00000-of-00001 | jq '.features.feature | to_entries | map( {(.key): .value[].value} ) | add'
{
  "age": [
    29
  ],
  "movie": [
    "The Shawshank Redemption",
    "Fight Club"
  ],
  "movie_ratings": [
    9,
    9.7
  ]
}
```
