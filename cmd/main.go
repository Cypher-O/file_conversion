package main

import (
	"github.com/gin-gonic/gin"
	"synth.com/file_converter/internal/service"
	"synth.com/file_converter/internal/response"
	"net/http"
	"log"
)

// UploadHandler is the handler for file upload and conversion
func UploadHandler(c *gin.Context) {
	// Get file from the form
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(400, "No file uploaded"))
		return
	}

	// Get the target format from the query params
	targetFormat := c.DefaultQuery("format", "")

	if targetFormat == "" {
		c.JSON(http.StatusBadRequest, response.NewErrorResponse(400, "Format query parameter is required"))
		return
	}

	// Call the service to process the file
	resp := service.ConvertFile(file, targetFormat)

	// Return the structured response
	c.JSON(http.StatusOK, resp)
}

func main() {
	r := gin.Default()

	// Define the POST endpoint for file conversion
	r.POST("/convert", UploadHandler)

	// Start the server
	err := r.Run(":8080") // Run on port 8080
	if err != nil {
		log.Fatal("Unable to start the server:", err)
	}
}
