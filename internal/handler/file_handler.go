package handler

import (
	"github.com/gin-gonic/gin"
	"synth.com/file_converter/internal/service"
	"net/http"
	"fmt"
)

func ConvertFileHandler(c *gin.Context) {
	// Parse form file
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse the file"})
		return
	}

	// Parse target format (could be 'jpg', 'png', 'pdf', etc.)
	targetFormat := c.DefaultQuery("format", "pdf")

	// Call the service layer to handle file conversion
	resp := service.ConvertFile(file, targetFormat)

	// Check if there is an error in the response (code != 0 means failure)
	if resp.Code != 0 {
		// Return an error response with the message from APIResponse
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// Send the converted file as the response
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=converted_file.%s", targetFormat))
	c.Data(http.StatusOK, "application/octet-stream", resp.Data.([]byte))
}
