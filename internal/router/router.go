package router

import (
	"github.com/gin-gonic/gin"
	"synth.com/file_converter/internal/handler"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/convert", handler.ConvertFileHandler) // POST request for file conversion
	return r
}
