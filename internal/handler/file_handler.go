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
	convertedFile, err := service.ConvertFile(file, targetFormat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Conversion failed: %s", err.Error())})
		return
	}

	// Send the converted file as the response
	c.Header("Content-Disposition", "attachment; filename=converted_file."+targetFormat)
	c.Data(http.StatusOK, "application/octet-stream", convertedFile)
}
