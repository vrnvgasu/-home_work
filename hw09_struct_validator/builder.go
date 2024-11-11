package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	errLen    = errors.New("length is failed")
	errMin    = errors.New("min is failed")
	errMax    = errors.New("max is failed")
	errRegexp = errors.New("regexp is failed")
	errIn     = errors.New("in is failed")
)

type stage func() error

type Validator struct {
	m       sync.Mutex
	vErrors ValidationErrors
	stages  []stage
}

func NewValidator() *Validator {
	return &Validator{
		vErrors: ValidationErrors{},
		stages:  []stage{},
	}
}

func (v *Validator) addValErr(vErr ValidationError) {
	v.m.Lock()
	defer v.m.Unlock()

	v.vErrors = append(v.vErrors, vErr)
}

func (v *Validator) Run() error {
	var eg errgroup.Group

	for _, s := range v.stages {
		eg.Go(func() error {
			return s()
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	if len(v.vErrors) == 0 {
		return nil
	}

	return v.vErrors
}

func (v *Validator) In(constrain string, vr reflect.Value, fieldName string) *Validator {
	f := func() error {
		rule, err := constrainValue(constrain)
		if err != nil {
			return err
		}

		ins := strings.Split(rule, ",")

		for _, in := range ins {
			switch vr.Kind() { //nolint:exhaustive
			case reflect.String:
				if vr.String() == in {
					return nil
				}
			case reflect.Int:
				inInt, err := strconv.Atoi(in)
				if err != nil {
					return fmt.Errorf("invalid constrain in: %s: %w", constrain, err)
				}
				if int(vr.Int()) == inInt {
					return nil
				}
			default:
				return fmt.Errorf("field is not string or int for in validation: %s", fieldName)
			}
		}

		v.addValErr(ValidationError{
			Field: fieldName,
			Err:   errIn,
		})

		return nil
	}
	v.stages = append(v.stages, f)

	return v
}

func (v *Validator) Regexp(constrain string, vr reflect.Value, fieldName string) *Validator {
	f := func() error {
		if vr.Kind() != reflect.String {
			return fmt.Errorf("field is not string for regexp validation: %s", fieldName)
		}

		rule, err := constrainValue(constrain)
		if err != nil {
			return err
		}

		ok, err := regexp.MatchString(rule, vr.String())
		if err != nil {
			return fmt.Errorf("invalid regexp rule: %s", fieldName)
		}

		if !ok {
			v.addValErr(ValidationError{
				Field: fieldName,
				Err:   errRegexp,
			})
		}

		return nil
	}
	v.stages = append(v.stages, f)

	return v
}

func (v *Validator) Max(constrain string, vr reflect.Value, fieldName string) *Validator {
	f := func() error {
		if vr.Kind() != reflect.Int && vr.Kind() != reflect.Int32 && vr.Kind() != reflect.Int64 {
			return fmt.Errorf("field is not int for max validation: %s", fieldName)
		}

		rule, err := constrainValue(constrain)
		if err != nil {
			return err
		}

		maxV, err := strconv.Atoi(rule)
		if err != nil {
			return fmt.Errorf("invalid constrain max: %s: %w", constrain, err)
		}

		if int(vr.Int()) > maxV {
			v.addValErr(ValidationError{
				Field: fieldName,
				Err:   errMax,
			})
		}

		return nil
	}
	v.stages = append(v.stages, f)

	return v
}

func (v *Validator) Min(constrain string, vr reflect.Value, fieldName string) *Validator {
	f := func() error {
		if vr.Kind() != reflect.Int && vr.Kind() != reflect.Int32 && vr.Kind() != reflect.Int64 {
			return fmt.Errorf("field is not int for min validation: %s", fieldName)
		}

		rule, err := constrainValue(constrain)
		if err != nil {
			return err
		}

		minV, err := strconv.Atoi(rule)
		if err != nil {
			return fmt.Errorf("invalid constrain min: %s: %w", constrain, err)
		}

		if int(vr.Int()) < minV {
			v.addValErr(ValidationError{
				Field: fieldName,
				Err:   errMin,
			})
		}

		return nil
	}
	v.stages = append(v.stages, f)

	return v
}

func (v *Validator) Len(constrain string, vr reflect.Value, fieldName string) *Validator {
	if vr.Kind() == reflect.Slice {
		for i := 0; i < vr.Len(); i++ {
			v.lenValidation(constrain, vr.Index(i), fieldName)
		}
	} else {
		v.lenValidation(constrain, vr, fieldName)
	}

	return v
}

func (v *Validator) lenValidation(constrain string, vr reflect.Value, fieldName string) {
	f := func() error {
		if vr.Kind() != reflect.String {
			return fmt.Errorf("field is not string for len validation: %s", fieldName)
		}

		rule, err := constrainValue(constrain)
		if err != nil {
			return err
		}

		length, err := strconv.Atoi(rule)
		if err != nil {
			return fmt.Errorf("invalid constrain len: %s: %w", constrain, err)
		}

		if len(vr.String()) != length {
			v.addValErr(ValidationError{
				Field: fieldName,
				Err:   errLen,
			})
		}

		return nil
	}
	v.stages = append(v.stages, f)
}

func constrainValue(constrain string) (string, error) {
	args := strings.Split(constrain, ":")
	if len(args) != 2 {
		return "", fmt.Errorf("invalid constrain: %s", constrain)
	}
	if args[1] == "" {
		return "", fmt.Errorf("invalid constrain value: %s", constrain)
	}

	return args[1], nil
}
