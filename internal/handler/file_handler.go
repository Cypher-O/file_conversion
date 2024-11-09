package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"synth.com/file_converter/internal/response"
	"synth.com/file_converter/internal/service"
)

// ConvertFileHandler handles the file upload and conversion
func ConvertFileHandler(c *gin.Context) {
	// Log to check if the request is reaching the handler
	log.Println("Received request for file conversion")

	// Parse the file from the form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		// Log the error for debugging purposes
		log.Println("Error parsing file:", err)
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(400, "Unable to parse the file"))
		return
	}
	defer file.Close()

	// Parse the target format (e.g., 'jpg', 'png', 'pdf', etc.)
	targetFormat := c.DefaultQuery("format", "")
	if targetFormat == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(400, "Format query parameter is required"))
		return
	}

	// Call the service layer to handle file conversion
	resp := service.ConvertFile(file, header.Filename, targetFormat)

	// Check if the conversion was successful
	if resp.Code != 0 {
		// Log the error from the service layer
		log.Println("Error during conversion:", resp.Message)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// Send the converted file as the response
	c.Header("Content-Disposition", "attachment; filename=converted_file."+targetFormat)
	c.Data(http.StatusOK, "application/octet-stream", resp.Data.([]byte))
}
