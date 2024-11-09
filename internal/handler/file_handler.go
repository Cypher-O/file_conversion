package handler

import (
	"github.com/gin-gonic/gin"
	"synth.com/file_converter/internal/service"
	"net/http"
	"fmt"
	"log"
)

func ConvertFileHandler(c *gin.Context) {
	// Log to check if the request is reaching the handler
	log.Println("Received request for file conversion")
	
	// Parse the file from the form
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		// Log the error for debugging purposes
		log.Println("Error parsing file:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse the file"})
		return
	}

	// Log the filename for debugging purposes
	log.Println("File received:", file)

	// Parse the target format (could be 'jpg', 'png', 'pdf', etc.)
	targetFormat := c.DefaultQuery("format", "pdf")

	// Log the requested format
	log.Println("Requested target format:", targetFormat)

	// Call the service layer to handle file conversion
	resp := service.ConvertFile(file, targetFormat)

	// Check if the conversion was successful
	if resp.Code != 0 {
		// Log the error from the service layer
		log.Println("Error during conversion:", resp.Message)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// Send the converted file as the response
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=converted_file.%s", targetFormat))
	c.Data(http.StatusOK, "application/octet-stream", resp.Data.([]byte))
}
