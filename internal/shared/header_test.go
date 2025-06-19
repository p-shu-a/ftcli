package shared

import (
	"encoding/binary"
	"ftcli/models"
	"reflect"
	"testing"
)

func TestHeaderEncodingAndDecoding(t *testing.T) {
	hdr := models.Header{FileName: "file.txt", CheckSum: "abc", Nonce: []byte{1, 2, 3}, Salt: []byte{4, 5}, IV: []byte{6}}
	jsonBytes, err := HeaderToJsonB(hdr)
	if err != nil {
		t.Fatalf("HeaderToJsonB failed: %v", err)
	}
	decoded, err := JsonBToHeader(jsonBytes)
	if err != nil {
		t.Fatalf("JsonBToHeader failed: %v", err)
	}
	if !reflect.DeepEqual(hdr, *decoded) {
		t.Errorf("header mismatch after encode/decode\nexpected:%+v\n got:%+v", hdr, *decoded)
	}
}

func TestGetHeaderLength(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	lenBytes := GetHeaderLength(data)
	if len(lenBytes) != 4 {
		t.Fatalf("expected 4 bytes, got %d", len(lenBytes))
	}
	if l := binary.BigEndian.Uint32(lenBytes); l != uint32(len(data)) {
		t.Errorf("expected length %d, got %d", len(data), l)
	}
}
