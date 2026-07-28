package util

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func Parse(ctx fiber.Ctx, req any) error {
	if err := ctx.Bind().Body(req); err != nil {
		return err
	}

	if err := validator.New().Struct(req); err != nil {
		return ParseValidationError(err)
	}

	return nil
}

func ParseValidationError(err error) error {
	if errs, ok := err.(validator.ValidationErrors); ok {
		e := errs[0]

		field := e.Field()
		param := e.Param()

		switch e.Tag() {
		case "required":
			return fmt.Errorf("%s is required", field)

		case "required_if":
			return fmt.Errorf("%s is required", field)

		case "required_without":
			return fmt.Errorf("%s is required", field)

		case "email":
			return fmt.Errorf("%s must be a valid email address", field)

		case "url":
			return fmt.Errorf("%s must be a valid URL", field)

		case "uri":
			return fmt.Errorf("%s must be a valid URI", field)

		case "uuid":
			return fmt.Errorf("%s must be a valid UUID", field)

		case "uuid4":
			return fmt.Errorf("%s must be a valid UUID v4", field)

		case "oneof":
			return fmt.Errorf("%s must be one of: %s", field, param)

		case "eq":
			return fmt.Errorf("%s must be equal to %s", field, param)

		case "ne":
			return fmt.Errorf("%s must not be equal to %s", field, param)

		case "gt":
			return fmt.Errorf("%s must be greater than %s", field, param)

		case "gte":
			return fmt.Errorf("%s must be greater than or equal to %s", field, param)

		case "lt":
			return fmt.Errorf("%s must be less than %s", field, param)

		case "lte":
			return fmt.Errorf("%s must be less than or equal to %s", field, param)

		case "min":
			switch e.Kind() {
			case reflect.String:
				return fmt.Errorf("%s must be at least %s characters", field, param)
			case reflect.Slice, reflect.Array:
				return fmt.Errorf("%s must contain at least %s items", field, param)
			default:
				return fmt.Errorf("%s must be at least %s", field, param)
			}

		case "max":
			switch e.Kind() {
			case reflect.String:
				return fmt.Errorf("%s must be at most %s characters", field, param)
			case reflect.Slice, reflect.Array:
				return fmt.Errorf("%s must contain at most %s items", field, param)
			default:
				return fmt.Errorf("%s must be at most %s", field, param)
			}

		case "len":
			switch e.Kind() {
			case reflect.String:
				return fmt.Errorf("%s must be exactly %s characters", field, param)
			case reflect.Slice, reflect.Array:
				return fmt.Errorf("%s must contain exactly %s items", field, param)
			default:
				return fmt.Errorf("%s must be exactly %s", field, param)
			}

		case "numeric":
			return fmt.Errorf("%s must contain only numeric characters", field)

		case "number":
			return fmt.Errorf("%s must be a valid number", field)

		case "alpha":
			return fmt.Errorf("%s must contain only letters", field)

		case "alphanum":
			return fmt.Errorf("%s must contain only letters and numbers", field)

		case "boolean":
			return fmt.Errorf("%s must be true or false", field)

		case "datetime":
			return fmt.Errorf("%s must be a valid datetime", field)

		case "startswith":
			return fmt.Errorf("%s must start with %s", field, param)

		case "endswith":
			return fmt.Errorf("%s must end with %s", field, param)

		case "contains":
			return fmt.Errorf("%s must contain %s", field, param)

		case "excludes":
			return fmt.Errorf("%s must not contain %s", field, param)

		default:
			return fmt.Errorf("%s is invalid", field)
		}
	}

	return err
}
