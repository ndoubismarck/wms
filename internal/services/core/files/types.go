package files

import (
	"mime/multipart"
	"server/internal/core/shared/types"
	"server/internal/services/providers/files"
)

type (
	UploadData struct {
		File *multipart.FileHeader `form:"file" binding:"required"`
	}
	UploadResult struct {
		Code       types.ServiceResultCode
		Payload    UploadResultPayload
		Validation types.ValidationResult
	}

	UploadResultPayload struct {
		File files.UploadFileResult `json:"file"`
	}
)
