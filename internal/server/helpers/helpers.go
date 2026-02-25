package helpers

import (
	"server/internal/core/shared/types"
	"server/internal/services"
)

type Helper struct {
	sse     *SSEHelper
	auth    *AuthHelper
	client  *ClientHelper
	swagger *SwaggerHelper
}

func New(ctx types.IContext, services *services.Services) *Helper {
	return &Helper{
		sse: NewSSEHelper(ctx, 5000),
		auth: &AuthHelper{
			ctx:      ctx,
			services: services,
		},
		client: &ClientHelper{
			ctx:      ctx,
			services: services,
		},
		swagger: &SwaggerHelper{
			ctx:      ctx,
			services: services,
		},
	}
}

func (h *Helper) SSE() *SSEHelper {
	return h.sse
}

func (h *Helper) Auth() *AuthHelper {
	return h.auth
}

func (h *Helper) Client() *ClientHelper {
	return h.client
}

func (h *Helper) Swagger() *SwaggerHelper {
	return h.swagger
}
