// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/models"
	"github.com/google/uuid"
)

type EstateRepository interface {
	CreateEstate(
		ctx context.Context,
		width int,
		length int,
	) (*models.Estate, error)

	GetEstate(
		ctx context.Context,
		id uuid.UUID,
	) (*models.Estate, error)

	CreateTree(
		ctx context.Context,
		estateID uuid.UUID,
		x int,
		y int,
		height int,
	) (*models.Tree, error)

	GetTreesByEstate(
		ctx context.Context,
		estateID uuid.UUID,
	) ([]models.Tree, error)
}
