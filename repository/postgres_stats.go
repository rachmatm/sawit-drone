package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/models"
	"github.com/google/uuid"
)

func (r *PostgresRepository) GetTreesByEstate(ctx context.Context, estateID uuid.UUID) ([]models.Tree, error) {

	query := `
		SELECT
			id,
			estate_id,
			x,
			y,
			height
		FROM trees
		WHERE estate_id = $1
		ORDER BY height
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		estateID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trees []models.Tree

	for rows.Next() {

		var tree models.Tree

		if err := rows.Scan(
			&tree.ID,
			&tree.EstateID,
			&tree.X,
			&tree.Y,
			&tree.Height,
		); err != nil {
			return nil, err
		}

		trees = append(trees, tree)
	}

	return trees, nil
}
