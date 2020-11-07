package utils


import (
	"bufio"
	"encoding/binary"
	"errors"
  "hash/crc32"
	"io"

	"github.com/golang/protobuf/proto"
	protobuf "github.com/emla2805/tfr/protobuf"
)

const (
	kMaskDelta = 0xa282ead8
)

var (
  crc32c = crc32.MakeTable(crc32.Castagnoli)
)


// Reader implements a reader for TFRecords with Example protos
type Reader struct {
	reader *bufio.Reader
}

// NewReader returns a new Reader
func NewReader(r io.Reader) *Reader {
	return &Reader{
		reader: bufio.NewReader(r),
	}
}

// Verify checksum
func (w *Reader) verifyChecksum(data []byte, crcMasked uint32) bool {
	rot := crcMasked - kMaskDelta
	unmaskedCrc := ((rot >> 17) | (rot << 15))

	crc := crc32.Checksum(data, crc32c)

	return crc == unmaskedCrc
}

// Next reads the next Example from the TFRecords input
func (r *Reader) Next() (*protobuf.Example, error) {
	header := make([]byte, 12)
	_, err := io.ReadFull(r.reader, header)
	if err != nil {
		return nil, err
	}

	crc := binary.LittleEndian.Uint32(header[8:12])
	if !r.verifyChecksum(header[0:8], crc) {
		return nil, errors.New("Invalid crc for length")
	}

	length := binary.LittleEndian.Uint64(header[0:8])

	payload := make([]byte, length)
	_, err = io.ReadFull(r.reader, payload)
	if err != nil {
		return nil, err
	}

	footer := make([]byte, 4)
	_, err = io.ReadFull(r.reader, footer)
	if err != nil {
		return nil, err
	}

	crc = binary.LittleEndian.Uint32(footer[0:4])
	if !r.verifyChecksum(payload, crc) {
		return nil, errors.New("Invalid crc for payload")
	}

	ex := &protobuf.Example{}
	err = proto.Unmarshal(payload, ex)
	if err != nil {
		return nil, err
	}

	return ex, nil
}

// ExampleFeatureBytes is a helper function for decoding proto Bytes feature
// from a TensorFlow Example. If key is not found it returns default value
func ExampleFeatureBytes(example *protobuf.Example, key string) []byte {
  print(example.GetFeatures())
	// TODO: return error if key is not found?
	f, ok := example.Features.Feature[key]
	if !ok {
		return nil
	}

	val, ok := f.Kind.(*protobuf.Feature_BytesList)
	if !ok {
		return nil
	}

	return val.BytesList.Value[0]
}
