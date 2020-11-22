package utils

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
)

const (
	maskDelta = 0xa282ead8
	headerLen = 12
	footerLen = 4
)

var (
	crc32c = crc32.MakeTable(crc32.Castagnoli)
)

// Verify checksum
func verifyChecksum(data []byte, crcMasked uint32) bool {
	rot := crcMasked - maskDelta
	unmaskedCrc := ((rot >> 17) | (rot << 15))

	crc := crc32.Checksum(data, crc32c)

	return crc == unmaskedCrc
}

// ScanTFRecord scans a single record
func ScanTFRecord(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if len(data) < headerLen {
		return 0, nil, nil
	}
	header := data[:headerLen]
	recordLen := int(binary.LittleEndian.Uint64(header[0:8]))
	crc := binary.LittleEndian.Uint32(header[8:12])
	if !verifyChecksum(header[0:8], crc) {
		return 0, nil, errors.New("Invalid crc for length")
	}
	if len(data) < headerLen+recordLen+footerLen {
		return 0, nil, nil
	}
	payload := data[headerLen : headerLen+recordLen]
	footer := data[headerLen+recordLen : headerLen+recordLen+footerLen]

	crc = binary.LittleEndian.Uint32(footer[0:4])
	if !verifyChecksum(payload, crc) {
		return 0, nil, errors.New("Invalid crc for payload")
	}
	return headerLen + recordLen + footerLen, payload, nil
}
