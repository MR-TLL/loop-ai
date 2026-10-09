package repository

import (
	"context"
	"database/sql"
	"time"
	"backend/models"
)

type SlotRepository interface {
	GetByClinicID(ctx context.Context, clinicID string) ([]models.Slot, error)
	GetByID(ctx context.Context, id string) (*models.Slot, error)
	LockForUpdate(ctx context.Context, tx *sql.Tx, slotID string, expectedState models.SlotState) (*models.Slot, error)
	UpdateStateTx(ctx context.Context, tx *sql.Tx, slotID string, newState models.SlotState) error
	GetExpiredHoldIDs(ctx context.Context, olderThan time.Duration) ([]string, error)
}

type slotRepo struct {
	db *sql.DB
}

func NewSlotRepository(db *sql.DB) SlotRepository {
	return &slotRepo{db: db}
}

func (r *slotRepo) GetByClinicID(ctx context.Context, clinicID string) ([]models.Slot, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, clinic_id, start_time, end_time, state, amount, created_at, updated_at FROM slots WHERE clinic_id = $1", clinicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []models.Slot
	for rows.Next() {
		var s models.Slot
		if err := rows.Scan(&s.ID, &s.ClinicID, &s.StartTime, &s.EndTime, &s.State, &s.Amount, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return slots, nil
}

func (r *slotRepo) GetByID(ctx context.Context, id string) (*models.Slot, error) {
	var s models.Slot
	err := r.db.QueryRowContext(ctx, "SELECT id, clinic_id, start_time, end_time, state, amount, created_at, updated_at FROM slots WHERE id = $1", id).Scan(&s.ID, &s.ClinicID, &s.StartTime, &s.EndTime, &s.State, &s.Amount, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *slotRepo) LockForUpdate(ctx context.Context, tx *sql.Tx, slotID string, expectedState models.SlotState) (*models.Slot, error) {
	var s models.Slot
	err := tx.QueryRowContext(ctx, "SELECT id, clinic_id, start_time, end_time, state, amount, created_at, updated_at FROM slots WHERE id = $1 AND state = $2 FOR UPDATE SKIP LOCKED", slotID, expectedState).Scan(&s.ID, &s.ClinicID, &s.StartTime, &s.EndTime, &s.State, &s.Amount, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrSlotNotAvailable
		}
		return nil, err
	}
	return &s, nil
}

func (r *slotRepo) UpdateStateTx(ctx context.Context, tx *sql.Tx, slotID string, newState models.SlotState) error {
	_, err := tx.ExecContext(ctx, "UPDATE slots SET state = $1, updated_at = NOW() WHERE id = $2", newState, slotID)
	return err
}

func (r *slotRepo) GetExpiredHoldIDs(ctx context.Context, olderThan time.Duration) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT s.id FROM slots s JOIN reservations r ON r.slot_id = s.id WHERE s.state='hold' AND r.status='hold' AND r.created_at < NOW() - ($1 || ' seconds')::interval", int(olderThan.Seconds()))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}
