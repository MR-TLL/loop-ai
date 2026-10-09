package controller

import (
	"context"
	"database/sql"
	"time"
	"backend/repository"
	"backend/db"
)

type ExpiryController struct {
	db          *sql.DB
	resRepo     repository.ReservationRepository
	holdTimeout time.Duration
}

func NewExpiryController(d *sql.DB, slotRepo repository.SlotRepository, resRepo repository.ReservationRepository, paymentRepo repository.PaymentRepository, holdTimeout time.Duration) *ExpiryController {
	return &ExpiryController{
		db:          d,
		resRepo:     resRepo,
		holdTimeout: holdTimeout,
	}
}

// ExpireHolds issues a single bulk SQL statement using a CTE to release expired holds,
// update reservations to cancelled state (for audit retention), and mark payment_status failed.
// This executes in a single DB round trip and transaction, minimizing DB TPS and latency.
func (c *ExpiryController) ExpireHolds(ctx context.Context) (int64, error) {
	tx, err := db.BeginTx(ctx, c.db)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	count, err := c.resRepo.BulkExpireHoldsTx(ctx, tx, int(c.holdTimeout.Seconds()))
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return count, nil
}
