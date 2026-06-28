package repository

import (
	"context"
	"database/sql"

	"github.com/SawitProRecruitment/UserService/models"
	"github.com/google/uuid"
)

func (r *PostgresRepository) CreateEstate(ctx context.Context, width int, length int) (*models.Estate, error) {

	estate := &models.Estate{
		ID:     uuid.New(),
		Width:  width,
		Length: length,
	}

	query := `
    INSERT INTO estates (
			id,
			width,
			length
		)
		VALUES ($1, $2, $3)
  `

	_, err := r.db.ExecContext(ctx, query, estate.ID, estate.Width, estate.Length)

	if err != nil {
		return nil, err
	}

	return estate, nil
}

func (r *PostgresRepository) GetEstate(ctx context.Context, id uuid.UUID) (*models.Estate, error) {

	var estate models.Estate

	query := `
    SELECT
			id,
			width,
			length
		FROM estates
		WHERE id = $1
  `

	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&estate.ID, &estate.Width, &estate.Length)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &estate, nil
}
