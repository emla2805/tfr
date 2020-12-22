module github.com/emla2805/tfr

go 1.15

require (
	github.com/spf13/cobra v1.1.1
	google.golang.org/protobuf v1.25.0
)

replace github.com/tensorflow/tensorflow/tensorflow/go/core => ./proto/tensorflow/core
