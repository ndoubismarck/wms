package files

import "mime/multipart"

const UploadsDir = "tmp/uploads"

type (
	UploadFileData struct {
		File *multipart.FileHeader `json:"file"`
	}
	UploadFileResult struct {
		Size int64  `json:"size"`
		Type string `json:"type"`
		Path string `json:"path"`
		Name string `json:"name"`
	}
)
