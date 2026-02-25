package ebay

import (
	"server/internal/core/shared/types"
)

type Provider struct {
	ctx types.IContext
}

func New(ctx types.IContext) *Provider {
	return &Provider{ctx: ctx}
}
