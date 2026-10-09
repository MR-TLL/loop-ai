package worker

import (
	"context"
	"database/sql"
	"time"

	"backend/controller"
)

const advisoryLockKey = int64(20261009)

type ExpiryWorker struct {
	db       *sql.DB
	ctrl     *controller.ExpiryController
	interval time.Duration
}

func NewExpiryWorker(db *sql.DB, ctrl *controller.ExpiryController, interval time.Duration) *ExpiryWorker {
	return &ExpiryWorker{
		db:       db,
		ctrl:     ctrl,
		interval: interval,
	}
}

func (w *ExpiryWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.runOnce(ctx)
		}
	}
}

func (w *ExpiryWorker) runOnce(ctx context.Context) {
	var acquired bool
	err := w.db.QueryRowContext(ctx, `SELECT pg_try_advisory_lock($1)`, advisoryLockKey).Scan(&acquired)
	if err != nil || !acquired {
		return
	}
	defer w.db.ExecContext(ctx, `SELECT pg_advisory_unlock($1)`, advisoryLockKey)
	w.ctrl.ExpireHolds(ctx)
}
