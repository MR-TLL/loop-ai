package repository

import (
	"context"
	"database/sql"
	"backend/models"
)

type ClinicRepository interface {
	GetByID(ctx context.Context, id string) (*models.Clinic, error)
	ListAll(ctx context.Context) ([]models.Clinic, error)
}

type clinicRepo struct {
	db *sql.DB
}

func NewClinicRepository(db *sql.DB) ClinicRepository {
	return &clinicRepo{db: db}
}

func (r *clinicRepo) GetByID(ctx context.Context, id string) (*models.Clinic, error) {
	var c models.Clinic
	err := r.db.QueryRowContext(ctx, "SELECT id, name, address, created_at, updated_at FROM clinics WHERE id = $1", id).Scan(&c.ID, &c.Name, &c.Address, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &c, nil
}

func (r *clinicRepo) ListAll(ctx context.Context) ([]models.Clinic, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, address, created_at, updated_at FROM clinics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clinics []models.Clinic
	for rows.Next() {
		var c models.Clinic
		if err := rows.Scan(&c.ID, &c.Name, &c.Address, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		clinics = append(clinics, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return clinics, nil
}
