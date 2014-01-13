package survana

import (
	"errors"
)

var (
	ErrEmptyRequest      = errors.New("Empty request")
	ErrNotFound          = errors.New("Not found")
	ErrInvalidId         = errors.New("Invalid id")
	ErrInvalidCollection = errors.New("Invalid collection")
	ErrUnauthorized      = errors.New("Unauthorized request")
    ErrNoSuchForm        = errors.New("Form does not exist")
)
