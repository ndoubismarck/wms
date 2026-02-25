package app

import (
	"net/http"
	"server/internal/core/shared/types"
	"server/internal/services/core/stats"

	"github.com/gin-gonic/gin"
)

type getStatsResponseBody struct {
	apiResponseBody
	Data stats.GetResultPayload `json:"data"`
}

func (h *handler) getStats(context *gin.Context) {
	var (
		requestBody  stats.GetData
		responseBody getStatsResponseBody
	)
	account, auth := h.helpers.Auth().GetUser(context)
	if auth && account != nil {
		requestBody.UserID = account.ID
	}
	requestBody.LocationID = context.Param("lid")
	result, err := h.services.Stats().Get(requestBody)
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
