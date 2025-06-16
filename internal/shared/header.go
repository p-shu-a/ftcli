package shared

import (
	"encoding/binary"
	"encoding/json"
	"ftcli/models"
)

// Takes a Header struct, and returns json-encoded bytes of the header and error
func HeaderToJsonB(header models.Header) ([]byte, error) {
	// marshal to json formatted bytes
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return nil, err
	}
	return headerBytes, nil
}

// Returns a byte-encoded length (as a uint32) of the length of the provided json byte-encoded header
func GetHeaderLength(hdrJsonBytes []byte) []byte {
	// calcualte the length of the header
	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(hdrJsonBytes)))
	return lenBuf[:]
}

// Takes the header in the from of json-encoded bytes and returns a Header struct
func JsonBToHeader(hdrJsonBytes []byte) (*models.Header, error) {
	var hdr models.Header
	if err := json.Unmarshal(hdrJsonBytes, &hdr); err != nil {
		return &models.Header{}, err
	}
	return &hdr, nil
}
