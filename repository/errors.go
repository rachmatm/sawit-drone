package repository

import "errors"

var (
	ErrEstateNotFound    = errors.New("estate not found")
	ErrInvalidCoordinate = errors.New("invalid coordinate")
	ErrTreeAlreadyExists = errors.New("tree already exists")
)
