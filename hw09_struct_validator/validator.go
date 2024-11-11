package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

const validateTag = "validate"

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

func Validate(v interface{}) error {
	validator := NewValidator()

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("cannot validate non-struct %T", v)
	}

	rt := rv.Type()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		tag := field.Tag.Get(validateTag)
		if tag == "" {
			continue
		}

		for _, constrain := range strings.Split(tag, "|") {
			if constrain == "" {
				continue
			}

			switch {
			case strings.HasPrefix(constrain, "len:"):
				validator = validator.Len(constrain, rv.Field(i), rt.Field(i).Name)
			case strings.HasPrefix(constrain, "min:"):
				validator = validator.Min(constrain, rv.Field(i), rt.Field(i).Name)
			case strings.HasPrefix(constrain, "max:"):
				validator = validator.Max(constrain, rv.Field(i), rt.Field(i).Name)
			case strings.HasPrefix(constrain, "regexp:"):
				validator = validator.Regexp(constrain, rv.Field(i), rt.Field(i).Name)
			case strings.HasPrefix(constrain, "in:"):
				validator = validator.In(constrain, rv.Field(i), rt.Field(i).Name)
			}
		}
	}

	return validator.Run()
}
