package models

import "github.com/google/uuid"

type Tree struct {
	ID       uuid.UUID
	EstateID uuid.UUID
	X        int
	Y        int
	Height   int
}
