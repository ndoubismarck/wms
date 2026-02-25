package types

type ServiceResultCode string

const (
	ServiceResultCodeSuccess      ServiceResultCode = "success"
	ServiceResultCodeNotFound     ServiceResultCode = "error.notfound"
	ServiceResultCodeFailed       ServiceResultCode = "error.failed"
	ServiceResultCodeInvalid      ServiceResultCode = "error.invalid"
	ServiceResultCodeInternal     ServiceResultCode = "error.internal"
	ServiceResultCodeForbidden    ServiceResultCode = "error.forbidden"
	ServiceResultCodeUnauthorized ServiceResultCode = "error.unauthorized"
	ServiceResultCodeInvalidSetup ServiceResultCode = "error.invalid.setup"
)

type ValidationResult map[string][]string

type IPagination interface {
	GetOrder() string
	GetLimit() int
	GetOffset() int
	GetResult(totalRecords int64) PaginationResult
	GetOrderWithPrefix(prefix string) string
}

type PaginationParams struct {
	Page  uint32 `json:"page" form:"page" validate:"required" example:"1" default:"1" minimum:"1"`     // Current page
	Limit uint32 `json:"limit" form:"limit" validate:"optional" example:"50" default:"50" minimum:"1"` // Max results per page
	Order string `json:"order" form:"order" validate:"optional" example:"id.desc" default:"id.desc"`   // Results order
}

type PaginationResult struct {
	Page         uint32 `json:"page"`
	Prev         uint32 `json:"prev"`
	Next         uint32 `json:"next"`
	Limit        uint32 `json:"limit"`
	Order        string `json:"order"`
	HasPrev      bool   `json:"has_prev"`
	HasNext      bool   `json:"has_next"`
	TotalPages   int64  `json:"total_pages"`
	TotalRecords int64  `json:"total_records"`
}
