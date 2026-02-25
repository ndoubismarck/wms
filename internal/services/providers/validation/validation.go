package validation

import (
	"errors"
	"fmt"
	"server/internal/core/shared/types"
	"server/internal/core/shared/utils/stringutil"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ent "github.com/go-playground/validator/v10/translations/en"
)

type Provider struct {
	ctx          types.IContext
	validators   map[string]validator.Func
	translations map[string]string
}

func New(ctx types.IContext) *Provider {
	val := newValidators(ctx)
	return &Provider{
		ctx: ctx,
		validators: map[string]validator.Func{
			"date":     val.date,
			"datetime": val.datetime,
			"password": val.password,
			"notblank": val.notBlank,
		},
		translations: map[string]string{
			"e164":     "Must be a valid phone number",
			"email":    "Must be a valid email address",
			"notblank": "This field cannot be empty, blank or null.",
			"password": "Password must contain at least 3 letters, 2 numbers and 1 special character.",
			"datetime": "Invalid datetime format, expected YYYY-MM-DD HH:MM.",
		},
	}
}

func (p *Provider) ValidateStruct(value any) (types.ValidationResult, bool, error) {
	validate := validator.New()
	return p.validate(validate, value)
}

func (p *Provider) ValidateWithMapRules(value any, rules map[string]string, types any) (types.ValidationResult, bool, error) {
	validate := validator.New()
	validate.RegisterStructValidationMapRules(rules, types)
	return p.validate(validate, value)
}

func (p *Provider) validate(validate *validator.Validate, value any) (types.ValidationResult, bool, error) {
	eng := en.New()
	unt := ut.New(eng, eng)
	translator, ok := unt.GetTranslator("en")
	if !ok {
		return nil, false, errors.New("validation translator not found")
	}
	for key, val := range p.validators {
		if err := validate.RegisterValidation(key, val); err != nil {
			return nil, false, err
		}
	}
	if err := ent.RegisterDefaultTranslations(validate, translator); err != nil {
		return nil, false, err
	}
	valid := true
	var result types.ValidationResult
	if err := validate.Struct(value); err != nil {
		for _, val := range err.(validator.ValidationErrors) {
			field := stringutil.ToSnakeCase(val.StructField())
			message := stringutil.SplitCamelCase(val.Translate(translator))
			if !strings.HasSuffix(message, ".") {
				message = fmt.Sprintf("%s.", message)
			}
			if _, ok := result[field]; ok {
				result[field] = append(result[field], message)
			} else {
				result[field] = []string{
					message,
				}
			}
			valid = false
		}
	}
	return result, valid, nil
}

func (p *Provider) translateFunc(ut ut.Translator, fe validator.FieldError) string {
	t, err := ut.T(fe.Tag(), fe.Field())
	if err != nil {
		p.ctx.Logger().Error(err)
		return fe.(error).Error()
	}
	return t
}

func (p *Provider) registrationFunc(tag string, translation string) validator.RegisterTranslationsFunc {
	return func(ut ut.Translator) (err error) {
		if err = ut.Add(tag, translation, true); err != nil {
			p.ctx.Logger().Error(err)
			return
		}
		return
	}
}
