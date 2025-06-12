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

func FileChecksum(file *os.File) (string, error){

	// move 0 bytes from the current cursor position
	// so in effect, the cursor stays at the current position
	// get current position
	curr, _ := file.Seek(0, io.SeekCurrent)

	// at the end of the func, move the cursor to where it was when we recieve this file
	// we're saying "set cursor position to curr (0) from the start of the file."
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

/*
	doing this: curr, _ := file.Seek(0, io.SeekCurrent)
	is like a hedge. "we don't know where we are in the file, but record the current position so we can come back to it later"
*/
