package shared

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// Test CopyAndHash's copy from reader to writer and hash generation
func TestCopyAndHash(t *testing.T){

	// create reader and writer
	fileText := "This is a test\n"
	reader := strings.NewReader(fileText)
	var writer bytes.Buffer

	// calculate hash for the contents of the reader
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		t.Fatal(err)
	}
	fileChkSum := fmt.Sprintf("%x", hash.Sum(nil))
	
	// reset to the top of the reader
	reader.Seek(0, io.SeekStart)

	// copy from reader to writer and get hash of transfered contents
	retHashStr, _ , err := CopyAndHash(&writer, reader)
	if err != nil {
		t.Fatal(err)
	}

	// compare the values
	if retHashStr != fileChkSum {
		t.Errorf("Hashes don't match.\nExpected Hash: %v\nReturned Hash: %v", fileChkSum, retHashStr)
	}else{
		t.Logf("Hashes match.\nExpected Hash: %v\nReturned Hash: %v", fileChkSum, retHashStr)
	}
}

// Test FileChecksum ability to return real hash
func TestFileChecksum(t *testing.T){

	// create temp file in tempdir
	file, err := os.CreateTemp(os.TempDir(), "testfile.*")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	/*
		in the terminal you can do:
		echo "This is a test." > /tmp/foo
		sha256 /tmp/foo
		SHA256 (/tmp/foo) = 11586d2eb43b73e539caa3d158c883336c0e2c904b309c0c5ffe2c9b83d562a1
	*/

	// write some text to the file
	fileText := "This is a test\n"
	byteText := []byte(fileText) 
	file.Write(byteText)
	if err != nil{
		t.Fatal(err)
	}
	
	// move cursor to the top of the file
	file.Seek(0,io.SeekStart)

	// get the hash for the file
	hash := sha256.New()
	if _ , err = io.Copy(hash, file); err != nil {
		t.Fatal(err)
	}
	sum := fmt.Sprintf("%x",hash.Sum(nil))

	// use the function to get the checksum of the file
	chkSum, err := FileChecksumSHA265(file)
	if err != nil {
		t.Fatal(err)
	}

	if chkSum != sum {
		t.Errorf("Hashes don't match.\nExpected Hash: %v\nReturned Hash: %v", sum, chkSum)
	}else{
		t.Logf("Hashes match.\nExpected Hash: %v\nReturned Hash: %v", sum, chkSum)
	}
}