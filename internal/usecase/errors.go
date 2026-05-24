package usecase

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrAlreadyExists       = errors.New("already exists")
	ErrPermissionDenied    = errors.New("permission denied")
	ErrInvalidInput        = errors.New("invalid input")
	ErrEncounterDuplicated = errors.New("encounter duplicated")
	ErrUnsupportedProvider = errors.New("unsupported provider")
)
