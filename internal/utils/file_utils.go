package utils

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/signintech/gopdf"
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

// GeneratePDFFromImage generates a PDF from the provided image
func GeneratePDFFromImage(imageFile *os.File) ([]byte, error) {
	// Open the image file
	img, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Initialize a new PDF document
	var buf bytes.Buffer
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: gopdf.Rect{W: 595.28, H: 841.89}, // A4 size
		Unit:     gopdf.Unit_PT,
	})

	// Add a page
	pdf.AddPage()

	// Convert the image to fit within the A4 page dimensions
	// Resize the image if necessary to fit within the page
	pdf.ImageFrom(img, 0, 0, nil)

	// Write PDF to buffer
	err = pdf.Write(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to write PDF: %w", err)
	}

	return buf.Bytes(), nil
}
