package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"synth.com/file_converter/internal/response" 
	"synth.com/file_converter/internal/utils"
	"github.com/chai2010/webp" // WebP support
	"github.com/nfnt/resize"
	"github.com/signintech/gopdf"
	"image"
	"image/jpeg"
	"image/png"
)

// ConvertFile handles the logic to convert the file based on target format.
func ConvertFile(file io.Reader, targetFormat string) response.APIResponse {
	validFormats := []string{"pdf", "jpg", "png", "webp"}
	if !utils.Contains(validFormats, targetFormat) {
		return response.NewErrorResponse(400, "Invalid target format")
	}

	if targetFormat == "webp" || targetFormat == "png" || targetFormat == "jpg" {
		data, err := ConvertImage(file, targetFormat)
		if err != nil {
			return response.NewErrorResponse(500, err.Error())
		}
		return response.NewSuccessResponse("Image converted successfully", data)
	}

	if targetFormat == "pdf" {
		data, err := ConvertToPDF(file)
		if err != nil {
			return response.NewErrorResponse(500, err.Error())
		}
		return response.NewSuccessResponse("PDF generated successfully", data)
	}

	return response.NewErrorResponse(400, "Unsupported file format")
}

// ConvertImage converts an image file to the target format (PNG, JPEG, or WebP).
func ConvertImage(file io.Reader, targetFormat string) ([]byte, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("image decode failed: %w", err)
	}

	resizedImg := resize.Resize(800, 0, img, resize.Lanczos3)

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
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("image decode failed: %w", err)
	}

	tmpFile, err := saveImageToTempFile(img)
	if err != nil {
		return nil, fmt.Errorf("unable to save image to temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up temp file

	var buf bytes.Buffer
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 595, H: 842}}) // A4 page size
	pdf.AddPage()

	err = pdf.Image(tmpFile.Name(), 10, 10, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to insert image into PDF: %w", err)
	}

	err = pdf.Write(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to write PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// saveImageToTempFile saves the image to a temporary file.
func saveImageToTempFile(img image.Image) (*os.File, error) {
	tmpFile, err := os.CreateTemp("", "image-*.png")
	if err != nil {
		return nil, err
	}

	err = png.Encode(tmpFile, img)
	if err != nil {
		tmpFile.Close()
		return nil, err
	}

	tmpFile.Seek(0, 0)

	return tmpFile, nil
}
