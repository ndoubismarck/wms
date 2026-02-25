package validation

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"server/internal/core/shared/types"
	"strings"
	"time"
)

type validators struct {
	ctx types.IContext
}

func newValidators(ctx types.IContext) *validators {
	return &validators{ctx: ctx}
}

func (v *validators) date(fl validator.FieldLevel) bool {
	if _, err := time.Parse("2006-01-02", fl.Field().String()); err != nil {
		return false
	}
	return true
}

func (v *validators) datetime(fl validator.FieldLevel) bool {
	if _, err := time.Parse("2006-01-02 15:04", fl.Field().String()); err != nil {
		return false
	}
	return true
}

func (v *validators) password(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) >= 6
}

func (v *validators) notBlank(fl validator.FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		return len(strings.Trim(strings.TrimSpace(field.String()), "\x1c\x1d\x1e\x1f")) > 0
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}
