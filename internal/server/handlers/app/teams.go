package app

import (
	"net/http"
	"server/internal/core/shared/types"
	"server/internal/services/core/teams"

	"github.com/gin-gonic/gin"
)

type (
	getTeamsRequestBody  teams.GetManyData
	getTeamsResponseBody struct {
		apiResponseBody
		Data teams.GetManyResultPayload `json:"data"`
	}
	addTeamRequestBody  teams.AddData
	addTeamResponseBody struct {
		apiResponseBody
		Data teams.AddResultPayload `json:"data"`
	}
	updateTeamRequestBody  teams.UpdateData
	updateTeamResponseBody struct {
		apiResponseBody
		Data teams.UpdateResultPayload `json:"data"`
	}
	deleteTeamResponseBody struct {
		apiResponseBody
	}
	getTeamResponseBody struct {
		apiResponseBody
		Data teams.GetResultPayload `json:"data"`
	}
)

type (
	setTeamMembersRequestBody  teams.SetMembersData
	setTeamMembersResponseBody struct {
		apiResponseBody
		Data teams.SetMembersResultPayload `json:"data"`
	}
	getTeamMembersResponseBody struct {
		apiResponseBody
		Data teams.GetMembersResultPayload `json:"data"`
	}
)

func (h *handler) getTeams(context *gin.Context) {
	var (
		requestBody  getTeamsRequestBody
		responseBody getTeamsResponseBody
	)
	if err := context.BindQuery(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Teams().GetMany(teams.GetManyData(requestBody))
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

func (h *handler) addTeam(context *gin.Context) {
	var (
		requestBody  addTeamRequestBody
		responseBody addTeamResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Teams().Add(teams.AddData(requestBody))
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

func (h *handler) updateTeam(context *gin.Context) {
	var (
		requestBody  updateTeamRequestBody
		responseBody updateTeamResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.ID = context.Param("id")
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Teams().Update(teams.UpdateData(requestBody))
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

func (h *handler) deleteTeam(context *gin.Context) {
	var responseBody deleteTeamResponseBody
	result, err := h.services.Teams().Delete(teams.DeleteData{LocationID: context.Param("lid"), ID: context.Param("id")})
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

func (h *handler) getTeam(context *gin.Context) {
	var responseBody getTeamResponseBody
	id := context.Param("id")
	result, err := h.services.Teams().Get(teams.GetData{LocationID: context.Param("lid"), ID: &id})
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

func (h *handler) setTeamMembers(context *gin.Context) {
	var (
		requestBody  setTeamMembersRequestBody
		responseBody setTeamMembersResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.LocationID = context.Param("lid")
	requestBody.TeamID = context.Param("id")
	result, err := h.services.Teams().SetMembers(teams.SetMembersData(requestBody))
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

func (h *handler) getTeamMembers(context *gin.Context) {
	var responseBody getTeamMembersResponseBody
	result, err := h.services.Teams().GetMembers(teams.GetMembersData{LocationID: context.Param("lid"), TeamID: context.Param("id")})
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
