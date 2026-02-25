package app

import (
	"net/http"
	"server/internal/core/shared/types"
	"server/internal/data/entities"
	coreusers "server/internal/services/core/users"

	"github.com/gin-gonic/gin"
)

type getAuthUserResponseBody struct {
	apiResponseBody
	Data struct {
		User entities.User `json:"user"`
	} `json:"data"`
}

type (
	getUsersRequestBody  coreusers.GetManyData
	getUsersResponseBody struct {
		apiResponseBody
		Data coreusers.GetManyResultPayload `json:"data"`
	}
	addUserRequestBody  coreusers.AddData
	addUserResponseBody struct {
		apiResponseBody
		Data coreusers.AddResultPayload `json:"data"`
	}
	updateUserRequestBody  coreusers.UpdateData
	updateUserResponseBody struct {
		apiResponseBody
		Data coreusers.UpdateResultPayload `json:"data"`
	}
	deleteUserResponseBody struct {
		apiResponseBody
	}
)

func (h *handler) getAuthUser(context *gin.Context) {
	var responseBody getAuthUserResponseBody
	user, auth := h.helpers.Auth().GetUser(context)
	if !auth || user == nil {
		responseBody.Code = types.ServiceResultCodeUnauthorized
		context.IndentedJSON(http.StatusUnauthorized, responseBody)
		return
	}
	responseBody.Code = types.ServiceResultCodeSuccess
	responseBody.Data = struct {
		User entities.User `json:"user"`
	}{User: *user}
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) updateAuthUser(context *gin.Context) {
}

func (h *handler) getUsers(context *gin.Context) {
	var (
		requestBody  getUsersRequestBody
		responseBody getUsersResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Users().GetMany(coreusers.GetManyData(requestBody))
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
	responseBody.Pagination = result.Pagination
	context.IndentedJSON(http.StatusOK, responseBody)
}

func (h *handler) addUser(context *gin.Context) {
	var (
		requestBody  addUserRequestBody
		responseBody addUserResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Users().Add(coreusers.AddData(requestBody))
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

func (h *handler) updateUser(context *gin.Context) {
	var (
		requestBody  updateUserRequestBody
		responseBody updateUserResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.ID = context.Param("id")
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Users().Update(coreusers.UpdateData(requestBody))
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

func (h *handler) deleteUser(context *gin.Context) {
	var responseBody deleteUserResponseBody
	result, err := h.services.Users().Delete(coreusers.DeleteData{
		LocationID: context.Param("lid"),
		ID:         context.Param("id"),
	})
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
	context.IndentedJSON(http.StatusOK, responseBody)
}
