package router

import (
	"github.com/gin-gonic/gin"
	"synth.com/file_converter/internal/handler"
)

// NewRouter sets up the routes and returns a gin.Engine instance
func NewRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/convert", handler.ConvertFileHandler) // POST request for file conversion
	return r
}
