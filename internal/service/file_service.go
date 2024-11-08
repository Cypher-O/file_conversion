// package service

// import (
// 	"bytes"
// 	"errors"
// 	"image"
//     "image/png"
// 	"synth.com/file_converter/internal/utils"
// 	"io"
// 	// "os"
// 	// "strings"
// 	"github.com/nfnt/resize"
// 	"github.com/signintech/gopdf"
// )

// func ConvertFile(file io.Reader, targetFormat string) ([]byte, error) {
// 	// Here we are just implementing a basic image conversion example. You could extend this for other file types.
// 	// Validate target format
// 	validFormats := []string{"pdf", "jpg", "png", "webp"}
// 	if !utils.Contains(validFormats, targetFormat) {
// 		return nil, errors.New("invalid target format")
// 	}

// 	// Example: Converting an image to WebP or PNG
// 	if targetFormat == "webp" || targetFormat == "png" {
// 		return ConvertImage(file, targetFormat)
// 	}

// 	// Example: Convert to PDF
// 	if targetFormat == "pdf" {
// 		return ConvertToPDF(file)
// 	}

// 	return nil, errors.New("unsupported file format")
// }

// func ConvertImage(file io.Reader, targetFormat string) ([]byte, error) {
// 	// Example: Convert a file to image format (e.g., PNG or WebP)
// 	img, _, err := image.Decode(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Resize image if needed (optional)
// 	resizedImg := resize.Resize(800, 0, img, resize.Lanczos3)

// 	// Convert to the desired format (PNG or WebP)
// 	var buf bytes.Buffer
// 	switch targetFormat {
// 	case "png":
// 		err = png.Encode(&buf, resizedImg)
// 	case "webp":
// 		err = webp.Encode(&buf, resizedImg, nil)
// 	default:
// 		return nil, errors.New("unsupported image format")
// 	}

// 	if err != nil {
// 		return nil, err
// 	}
// 	return buf.Bytes(), nil
// }

// func ConvertToPDF(file io.Reader) ([]byte, error) {
// 	// Example: Convert an image file to a PDF document
// 	img, _, err := image.Decode(file)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var buf bytes.Buffer
// 	pdf := gopdf.GoPdf{}
// 	pdf.Start(gopdf.Config{PageSize: gopdf.PageSizeA4})
// 	pdf.AddPage()

// 	// Insert image into PDF (simple example)
// 	err = pdf.ImageFromReader(img, 10, 10, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Write the PDF to a buffer
// 	err = pdf.Write(&buf)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return buf.Bytes(), nil
// }

package service

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"synth.com/file_converter/internal/utils"
	"github.com/nfnt/resize"
	"github.com/signintech/gopdf"
	_ "image/jpeg"    // JPEG support
	_ "image/png"     // PNG support
	"github.com/chai2010/webp" // WebP support
)

// ConvertFile handles the logic to convert the file based on target format.
func ConvertFile(file io.Reader, targetFormat string) ([]byte, error) {
	// Validate target format
	validFormats := []string{"pdf", "jpg", "png", "webp"}
	if !utils.Contains(validFormats, targetFormat) {
		return nil, errors.New("invalid target format")
	}

	// Example: Converting an image to WebP, PNG, or JPEG
	if targetFormat == "webp" || targetFormat == "png" || targetFormat == "jpg" {
		return ConvertImage(file, targetFormat)
	}

	// Example: Convert to PDF
	if targetFormat == "pdf" {
		return ConvertToPDF(file)
	}

	return nil, errors.New("unsupported file format")
}

// ConvertImage converts an image file to the target format (PNG, JPEG, or WebP).
func ConvertImage(file io.Reader, targetFormat string) ([]byte, error) {
	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("image decode failed: %w", err)
	}

	// Resize the image (optional, here resizing it to width 800px)
	resizedImg := resize.Resize(800, 0, img, resize.Lanczos3)

	// Convert the image to the desired format
	var buf bytes.Buffer
	switch targetFormat {
	case "png":
		err = png.Encode(&buf, resizedImg)
	case "webp":
		err = webp.Encode(&buf, resizedImg, nil)
	case "jpg":
		err = jpeg.Encode(&buf, resizedImg, nil)
	default:
		return nil, errors.New("unsupported image format")
	}

	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ConvertToPDF converts an image to a PDF document.
func ConvertToPDF(file io.Reader) ([]byte, error) {
	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("image decode failed: %w", err)
	}

	// Create a temporary file to save the image
	tmpFile, err := saveImageToTempFile(img)
	if err != nil {
		return nil, fmt.Errorf("unable to save image to temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up temp file

	// Generate the PDF
	var buf bytes.Buffer
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 595, H: 842}}) // A4 page size
	pdf.AddPage()

	// Insert the image into the PDF
	err = pdf.Image(tmpFile.Name(), 10, 10, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to insert image into PDF: %w", err)
	}

	// Write the PDF to a buffer
	err = pdf.Write(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to write PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// saveImageToTempFile saves the image to a temporary file.
func saveImageToTempFile(img image.Image) (*os.File, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "image-*.png")
	if err != nil {
		return nil, err
	}

	// Save the image to the temp file
	err = png.Encode(tmpFile, img)
	if err != nil {
		tmpFile.Close()
		return nil, err
	}

	// Rewind the file to the beginning so gopdf can read it
	tmpFile.Seek(0, 0)

	return tmpFile, nil
}
