package shared

import (
	"crypto/sha256"
	"fmt"
	"io"
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