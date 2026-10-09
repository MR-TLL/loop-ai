package repository

import (
	"context"
	"database/sql"
	"backend/models"
)

type ReservationRepository interface {
	CreateTx(ctx context.Context, tx *sql.Tx, r models.Reservation) (*models.Reservation, error)
	GetBySlotID(ctx context.Context, slotID string) (*models.Reservation, error)
	GetBySlotIDTx(ctx context.Context, tx *sql.Tx, slotID string) (*models.Reservation, error)
	UpdateStatusTx(ctx context.Context, tx *sql.Tx, reservationID string, status models.ReservationState) error
	BulkExpireHoldsTx(ctx context.Context, tx *sql.Tx, olderThanSecs int) (int64, error)
}

type resRepo struct {
	db *sql.DB
}

func NewReservationRepository(db *sql.DB) ReservationRepository {
	return &resRepo{db: db}
}

func (r *resRepo) CreateTx(ctx context.Context, tx *sql.Tx, res models.Reservation) (*models.Reservation, error) {
	if res.Status == "" {
		res.Status = models.ReservationHold
	}
	err := tx.QueryRowContext(ctx, "INSERT INTO reservations(id, user_id, clinic_id, slot_id, payment_status_id, status, created_at, updated_at) VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id, status, created_at, updated_at", res.UserID, res.ClinicID, res.SlotID, res.PaymentStatusID, res.Status).Scan(&res.ID, &res.Status, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *resRepo) GetBySlotID(ctx context.Context, slotID string) (*models.Reservation, error) {
	var res models.Reservation
	err := r.db.QueryRowContext(ctx, "SELECT id, user_id, clinic_id, slot_id, payment_status_id, status, created_at, updated_at FROM reservations WHERE slot_id = $1 AND status IN ('hold', 'booked')", slotID).Scan(&res.ID, &res.UserID, &res.ClinicID, &res.SlotID, &res.PaymentStatusID, &res.Status, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &res, nil
}

func (r *resRepo) GetBySlotIDTx(ctx context.Context, tx *sql.Tx, slotID string) (*models.Reservation, error) {
	var res models.Reservation
	err := tx.QueryRowContext(ctx, "SELECT id, user_id, clinic_id, slot_id, payment_status_id, status, created_at, updated_at FROM reservations WHERE slot_id = $1 AND status IN ('hold', 'booked')", slotID).Scan(&res.ID, &res.UserID, &res.ClinicID, &res.SlotID, &res.PaymentStatusID, &res.Status, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &res, nil
}

func (r *resRepo) UpdateStatusTx(ctx context.Context, tx *sql.Tx, reservationID string, status models.ReservationState) error {
	_, err := tx.ExecContext(ctx, "UPDATE reservations SET status = $1, updated_at = NOW() WHERE id = $2", status, reservationID)
	return err
}

// BulkExpireHoldsTx runs a single bulk SQL statement using CTEs to release expired holds,
// update reservations to cancelled for audit purposes, and mark payments failed in 1 DB round trip.
func (r *resRepo) BulkExpireHoldsTx(ctx context.Context, tx *sql.Tx, olderThanSecs int) (int64, error) {
	query := `
		WITH expired_candidates AS (
			SELECT r.id AS res_id, r.slot_id, r.payment_status_id
			FROM reservations r
			JOIN slots s ON s.id = r.slot_id
			WHERE r.status = 'hold'
			  AND s.state = 'hold'
			  AND r.created_at < NOW() - ($1 || ' seconds')::interval
			FOR UPDATE OF s SKIP LOCKED
		),
		updated_slots AS (
			UPDATE slots s
			SET state = 'available', updated_at = NOW()
			FROM expired_candidates ec
			WHERE s.id = ec.slot_id
			RETURNING s.id
		),
		updated_payments AS (
			UPDATE payment_status p
			SET status = 'failed', updated_at = NOW()
			FROM expired_candidates ec
			WHERE p.id = ec.payment_status_id
			RETURNING p.id
		),
		updated_reservations AS (
			UPDATE reservations r
			SET status = 'cancelled', updated_at = NOW()
			FROM expired_candidates ec
			WHERE r.id = ec.res_id
			RETURNING r.id
		)
		SELECT count(*) FROM updated_reservations;
	`
	var count int64
	err := tx.QueryRowContext(ctx, query, olderThanSecs).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
