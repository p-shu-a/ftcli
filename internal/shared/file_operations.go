package shared

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
)

// Copies the file from the source to a destination and calcualates the sha265 checksum of the file
func CopyAndHash(dst io.Writer, src io.Reader) (string, int64, error) {

	h := sha256.New()
	mw := io.MultiWriter(dst, h)
	bytesWritten, err := io.Copy(mw, src)
	if err != nil {
		return "", 0, err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), bytesWritten, nil

}

// Calculates the sha256 hash for the file
func FileChecksumSHA265(file *os.File) (string, error) {

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
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}

	checksum := fmt.Sprintf("%x", h.Sum(nil))
	// defer executes...

	return checksum, nil
}

// Given the original file name, and the counter, sends back a new filename
// Just to keep the code clean :P
func SuggestNewFileName(ogFilename string, ctr int) string {
	fileSplit := strings.Split(ogFilename, ".")
	// save the last index as the ext.
	ext := fileSplit[len(fileSplit)-1]
	// replace last index with counter
	fileSplit[len(fileSplit)-1] = fmt.Sprintf("%d", ctr)
	// join file and add extension
	newName := strings.Join(fileSplit, ".")
	newName = fmt.Sprintf(newName+".%v", ext)
	return newName
}
