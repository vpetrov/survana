package db

import (
	"errors"
)

var (
	ErrNotFound          = errors.New("Not found")
	ErrInvalidId         = errors.New("Invalid id")
	ErrInvalidCollection = errors.New("Invalid collection")
)
