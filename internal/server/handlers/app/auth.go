package app

import (
	"fmt"
	"net/http"
	"server/internal/core/shared/types"
	"server/internal/services/core/auth"

	"github.com/gin-gonic/gin"
)

type (
	authLoginFormData     auth.LoginData          // @name AuthLoginData
	authLoginResponseData auth.LoginResultPayload // @name AuthLoginResponseData
	authLoginResponseBody struct {
		apiResponseBody
		Data *authLoginResponseData `json:"data,omitempty"`
	} // @name AuthLoginResponseBody
)

func (h *handler) authLogin(context *gin.Context) {
	var (
		requestBody  authLoginFormData
		responseBody authLoginResponseBody
	)
	if err := context.BindJSON(&requestBody); err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	ipAddress, err := h.helpers.Client().GetIPAddress(context)
	if err != nil {
		responseBody.Code = types.ServiceResultCodeInternal
		context.IndentedJSON(http.StatusInternalServerError, responseBody)
		h.ctx.Logger().Error(err)
		return
	}
	requestBody.IPAddress = ipAddress
	result, err := h.services.Auth().Login(auth.LoginData(requestBody))
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
		context.IndentedJSON(http.StatusUnauthorized, responseBody)
		return
	}
	responseBody.Code = result.Code
	responseBody.Data = (*authLoginResponseData)(result.Payload)
	context.Header("Authorization", fmt.Sprintf("Bearer %s", result.Payload.Token))
	context.IndentedJSON(http.StatusOK, responseBody)
}
