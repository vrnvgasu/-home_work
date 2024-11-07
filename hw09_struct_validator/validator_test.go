package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr ValidationErrors
	}{
		{
			in: User{
				ID:     "3422b448-2460-4fd2-9183-8000de6f8343",
				Name:   "test",
				Age:    18,
				Email:  "test@test.com",
				Role:   "stuff",
				Phones: []string{"79001001010", "79002002020"},
				meta:   json.RawMessage(`{"number":123}`),
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "1",
				Name:   "test",
				Age:    10,
				Email:  "test@testcom",
				Role:   "dummy",
				Phones: []string{"555", "222"},
				meta:   json.RawMessage(nil),
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   errLen,
				},
				ValidationError{
					Field: "Age",
					Err:   errMin,
				},
				ValidationError{
					Field: "Email",
					Err:   errRegexp,
				},
				ValidationError{
					Field: "Role",
					Err:   errIn,
				},
				ValidationError{
					Field: "Phones",
					Err:   errLen,
				},
				ValidationError{
					Field: "Phones",
					Err:   errLen,
				},
			},
		},
		{
			in: App{
				Version: "12345",
			},
			expectedErr: nil,
		},
		{
			in: App{
				Version: "",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Version",
					Err:   errLen,
				},
			},
		},
		{
			in:          Token{},
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte("Header"),
				Payload:   []byte("Payload"),
				Signature: []byte("Signature"),
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 404,
				Body: "not found",
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 201,
				Body: "",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   errIn,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)

				return
			}

			var vErrs ValidationErrors
			ok := errors.As(err, &vErrs)
			require.True(t, ok)

			slices.SortFunc(vErrs, func(a, b ValidationError) int {
				if a.Field > b.Field {
					return 1
				}

				return -1
			})
			slices.SortFunc(tt.expectedErr, func(a, b ValidationError) int {
				if a.Field > b.Field {
					return 1
				}

				return -1
			})

			require.Equal(t, tt.expectedErr, vErrs)
		})
	}
}

func TestValidateLangError(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr require.ErrorAssertionFunc
	}{
		{
			in: struct {
				ID     string `json:"id" validate:"len:36"`
				Name   string
				Age    int      `validate:"min:18|max:50"`
				Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
				Role   UserRole `validate:"in:admin,stuff"`
				Phones []string `validate:"len:11"`
			}{
				ID:     "3422b448-2460-4fd2-9183-8000de6f8343",
				Name:   "test",
				Age:    18,
				Email:  "test@test.com",
				Role:   "stuff",
				Phones: []string{"79001001010", "79002002020"},
			},
			expectedErr: require.NoError,
		},
		{
			in: struct {
				ID     int `json:"id" validate:"len:36"`
				Name   string
				Age    int      `validate:"min:18|max:50"`
				Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
				Role   UserRole `validate:"in:admin,stuff"`
				Phones []string `validate:"len:11"`
			}{
				ID:     123,
				Name:   "test",
				Age:    18,
				Email:  "test@test.com",
				Role:   "stuff",
				Phones: []string{"79001001010", "79002002020"},
			},
			expectedErr: require.Error,
		},
		{
			in: struct {
				ID     string `json:"id" validate:"len:36"`
				Name   string
				Age    string   `validate:"min:18|max:50"`
				Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
				Role   UserRole `validate:"in:admin,stuff"`
				Phones []string `validate:"len:11"`
			}{
				ID:     "3422b448-2460-4fd2-9183-8000de6f8343",
				Name:   "test",
				Age:    "18",
				Email:  "test@test.com",
				Role:   "stuff",
				Phones: []string{"79001001010", "79002002020"},
			},
			expectedErr: require.Error,
		},
		{
			in: struct {
				ID     string `json:"id" validate:"len:36"`
				Name   string
				Age    int      `validate:"min:18|max:50"`
				Email  int      `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
				Role   UserRole `validate:"in:admin,stuff"`
				Phones []string `validate:"len:11"`
			}{
				ID:     "3422b448-2460-4fd2-9183-8000de6f8343",
				Name:   "test",
				Age:    18,
				Email:  99,
				Role:   "stuff",
				Phones: []string{"79001001010", "79002002020"},
			},
			expectedErr: require.Error,
		},
		{
			in: struct {
				ID     string `json:"id" validate:"len:36"`
				Name   string
				Age    int      `validate:"min:18|max:50"`
				Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
				Role   bool     `validate:"in:admin,stuff"`
				Phones []string `validate:"len:11"`
			}{
				ID:     "3422b448-2460-4fd2-9183-8000de6f8343",
				Name:   "test",
				Age:    18,
				Email:  "test@test.com",
				Role:   true,
				Phones: []string{"79001001010", "79002002020"},
			},
			expectedErr: require.Error,
		},
		{
			in: struct {
				ID     string `json:"id" validate:"len:36"`
				Name   string
				Age    int      `validate:"min:18|max:50"`
				Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
				Role   UserRole `validate:"in:admin,stuff"`
				Phones []int    `validate:"len:11"`
			}{
				ID:     "3422b448-2460-4fd2-9183-8000de6f8343",
				Name:   "test",
				Age:    18,
				Email:  "test@test.com",
				Role:   "stuff",
				Phones: []int{555, 777},
			},
			expectedErr: require.Error,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			tt.expectedErr(t, err)

			var vErrs ValidationErrors
			require.False(t, errors.As(err, &vErrs))
		})
	}
}
