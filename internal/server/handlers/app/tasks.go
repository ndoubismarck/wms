package app

import (
	"net/http"
	"server/internal/core/shared/types"
	"server/internal/services/core/tasks"

	"github.com/gin-gonic/gin"
)

type (
	getTasksRequestBody  tasks.GetManyData
	getTasksResponseBody struct {
		apiResponseBody
		Data tasks.GetManyResultPayload `json:"data"`
	}
	addTaskRequestBody  tasks.AddData
	addTaskResponseBody struct {
		apiResponseBody
		Data tasks.AddResultPayload `json:"data"`
	}
	updateTaskRequestBody  tasks.UpdateData
	updateTaskResponseBody struct {
		apiResponseBody
		Data tasks.UpdateResultPayload `json:"data"`
	}
	deleteTaskResponseBody struct {
		apiResponseBody
	}
	getTaskResponseBody struct {
		apiResponseBody
		Data tasks.GetResultPayload `json:"data"`
	}
)

func (h *handler) getTasks(context *gin.Context) {
	var (
		requestBody  getTasksRequestBody
		responseBody getTasksResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Tasks().GetMany(tasks.GetManyData(requestBody))
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

func (h *handler) addTask(context *gin.Context) {
	var (
		requestBody  addTaskRequestBody
		responseBody addTaskResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	user, ok := h.helpers.Auth().GetUser(context)
	if !ok || user == nil {
		responseBody.Code = types.ServiceResultCodeUnauthorized
		context.IndentedJSON(http.StatusUnauthorized, responseBody)
		return
	}
	requestBody.LocationID = context.Param("lid")
	requestBody.ScheduledByUserID = &user.ID
	result, err := h.services.Tasks().Add(tasks.AddData(requestBody))
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

func (h *handler) updateTask(context *gin.Context) {
	var (
		requestBody  updateTaskRequestBody
		responseBody updateTaskResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	user, ok := h.helpers.Auth().GetUser(context)
	if !ok || user == nil {
		responseBody.Code = types.ServiceResultCodeUnauthorized
		context.IndentedJSON(http.StatusUnauthorized, responseBody)
		return
	}
	requestBody.ID = context.Param("id")
	requestBody.LocationID = context.Param("lid")
	requestBody.ScheduledByUserID = &user.ID
	result, err := h.services.Tasks().Update(tasks.UpdateData(requestBody))
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

func (h *handler) deleteTask(context *gin.Context) {
	var responseBody deleteTaskResponseBody
	result, err := h.services.Tasks().Delete(tasks.DeleteData{LocationID: context.Param("lid"), ID: context.Param("id")})
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

func (h *handler) getTask(context *gin.Context) {
	var responseBody getTaskResponseBody
	id := context.Param("id")
	result, err := h.services.Tasks().Get(tasks.GetData{LocationID: context.Param("lid"), ID: &id})
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
