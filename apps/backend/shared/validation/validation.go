package validation

import (
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
)

var transformer = func() *mold.Transformer {
	return modifiers.New()
}()

var validate = func() *validator.Validate {
	v := validator.New()
	if err := v.RegisterValidation("emailmx", ValidateEmailMX); err != nil {
		panic(fmt.Sprintf("failed to register emailmx validation: %v", err))
	}
	return v
}()

// TransformAndValidate cleans the fields of request tagged with mold modifiers,
// then validates the fields tagged with validate rules. It returns a cleaned
// copy of request (the original is left untouched), or the first validation
// error.
func TransformAndValidate(request any) (any, error) {
	value := reflect.ValueOf(request)
	polished := reflect.New(value.Type())
	polished.Elem().Set(value)

	if err := transformer.Struct(context.Background(), polished.Interface()); err != nil {
		return request, errors.New("failed to process request")
	}

	if err := validate.Struct(polished.Interface()); err != nil {
		var validation_errors validator.ValidationErrors
		if errors.As(err, &validation_errors) && len(validation_errors) > 0 {
			return request, errors.New(validation_errors[0].Error())
		}
		return request, errors.New("failed to process request")
	}

	return polished.Elem().Interface(), nil
}

func ValidateEmailMX(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]

	resolver := net.Resolver{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mx, err := resolver.LookupMX(ctx, domain)
	if err == nil && len(mx) > 0 {
		return true
	}

	ips, err := resolver.LookupIPAddr(ctx, domain)
	if err == nil && len(ips) > 0 {
		return true
	}

	return false
}
