package pagination

import (
	"fmt"
	"server/internal/core/shared/types"
	"strings"
)

type pagination struct {
	ctx    types.IContext
	params types.PaginationParams
}

func New(params types.PaginationParams) types.IPagination {
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 100
	}
	if len(params.Order) == 0 {
		params.Order = "created_at.desc"
	}
	return &pagination{
		params: params,
	}
}

func (p *pagination) GetOrder() string {
	if len(p.params.Order) == 0 {
		p.params.Order = "created_at.desc"
	}
	return strings.ReplaceAll(p.params.Order, ".", " ")
}

func (p *pagination) GetLimit() int {
	if p.params.Limit == 0 {
		p.params.Limit = 100
	}
	return int(p.params.Limit)
}

func (p *pagination) GetOffset() int {
	if p.params.Page == 0 {
		p.params.Page = 1
	}
	if p.params.Limit == 0 {
		p.params.Limit = 100
	}
	page := p.params.Page
	limit := p.params.Limit
	return int((page - 1) * limit)
}

func (p *pagination) GetResult(totalRecords int64) types.PaginationResult {
	if p.params.Page == 0 {
		p.params.Page = 1
	}
	nextPage := p.params.Page
	prevPage := p.params.Page
	currentPage := p.params.Page
	totalPages := (totalRecords + int64(p.params.Limit) - 1) / int64(p.params.Limit)
	hasPrev := currentPage > 1
	hasNext := int64(nextPage) < totalPages
	if hasPrev {
		prevPage = currentPage - 1
	}
	if hasNext {
		nextPage = currentPage + 1
	}
	return types.PaginationResult{
		Page:         currentPage,
		Prev:         prevPage,
		Next:         nextPage,
		Limit:        p.params.Limit,
		Order:        p.params.Order,
		HasNext:      hasNext,
		HasPrev:      hasPrev,
		TotalPages:   totalPages,
		TotalRecords: totalRecords,
	}
}

func (p *pagination) GetOrderWithPrefix(prefix string) string {
	return fmt.Sprintf("%s.%s", prefix, p.GetOrder())
}
