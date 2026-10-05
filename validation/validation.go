package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var v = validator.New()

func init() {
	// backs the `password` tag used by the input DTOs.
	_ = v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) >= 8
	})

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name, _, _ := strings.Cut(fld.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		return name
	})
}

// Struct validates s against its `validate` tags and returns the validation
// error (nil when valid).
func Struct(s any) error { return v.Struct(s) }
