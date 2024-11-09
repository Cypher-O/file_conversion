package utils

import (
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Contains checks if a slice contains a specific element
func Contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

// SaveImageToTempFile saves an image to a temporary file and returns the file
func SaveImageToTempFile(img image.Image, filename string) (*os.File, error) {
	tmpDir := os.TempDir()
	tmpFile, err := ioutil.TempFile(tmpDir, filepath.Base(filename)+".*")
	if err != nil {
		return nil, err
	}

	err = png.Encode(tmpFile, img)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, err
	}
	tmpFile.Seek(0, 0)
	return tmpFile, nil
}
