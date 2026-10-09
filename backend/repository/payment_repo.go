package repository

import (
	"context"
	"database/sql"
	"backend/models"
)

type PaymentRepository interface {
	CreatePaymentStatusTx(ctx context.Context, tx *sql.Tx, p models.Payment) (*models.Payment, error)
	UpdateStatusTx(ctx context.Context, tx *sql.Tx, paymentID string, status models.PaymentState) error
}

type paymentRepo struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepo{db: db}
}

func (r *paymentRepo) CreatePaymentStatusTx(ctx context.Context, tx *sql.Tx, p models.Payment) (*models.Payment, error) {
	err := tx.QueryRowContext(ctx, "INSERT INTO payment_status(id, status, amount, created_at, updated_at) VALUES(gen_random_uuid(), $1, $2, NOW(), NOW()) RETURNING id, created_at, updated_at", p.Status, p.Amount).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepo) UpdateStatusTx(ctx context.Context, tx *sql.Tx, paymentID string, status models.PaymentState) error {
	_, err := tx.ExecContext(ctx, "UPDATE payment_status SET status = $1, updated_at = NOW() WHERE id = $2", status, paymentID)
	return err
}
