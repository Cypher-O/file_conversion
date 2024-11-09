package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"synth.com/file_converter/internal/response"
	"synth.com/file_converter/internal/utils"
	"github.com/nfnt/resize"
	"github.com/signintech/gopdf"
	"image"
	"image/jpeg"
	"image/png"
	"github.com/chai2010/webp" // WebP support
	"os"
)

// ConvertFile handles the logic to convert the file based on target format.
func ConvertFile(file io.Reader, targetFormat string) response.APIResponse {
	validFormats := []string{"pdf", "jpg", "png", "webp"}
	if !utils.Contains(validFormats, targetFormat) {
		// Return error response if the target format is invalid
		return response.NewErrorResponse(400, "Invalid target format")
	}

	// Convert the file based on the target format
	if targetFormat == "webp" || targetFormat == "png" || targetFormat == "jpg" {
		data, err := ConvertImage(file, targetFormat)
		if err != nil {
			// Log the error for debugging purposes
			return response.NewErrorResponse(500, fmt.Sprintf("Image conversion failed: %s", err.Error()))
		}
		return response.NewSuccessResponse("Image converted successfully", data)
	}

	if targetFormat == "pdf" {
		data, err := ConvertToPDF(file)
		if err != nil {
			// Log the error for debugging purposes
			return response.NewErrorResponse(500, fmt.Sprintf("PDF conversion failed: %s", err.Error()))
		}
		return response.NewSuccessResponse("PDF generated successfully", data)
	}

	// If the target format is unsupported, return an error
	return response.NewErrorResponse(400, "Unsupported file format")
}

// ConvertImage converts an image file to the target format (PNG, JPEG, or WebP).
func ConvertImage(file io.Reader, targetFormat string) ([]byte, error) {
	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		// Log the error for debugging purposes
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize the image (optional resizing)
	resizedImg := resize.Resize(800, 0, img, resize.Lanczos3)

	// Prepare the buffer to store the converted image
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
		// Log the error for debugging purposes
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}
	return buf.Bytes(), nil
}

// ConvertToPDF converts an image to a PDF document.
func ConvertToPDF(file io.Reader) ([]byte, error) {
	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		// Log the error for debugging purposes
		return nil, fmt.Errorf("image decode failed: %w", err)
	}

	// Create a temporary file to save the image
	tmpFile, err := saveImageToTempFile(img)
	if err != nil {
		// Log the error for debugging purposes
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
		// Log the error for debugging purposes
		return nil, fmt.Errorf("failed to insert image into PDF: %w", err)
	}

	// Write the PDF to a buffer
	err = pdf.Write(&buf)
	if err != nil {
		// Log the error for debugging purposes
		return nil, fmt.Errorf("failed to write PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// saveImageToTempFile saves the image to a temporary file.
func saveImageToTempFile(img image.Image) (*os.File, error) {
	tmpFile, err := os.CreateTemp("", "image-*.png")
	if err != nil {
		// Log the error for debugging purposes
		return nil, err
	}

	// Save the image to the temp file
	err = png.Encode(tmpFile, img)
	if err != nil {
		tmpFile.Close()
		// Log the error for debugging purposes
		return nil, err
	}

	// Rewind the file to the beginning so gopdf can read it
	tmpFile.Seek(0, 0)

	return tmpFile, nil
}
