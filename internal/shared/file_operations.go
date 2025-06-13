package shared

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
)

// Copies the file from the source to a destination and calcualates the sha265 checksum of the file
func CopyAndHash(dst io.Writer, src io.Reader) (string, int64, error) {

	h := sha256.New()
	mw := io.MultiWriter(dst, h)
	bytesWritten, err := io.Copy(mw, src)
	if err != nil{
		return "", 0, err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), bytesWritten, nil

}

// Calculates the sha256 hash for the file
func FileChecksumSHA265(file *os.File) (string, error){

	// record the current position of the cursor
	// we do this because we don't know at what stage the file maybe passed.
	// we record the cursor position so we can reset back to the location
	curr, _ := file.Seek(0, io.SeekCurrent)

	// at the end of the func, move the cursor to where it was when we recieve this file
	defer file.Seek(curr, io.SeekStart)

	// go to start of file
	file.Seek(0, io.SeekStart)

	// copy contents into the hash
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil{
		return  "", err
	}

	checksum := fmt.Sprintf("%x", h.Sum(nil))
	// defer executes...

	return checksum, nil
}