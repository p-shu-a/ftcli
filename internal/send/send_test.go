package send

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"ftcli/internal/encryption"
	"ftcli/internal/shared"
	"ftcli/models"
)

// Test send of header
func TestSendHeader(t *testing.T) {
	c1, c2 := net.Pipe()
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
		if err := sendHeader(c1, hdrJSON, hdrLen); err != nil {
			t.Errorf("sendHeader failed: %v", err)
		}
	}()

	var lenBuf [4]byte
	if _, err := io.ReadFull(c2, lenBuf[:]); err != nil {
		t.Fatalf("failed to read length prefix: %v", err)
	}
	length := binary.BigEndian.Uint32(lenBuf[:])

	jsonBuf := make([]byte, length)
	if _, err := io.ReadFull(c2, jsonBuf); err != nil {
		t.Fatalf("failed to read JSON data: %v", err)
	}

	var received models.Header
	if err := json.Unmarshal(jsonBuf, &received); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if received.FileName != header.FileName || received.CheckSum != header.CheckSum {
		t.Errorf("mismatch:\nexpected: %+v\n got: %+v", header, received)
	}
}

func TestDialRemote(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:7891")
	if err != nil {
		t.Fatalf("failed to listen on port: %v", err)
	}
	defer ln.Close()

	done := make(chan struct{})
	go func() {
		if conn, err := ln.Accept(); err == nil {
			conn.Close()
		}
		close(done)
	}()

	conn, err := dialRemote(net.ParseIP("127.0.0.1"))
	if err != nil {
		t.Fatalf("dialRemote returned error: %v", err)
	}
	conn.Close()
	<-done
}

func TestDialRemoteFail(t *testing.T) {
	if _, err := dialRemote(net.ParseIP("127.0.0.1")); err == nil {
		t.Fatalf("expected error when no server is listening")
	}
}

func testReceiveHeader(conn net.Conn) ([]byte, error) {
	var lenBuf [4]byte
	if _, err := io.ReadFull(conn, lenBuf[:]); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lenBuf[:])
	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, err
	}
	return data, nil
}

func TestSendFile(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:7891")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	file, err := os.CreateTemp("", "sendfile*.txt")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	defer os.Remove(file.Name())

	content := []byte("this is a test file")
	if _, err := file.Write(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	file.Seek(0, io.SeekStart)

	expectedHash, err := shared.FileChecksumSHA265(file)
	if err != nil {
		t.Fatalf("hash file: %v", err)
	}
	file.Seek(0, io.SeekStart)
	expectedName := filepath.Base(file.Name())

	recv := make(chan []byte)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			t.Error(err)
			close(recv)
			return
		}
		defer conn.Close()

		infoHdrBytes, err := testReceiveHeader(conn)
		if err != nil {
			t.Error(err)
			close(recv)
			return
		}
		infoHdr, err := shared.JsonBToHeader(infoHdrBytes)
		if err != nil {
			t.Error(err)
			close(recv)
			return
		}
		if infoHdr.FileName != expectedName || infoHdr.CheckSum != expectedHash {
			t.Errorf("info header mismatch")
			close(recv)
			return
		}

		chunkHdrBytes, err := testReceiveHeader(conn)
		if err != nil {
			t.Error(err)
			close(recv)
			return
		}
		chunkHdr, err := shared.JsonBToHeader(chunkHdrBytes)
		if err != nil {
			t.Error(err)
			close(recv)
			return
		}

		cipher, err := io.ReadAll(conn)
		if err != nil {
			t.Error(err)
			close(recv)
			return
		}

		key := encryption.GenerateMasterKey(infoHdr.Salt, "pass")
		plain, err := encryption.DecryptAEAD(chunkHdr.Nonce, key, cipher, chunkHdrBytes)
		if err != nil {
			t.Error(err)
			close(recv)
			return
		}
		recv <- append([]byte{}, plain...)
	}()

	wg := sync.WaitGroup{}
	wg.Add(1)
	if err := SendFile(context.Background(), &wg, file, net.ParseIP("127.0.0.1"), "pass"); err != nil {
		t.Fatalf("SendFile error: %v", err)
	}
	wg.Wait()

	got := <-recv
	if !bytes.Equal(got, content) {
		t.Errorf("received data mismatch: got %q want %q", got, content)
	}
}
