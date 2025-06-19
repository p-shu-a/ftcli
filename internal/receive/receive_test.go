package receive

import (
	"net"
	"reflect"
	"testing"

	"ftcli/internal/shared"
	"ftcli/models"
)

func TestReceiveHeader(t *testing.T) {
	c1, c2 := net.Pipe()
	defer c1.Close()
	defer c2.Close()

	header := models.Header{FileName: "foo.txt", CheckSum: "123"}
	hdrBytes, err := shared.HeaderToJsonB(header)
	if err != nil {
		t.Fatalf("failed to marshal header: %v", err)
	}
	lenBuf := shared.GetHeaderLength(hdrBytes)

	go func() {
		c1.Write(lenBuf)
		c1.Write(hdrBytes)
		c1.Close()
	}()

	data, err := receiveHeader(c2)
	if err != nil {
		t.Fatalf("receiveHeader returned error: %v", err)
	}
	decoded, err := shared.JsonBToHeader(data)
	if err != nil {
		t.Fatalf("failed to decode header: %v", err)
	}
	if !reflect.DeepEqual(header, *decoded) {
		t.Errorf("decoded header mismatch: expected %+v got %+v", header, *decoded)
	}
}
