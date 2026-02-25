package providers

import (
	"server/internal/core/shared/types"
	"server/internal/services/providers/database"
	"server/internal/services/providers/ebay"
	"server/internal/services/providers/files"
	"server/internal/services/providers/validation"
)

type Providers struct {
	ebay       *ebay.Provider
	files      *files.Provider
	database   *database.Provider
	validation *validation.Provider
}

func New(ctx types.IContext) *Providers {
	return &Providers{
		ebay:       ebay.New(ctx),
		files:      files.New(ctx),
		database:   database.New(ctx),
		validation: validation.New(ctx),
	}
}

func (p *Providers) Ebay() *ebay.Provider {
	return p.ebay
}

func (p *Providers) Files() *files.Provider {
	return p.files
}

func (p *Providers) Database() *database.Provider {
	return p.database
}

func (p *Providers) Validation() *validation.Provider {
	return p.validation
}
