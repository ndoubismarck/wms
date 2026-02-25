package app

import (
	"net/http"
	"server/internal/core/shared/types"
	"server/internal/services/core/setup"

	"github.com/gin-gonic/gin"
)

type getSetupResponseBody struct {
	apiResponseBody
	Data setup.GetResultPayload `json:"data"`
}

type createSetupResponseBody struct {
	apiResponseBody
	Data setup.CreateResultPayload `json:"data"`
}

func (h *handler) getSetup(context *gin.Context) {
	var (
		requestBody  setup.GetData
		responseBody getSetupResponseBody
	)
	user, auth := h.helpers.Auth().GetUser(context)
	if auth && user != nil {
		requestBody.User = user
	}
	result, err := h.services.Setup().Get(requestBody)
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if result.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = result.Code
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = result.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) createSetup(context *gin.Context) {
	var (
		requestBody  setup.CreateData
		responseBody createSetupResponseBody
	)
	checkResult, err := h.services.Setup().Get(setup.GetData{})
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	valid := true
	if !checkResult.Payload.Setup.Database {
		valid = false
	}
	if !checkResult.Payload.Setup.Admin {
		valid = false
	}
	if valid {
		//context.Status(http.StatusNotFound)
		//return
	}
	if err = context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	createResult, err := h.services.Setup().Create(requestBody)
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	if createResult.Code != types.ServiceResultCodeSuccess {
		responseBody.Code = createResult.Code
		if createResult.Code == types.ServiceResultCodeInvalid {
			responseBody.Validation = createResult.Validation
		}
		context.IndentedJSON(http.StatusForbidden, responseBody)
		return
	}
	responseBody.Code = createResult.Code
	responseBody.Data = createResult.Payload
	context.IndentedJSON(http.StatusOK, responseBody)
}
