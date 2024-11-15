package storage

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrDateBusy = errors.New("date busy")
)
