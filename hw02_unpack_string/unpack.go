package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (unpack string, err error) {
	if err = validate(s); err != nil {
		return "", err
	}

	result := make([]rune, 0, len(s))
	for _, c := range s {
		if result, err = add(result, c); err != nil {
			return "", err
		}
	}

	return string(result), nil
}

func validate(s string) error {
	var lastChar rune
	for i, c := range s {
		if i == 0 && unicode.IsNumber(c) {
			return ErrInvalidString
		}
		if unicode.IsNumber(c) && unicode.IsNumber(lastChar) {
			return ErrInvalidString
		}

		lastChar = c
	}

	return nil
}

func add(chars []rune, c rune) ([]rune, error) {
	if !unicode.IsNumber(c) {
		return append(chars, c), nil
	}

	repeat, err := strconv.Atoi(string(c))
	if err != nil {
		return nil, err
	}

	if repeat == 0 {
		newResult := make([]rune, len(chars)-1)
		copy(newResult, chars)

		return newResult, nil
	}

	prevC := chars[len(chars)-1]
	for j := 1; j < repeat; j++ {
		chars = append(chars, prevC)
	}

	return chars, nil
}
