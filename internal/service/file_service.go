// package service

// import (
// 	"bytes"
// 	"fmt"
// 	"image"
// 	"image/jpeg"
// 	"image/png"
// 	"io"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"strings"

// 	"github.com/chai2010/webp"
// 	"github.com/nfnt/resize"
// 	"github.com/signintech/gopdf"
// 	"github.com/xuri/excelize/v2"
// 	"synth.com/file_converter/internal/response"
// 	"synth.com/file_converter/internal/utils"
// )

// // ConvertFile handles the logic to convert the file based on target format
// func ConvertFile(file io.Reader, filename, targetFormat string) response.APIResponse {
// 	validFormats := []string{"pdf", "jpg", "png", "webp", "csv"}
// 	if !utils.Contains(validFormats, targetFormat) {
// 		return response.NewErrorResponse(400, "Invalid target format")
// 	}

// 	// Determine the file extension
// 	ext := filepath.Ext(filename)

// 	switch ext {
// 	case ".png", ".jpg", ".jpeg", ".webp":
// 		data, err := ConvertImage(file, targetFormat)
// 		if err != nil {
// 			return response.NewErrorResponse(500, fmt.Sprintf("Image conversion failed: %s", err.Error()))
// 		}
// 		return response.NewSuccessResponse("Image converted successfully", data)

// 	case ".docx":
// 		data, err := ConvertWordToPDF(file)
// 		if err != nil {
// 			return response.NewErrorResponse(500, fmt.Sprintf("Word to PDF conversion failed: %s", err.Error()))
// 		}
// 		return response.NewSuccessResponse("Word document converted to PDF successfully", data)

// 	case ".xlsx":
// 		if targetFormat == "csv" {
// 			data, err := ConvertExcelToCSV(file)
// 			if err != nil {
// 				return response.NewErrorResponse(500, fmt.Sprintf("Excel to CSV conversion failed: %s", err.Error()))
// 			}
// 			return response.NewSuccessResponse("Excel document converted to CSV successfully", data)
// 		}

// 	default:
// 		if targetFormat == "pdf" {
// 			data, err := ConvertToPDF(file, filename)
// 			if err != nil {
// 				return response.NewErrorResponse(500, fmt.Sprintf("PDF conversion failed: %s", err.Error()))
// 			}
// 			return response.NewSuccessResponse("PDF generated successfully", data)
// 		}
// 	}

// 	return response.NewErrorResponse(400, "Unsupported file format")
// }

// // ConvertImage converts an image file to the target format (PNG, JPEG, or WebP)
// func ConvertImage(file io.Reader, targetFormat string) ([]byte, error) {
// 	img, _, err := image.Decode(file)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to decode image: %w", err)
// 	}

// 	resizedImg := resize.Resize(800, 0, img, resize.Lanczos3)

// 	var buf bytes.Buffer
// 	switch targetFormat {
// 	case "png":
// 		err = png.Encode(&buf, resizedImg)
// 	case "webp":
// 		err = webp.Encode(&buf, resizedImg, nil)
// 	case "jpg":
// 		err = jpeg.Encode(&buf, resizedImg, nil)
// 	default:
// 		return nil, fmt.Errorf("unsupported image format")
// 	}

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to encode image: %w", err)
// 	}
// 	return buf.Bytes(), nil
// }

// // ConvertWordToPDF converts a Word document to a PDF document
// func ConvertWordToPDF(file io.Reader) ([]byte, error) {
//     // Create a temporary file to store the input
//     tmpInput, err := os.CreateTemp("", "input-*.docx")
//     if err != nil {
//         return nil, fmt.Errorf("failed to create temp file: %w", err)
//     }
//     defer os.Remove(tmpInput.Name())
    
//     // Copy the input to the temporary file
//     _, err = io.Copy(tmpInput, file)
//     if err != nil {
//         return nil, fmt.Errorf("failed to copy input to temp file: %w", err)
//     }
//     tmpInput.Close()

//     // Initialize PDF document
//     var buf bytes.Buffer
//     pdf := gopdf.GoPdf{}
//     pdf.Start(gopdf.Config{
//         PageSize: gopdf.Rect{W: 595.28, H: 841.89}, // A4 size
//         Unit:     gopdf.Unit_PT,
//     })

//     // Use pandoc for conversion (requires pandoc to be installed)
//     cmd := exec.Command("pandoc", 
//         tmpInput.Name(),
//         "-f", "docx",
//         "-t", "plain",
//         "--wrap=none")
    
//     output, err := cmd.Output()
//     if err != nil {
//         return nil, fmt.Errorf("failed to convert document: %w", err)
//     }

//     // Add content to PDF
//     pdf.AddPage()
//     err = pdf.AddTTFFont("default", "assets/fonts/Arial.ttf")
//     if err != nil {
//         return nil, fmt.Errorf("failed to add font: %w", err)
//     }

//     err = pdf.SetFont("default", "", 12)
//     if err != nil {
//         return nil, fmt.Errorf("failed to set font: %w", err)
//     }

//     // Split content into paragraphs and add to PDF
//     paragraphs := strings.Split(string(output), "\n\n")
//     for _, para := range paragraphs {
//         if strings.TrimSpace(para) != "" {
//             pdf.Cell(nil, para)
//             pdf.Br(20)
//         }
//     }

//     err = pdf.Write(&buf)
//     if err != nil {
//         return nil, fmt.Errorf("failed to write PDF: %w", err)
//     }

//     return buf.Bytes(), nil
// }

// // ConvertExcelToCSV converts an Excel document to a CSV format
// func ConvertExcelToCSV(file io.Reader) ([]byte, error) {
// 	xl, err := excelize.OpenReader(file)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read Excel document: %w", err)
// 	}

// 	var buf bytes.Buffer
// 	for _, sheetName := range xl.GetSheetList() {
// 		rows, err := xl.GetRows(sheetName)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get rows: %w", err)
// 		}
// 		for _, row := range rows {
// 			buf.WriteString(strings.Join(row, ","))
// 			buf.WriteString("\n")
// 		}
// 	}

// 	return buf.Bytes(), nil
// }

// // ConvertToPDF converts an image to a PDF document
// func ConvertToPDF(file io.Reader, filename string) ([]byte, error) {
// 	img, _, err := image.Decode(file)
// 	if err != nil {
// 		return nil, fmt.Errorf("image decode failed: %w", err)
// 	}

// 	tmpFile, err := utils.SaveImageToTempFile(img, filename)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to save image to temp file: %w", err)
// 	}
// 	defer os.Remove(tmpFile.Name())

// 	var buf bytes.Buffer
// 	pdf := gopdf.GoPdf{}
// 	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 595, H: 842}})
// 	pdf.AddPage()
// 	err = pdf.Image(tmpFile.Name(), 0, 0, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to add image to PDF: %w", err)
// 	}

// 	err = pdf.Write(&buf)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to write PDF: %w", err)
// 	}

// 	return buf.Bytes(), nil
// }



package service

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
	"github.com/signintech/gopdf"
	"github.com/xuri/excelize/v2"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"synth.com/file_converter/internal/response"
	"synth.com/file_converter/internal/utils"
)

// ConvertFile handles the logic to convert the file based on target format
func ConvertFile(file io.Reader, filename, targetFormat string) response.APIResponse {
	// List of valid formats
	validFormats := []string{"pdf", "txt", "jpg", "png", "webp", "csv"}
	if !utils.Contains(validFormats, targetFormat) {
		return response.NewErrorResponse(400, "Invalid target format")
	}

	// Determine the file extension
	ext := filepath.Ext(filename)

	// Process the file based on its extension
	switch ext {
	case ".png", ".jpg", ".jpeg", ".webp":
		return handleImageConversion(file, targetFormat)

	case ".docx":
		return handleWordToPDFConversion(file)

	case ".xlsx":
		if targetFormat == "csv" {
			return handleExcelToCSVConversion(file)
		}

	case ".pdf":
		if targetFormat == "txt" {
			return handlePDFToTextConversion(file) // Handle PDF to Text conversion
		}

	default:
		return handleDefaultPDFConversion(file, filename, targetFormat)
	}

	// Return error for unsupported formats
	return response.NewErrorResponse(400, "Unsupported file format")
}

// handleImageConversion processes image files (png, jpg, jpeg, webp)
func handleImageConversion(file io.Reader, targetFormat string) response.APIResponse {
	data, err := ConvertImage(file, targetFormat)
	if err != nil {
		return response.NewErrorResponse(500, fmt.Sprintf("Image conversion failed: %s", err.Error()))
	}
	return response.NewSuccessResponse("Image converted successfully", data)
}

// handleWordToPDFConversion processes Word document conversion to PDF
func handleWordToPDFConversion(file io.Reader) response.APIResponse {
	data, err := ConvertWordToPDF(file)
	if err != nil {
		return response.NewErrorResponse(500, fmt.Sprintf("Word to PDF conversion failed: %s", err.Error()))
	}
	return response.NewSuccessResponse("Word document converted to PDF successfully", data)
}

// handleExcelToCSVConversion processes Excel to CSV conversion
func handleExcelToCSVConversion(file io.Reader) response.APIResponse {
	data, err := ConvertExcelToCSV(file)
	if err != nil {
		return response.NewErrorResponse(500, fmt.Sprintf("Excel to CSV conversion failed: %s", err.Error()))
	}
	return response.NewSuccessResponse("Excel document converted to CSV successfully", data)
}

// handlePDFToTextConversion processes PDF to text conversion
func handlePDFToTextConversion(file io.Reader) response.APIResponse {
	data, err := ConvertPDFToText(file)
	if err != nil {
		return response.NewErrorResponse(500, fmt.Sprintf("PDF to text conversion failed: %s", err.Error()))
	}
	return response.NewSuccessResponse("PDF text extracted successfully", data)
}

// handleDefaultPDFConversion processes file conversion to PDF for unsupported formats
func handleDefaultPDFConversion(file io.Reader, filename, targetFormat string) response.APIResponse {
	if targetFormat == "pdf" {
		data, err := ConvertToPDF(file, filename)
		if err != nil {
			return response.NewErrorResponse(500, fmt.Sprintf("PDF conversion failed: %s", err.Error()))
		}
		return response.NewSuccessResponse("PDF generated successfully", data)
	}
	return response.NewErrorResponse(400, "Unsupported file format for PDF conversion")
}

// ConvertPDFToText extracts text from a PDF file
func ConvertPDFToText(file io.Reader) (string, error) {
	// Create a temporary file to store the input PDF
	tmpInput, err := os.CreateTemp("", "input-*.pdf")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpInput.Name())

	// Copy the input PDF to the temporary file
	_, err = io.Copy(tmpInput, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy input to temp file: %w", err)
	}
	tmpInput.Close()

	// Extract text from the PDF using pdfcpu
	var buf bytes.Buffer
	err = api.ExtractText(tmpInput.Name(), &buf)
	if err != nil {
		return "", fmt.Errorf("failed to extract text from PDF: %w", err)
	}

	// Return the extracted text as a string
	return buf.String(), nil
}

// ConvertImage converts an image file to the target format (PNG, JPEG, or WebP)
func ConvertImage(file io.Reader, targetFormat string) ([]byte, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize image while maintaining aspect ratio
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
		return nil, fmt.Errorf("unsupported image format")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}
	return buf.Bytes(), nil
}

// ConvertWordToPDF converts a Word document to a PDF document
func ConvertWordToPDF(file io.Reader) ([]byte, error) {
	// Create a temporary file to store the input
	tmpInput, err := os.CreateTemp("", "input-*.docx")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpInput.Name())

	// Copy the input to the temporary file
	_, err = io.Copy(tmpInput, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy input to temp file: %w", err)
	}
	tmpInput.Close()

	// Initialize PDF document
	var buf bytes.Buffer
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		PageSize: gopdf.Rect{W: 595.28, H: 841.89}, // A4 size
		Unit:     gopdf.Unit_PT,
	})

	// Use pandoc for conversion (requires pandoc to be installed)
	cmd := exec.Command("pandoc", tmpInput.Name(), "-f", "docx", "-t", "plain", "--wrap=none")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to convert document: %w", err)
	}

	// Add content to PDF
	pdf.AddPage()
	err = pdf.AddTTFFont("default", "assets/fonts/Arial.ttf")
	if err != nil {
		return nil, fmt.Errorf("failed to add font: %w", err)
	}

	err = pdf.SetFont("default", "", 12)
	if err != nil {
		return nil, fmt.Errorf("failed to set font: %w", err)
	}

	// Split content into paragraphs and add to PDF
	paragraphs := strings.Split(string(output), "\n\n")
	for _, para := range paragraphs {
		if strings.TrimSpace(para) != "" {
			pdf.Cell(nil, para)
			pdf.Br(20)
		}
	}

	err = pdf.Write(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to write PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// ConvertExcelToCSV converts an Excel document to CSV format
func ConvertExcelToCSV(file io.Reader) ([]byte, error) {
	xl, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read Excel document: %w", err)
	}

	var buf bytes.Buffer
	for _, sheetName := range xl.GetSheetList() {
		rows, err := xl.GetRows(sheetName)
		if err != nil {
			return nil, fmt.Errorf("failed to get rows: %w", err)
		}
		for _, row := range rows {
			buf.WriteString(strings.Join(row, ","))
			buf.WriteString("\n")
		}
	}

	return buf.Bytes(), nil
}

// ConvertToPDF converts an image to a PDF document
func ConvertToPDF(file io.Reader, filename string) ([]byte, error) {
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("image decode failed: %w", err)
	}

	tmpFile, err := utils.SaveImageToTempFile(img, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to save image to temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Convert the image to PDF (this could be extended with other logic)
	return utils.GeneratePDFFromImage(tmpFile)
}



// package service

// import (
// 	"bytes"
// 	"fmt"
// 	"image"
// 	"image/jpeg"
// 	"image/png"
// 	"io"
// 	"os"
// 	"os/exec"
// 	"path/filepath"
// 	"strings"

// 	"github.com/chai2010/webp"
// 	"github.com/nfnt/resize"
// 	"github.com/signintech/gopdf"
// 	"github.com/xuri/excelize/v2"
// 	"synth.com/file_converter/internal/response"
// 	"synth.com/file_converter/internal/utils"
// )

// // ConvertFile handles the logic to convert the file based on target format
// func ConvertFile(file io.Reader, filename, targetFormat string) response.APIResponse {
// 	// List of valid formats
// 	validFormats := []string{"pdf", "jpg", "png", "webp", "csv"}
// 	if !utils.Contains(validFormats, targetFormat) {
// 		return response.NewErrorResponse(400, "Invalid target format")
// 	}

// 	// Determine the file extension
// 	ext := filepath.Ext(filename)

// 	// Process the file based on its extension
// 	switch ext {
// 	case ".png", ".jpg", ".jpeg", ".webp":
// 		return handleImageConversion(file, targetFormat)

// 	case ".docx":
// 		return handleWordToPDFConversion(file)

// 	case ".xlsx":
// 		if targetFormat == "csv" {
// 			return handleExcelToCSVConversion(file)
// 		}

// 	default:
// 		return handleDefaultPDFConversion(file, filename, targetFormat)
// 	}

// 	// Return error for unsupported formats
// 	return response.NewErrorResponse(400, "Unsupported file format")
// }

// // handleImageConversion processes image files (png, jpg, jpeg, webp)
// func handleImageConversion(file io.Reader, targetFormat string) response.APIResponse {
// 	data, err := ConvertImage(file, targetFormat)
// 	if err != nil {
// 		return response.NewErrorResponse(500, fmt.Sprintf("Image conversion failed: %s", err.Error()))
// 	}
// 	return response.NewSuccessResponse("Image converted successfully", data)
// }

// // handleWordToPDFConversion processes Word document conversion to PDF
// func handleWordToPDFConversion(file io.Reader) response.APIResponse {
// 	data, err := ConvertWordToPDF(file)
// 	if err != nil {
// 		return response.NewErrorResponse(500, fmt.Sprintf("Word to PDF conversion failed: %s", err.Error()))
// 	}
// 	return response.NewSuccessResponse("Word document converted to PDF successfully", data)
// }

// // handleExcelToCSVConversion processes Excel to CSV conversion
// func handleExcelToCSVConversion(file io.Reader) response.APIResponse {
// 	data, err := ConvertExcelToCSV(file)
// 	if err != nil {
// 		return response.NewErrorResponse(500, fmt.Sprintf("Excel to CSV conversion failed: %s", err.Error()))
// 	}
// 	return response.NewSuccessResponse("Excel document converted to CSV successfully", data)
// }

// // handleDefaultPDFConversion processes file conversion to PDF for unsupported formats
// func handleDefaultPDFConversion(file io.Reader, filename, targetFormat string) response.APIResponse {
// 	if targetFormat == "pdf" {
// 		data, err := ConvertToPDF(file, filename)
// 		if err != nil {
// 			return response.NewErrorResponse(500, fmt.Sprintf("PDF conversion failed: %s", err.Error()))
// 		}
// 		return response.NewSuccessResponse("PDF generated successfully", data)
// 	}
// 	return response.NewErrorResponse(400, "Unsupported file format for PDF conversion")
// }

// // ConvertImage converts an image file to the target format (PNG, JPEG, or WebP)
// func ConvertImage(file io.Reader, targetFormat string) ([]byte, error) {
// 	img, _, err := image.Decode(file)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to decode image: %w", err)
// 	}

// 	// Resize image while maintaining aspect ratio
// 	resizedImg := resize.Resize(800, 0, img, resize.Lanczos3)

// 	var buf bytes.Buffer
// 	switch targetFormat {
// 	case "png":
// 		err = png.Encode(&buf, resizedImg)
// 	case "webp":
// 		err = webp.Encode(&buf, resizedImg, nil)
// 	case "jpg":
// 		err = jpeg.Encode(&buf, resizedImg, nil)
// 	default:
// 		return nil, fmt.Errorf("unsupported image format")
// 	}

// 	if err != nil {
// 		return nil, fmt.Errorf("failed to encode image: %w", err)
// 	}
// 	return buf.Bytes(), nil
// }

// // ConvertWordToPDF converts a Word document to a PDF document
// func ConvertWordToPDF(file io.Reader) ([]byte, error) {
// 	// Create a temporary file to store the input
// 	tmpInput, err := os.CreateTemp("", "input-*.docx")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to create temp file: %w", err)
// 	}
// 	defer os.Remove(tmpInput.Name())

// 	// Copy the input to the temporary file
// 	_, err = io.Copy(tmpInput, file)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to copy input to temp file: %w", err)
// 	}
// 	tmpInput.Close()

// 	// Initialize PDF document
// 	var buf bytes.Buffer
// 	pdf := gopdf.GoPdf{}
// 	pdf.Start(gopdf.Config{
// 		PageSize: gopdf.Rect{W: 595.28, H: 841.89}, // A4 size
// 		Unit:     gopdf.Unit_PT,
// 	})

// 	// Use pandoc for conversion (requires pandoc to be installed)
// 	cmd := exec.Command("pandoc", tmpInput.Name(), "-f", "docx", "-t", "plain", "--wrap=none")
// 	output, err := cmd.Output()
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to convert document: %w", err)
// 	}

// 	// Add content to PDF
// 	pdf.AddPage()
// 	err = pdf.AddTTFFont("default", "assets/fonts/Arial.ttf")
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to add font: %w", err)
// 	}

// 	err = pdf.SetFont("default", "", 12)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to set font: %w", err)
// 	}

// 	// Split content into paragraphs and add to PDF
// 	paragraphs := strings.Split(string(output), "\n\n")
// 	for _, para := range paragraphs {
// 		if strings.TrimSpace(para) != "" {
// 			pdf.Cell(nil, para)
// 			pdf.Br(20)
// 		}
// 	}

// 	err = pdf.Write(&buf)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to write PDF: %w", err)
// 	}

// 	return buf.Bytes(), nil
// }

// // ConvertExcelToCSV converts an Excel document to CSV format
// func ConvertExcelToCSV(file io.Reader) ([]byte, error) {
// 	xl, err := excelize.OpenReader(file)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read Excel document: %w", err)
// 	}

// 	var buf bytes.Buffer
// 	for _, sheetName := range xl.GetSheetList() {
// 		rows, err := xl.GetRows(sheetName)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get rows: %w", err)
// 		}
// 		for _, row := range rows {
// 			buf.WriteString(strings.Join(row, ","))
// 			buf.WriteString("\n")
// 		}
// 	}

// 	return buf.Bytes(), nil
// }

// // ConvertToPDF converts an image to a PDF document
// func ConvertToPDF(file io.Reader, filename string) ([]byte, error) {
// 	img, _, err := image.Decode(file)
// 	if err != nil {
// 		return nil, fmt.Errorf("image decode failed: %w", err)
// 	}

// 	tmpFile, err := utils.SaveImageToTempFile(img, filename)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to save image to temp file: %w", err)
// 	}
// 	defer os.Remove(tmpFile.Name())

// 	var buf bytes.Buffer
// 	pdf := gopdf.GoPdf{}
// 	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 595, H: 842}})
// 	pdf.AddPage()
// 	err = pdf.Image(tmpFile.Name(), 0, 0, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to add image to PDF: %w", err)
// 	}

// 	err = pdf.Write(&buf)
// 	if err != nil {
// 		return nil, fmt.Errorf("unable to write PDF: %w", err)
// 	}

// 	return buf.Bytes(), nil
// }
