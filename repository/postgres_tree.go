package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/SawitProRecruitment/UserService/models"
)

func (r *PostgresRepository) CreateTree(
	ctx context.Context,
	estateID uuid.UUID,
	x int,
	y int,
	height int,
) (*models.Tree, error) {

	estate, err := r.GetEstate(ctx, estateID)
	if err != nil {
		return nil, err
	}

	if estate == nil {
		return nil, ErrEstateNotFound
	}

	if x < 1 || x > estate.Width {
		return nil, ErrInvalidCoordinate
	}

	if y < 1 || y > estate.Length {
		return nil, ErrInvalidCoordinate
	}

	tree := &models.Tree{
		ID:       uuid.New(),
		EstateID: estateID,
		X:        x,
		Y:        y,
		Height:   height,
	}

	query := `
		INSERT INTO trees (
			id,
			estate_id,
			x,
			y,
			height
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	if _, err := r.db.ExecContext(
		ctx,
		query,
		tree.ID,
		tree.EstateID,
		tree.X,
		tree.Y,
		tree.Height,
	); err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" &&
				pqErr.Constraint == "idx_unique_tree_plot" {

				return nil, ErrTreeAlreadyExists
			}
		}

		return nil, err

	}

	return tree, nil
}
