package model

import "mime/multipart"

// File represents an uploaded file's metadata and contents.
type File struct {
	Filename   string
	Filetype   string
	Size       int64
	TargetType string
	FileContent *multipart.FileHeader
}
