package utils

import "fmt"

// Custom error for file conversion
type FileConversionError struct {
	Msg string
}

func (e *FileConversionError) Error() string {
	return fmt.Sprintf("File Conversion Error: %s", e.Msg)
}
