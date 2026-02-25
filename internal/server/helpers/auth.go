package helpers

import (
	"server/internal/core/shared/types"
	"server/internal/data/entities"
	"server/internal/services"
	"server/internal/services/core/auth"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHelper struct {
	ctx      types.IContext
	services *services.Services
}

type AuthType string

const authKey = "authAccount"
const (
	AuthTypePrivateApi AuthType = "private"
)

func (h *AuthHelper) Login(context *gin.Context, authType AuthType) (bool, error) {
	switch authType {
	case AuthTypePrivateApi:
		token := ""
		header := context.GetHeader("Authorization")
		if len(header) > 0 {
			token = strings.TrimPrefix(header, "Bearer")
			token = strings.TrimSpace(token)
		}
		if len(token) == 0 {
			return false, nil
		}
		result, err := h.services.Auth().JWTAuth(auth.JWTAuthData{
			Token: token,
		})
		if err != nil {
			return false, err
		}
		if result.Code != types.ServiceResultCodeSuccess {
			return false, nil
		}
		if result.Payload == nil {
			return false, nil
		}
		if len(result.Payload.User.ID) == 0 {
			return false, nil
		}
		context.Set(authKey, result.Payload.User)
	}
	return true, nil
}

func (h *AuthHelper) GetUser(context *gin.Context) (*entities.User, bool) {
	val, ok := context.Get(authKey)
	if !ok {
		ok, err := h.Login(context, AuthTypePrivateApi)
		if err != nil {
			h.ctx.Logger().Error(err)
			return nil, false
		}
		val, ok = context.Get(authKey)
		if !ok {
			return nil, false
		}
	}
	account, ok := val.(entities.User)
	if !ok {
		return nil, false
	}
	return &account, true
}
