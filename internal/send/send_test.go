package send

import (
	"encoding/binary"
	"encoding/json"
	"ftcli/internal/shared"
	"ftcli/models"
	"io"
	"net"
	"reflect"
	"testing"
)

// Test send of header
// This test was written by chatGPT
func TestSendHeader(t *testing.T) {
	c1, c2 := net.Pipe() // c1: writer, c2: reader
	defer c1.Close()
	defer c2.Close()

	header := models.Header{
		FileName: "foo.bar",
		CheckSum: "11586d2eb43b73e539caa3d158c883336c0e2c904b309c0c5ffe2c9b83d562a1",
	}

	hdrJSON, err := shared.HeaderToJsonB(header)
	if err != nil {
		t.Fatalf("failed to marshal header: %v", err)
	}
	hdrLen := shared.GetHeaderLength(hdrJSON)

	go func() {
		err := sendHeader(c1, hdrJSON, hdrLen)
		if err != nil {
			t.Errorf("sendHeader failed: %v", err)
		}
	}()

	// Read 4-byte length prefix
	var lenBuf [4]byte
	if _, err := io.ReadFull(c2, lenBuf[:]); err != nil {
		t.Fatalf("failed to read length prefix: %v", err)
	}
	length := binary.BigEndian.Uint32(lenBuf[:])

	// Read JSON payload
	jsonBuf := make([]byte, length)
	if _, err := io.ReadFull(c2, jsonBuf); err != nil {
		t.Fatalf("failed to read JSON data: %v", err)
	}

	var received models.Header
	if err := json.Unmarshal(jsonBuf, &received); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// Compare original and received headers
	if !reflect.DeepEqual(header, received) {
		t.Errorf("mismatch:\nexpected: %+v\ngot: %+v", header, received)
	}
}
