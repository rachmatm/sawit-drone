package models

import "github.com/google/uuid"

type Estate struct {
	ID     uuid.UUID
	Width  int
	Length int
	Trees  []Tree
}
