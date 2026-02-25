package app

import (
	"net/http"
	"server/internal/core/shared/types"
	"server/internal/services/core/files"

	"github.com/gin-gonic/gin"
)

type (
	uploadFileRequestBody  files.UploadData
	uploadFileResponseBody struct {
		apiResponseBody
		Data files.UploadResultPayload `json:"data"`
	}
)

func (h *handler) uploadFile(context *gin.Context) {
	var (
		requestBody  uploadFileRequestBody
		responseBody uploadFileResponseBody
	)
	if err := context.ShouldBind(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	result, err := h.services.Files().Upload(files.UploadData(requestBody))
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		if result.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = result.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}
