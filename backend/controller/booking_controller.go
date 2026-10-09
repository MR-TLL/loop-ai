package controller

import (
	"context"
	"database/sql"
	"time"

	"backend/models"
	"backend/repository"
	"backend/db"
)

type HoldResult struct {
	ReservationID string    `json:"reservation_id"`
	PaymentID     string    `json:"payment_id"`
	Amount        float64   `json:"amount"`
	ExpiresAt     time.Time `json:"expires_at"`
}

type BookingController struct {
	db          *sql.DB
	slotRepo    repository.SlotRepository
	resRepo     repository.ReservationRepository
	paymentRepo repository.PaymentRepository
}

func NewBookingController(d *sql.DB, slotRepo repository.SlotRepository, resRepo repository.ReservationRepository, paymentRepo repository.PaymentRepository) *BookingController {
	return &BookingController{
		db:          d,
		slotRepo:    slotRepo,
		resRepo:     resRepo,
		paymentRepo: paymentRepo,
	}
}

func (c *BookingController) HoldSlot(ctx context.Context, userID, slotID string) (*HoldResult, error) {
	tx, err := db.BeginTx(ctx, c.db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	slot, err := c.slotRepo.LockForUpdate(ctx, tx, slotID, models.SlotAvailable)
	if err != nil {
		if err == repository.ErrSlotNotAvailable {
			return nil, ErrSlotNotAvailable
		}
		return nil, err
	}

	// Amount comes dynamically from the slot entity
	pay, err := c.paymentRepo.CreatePaymentStatusTx(ctx, tx, models.Payment{
		Status: models.PaymentPending,
		Amount: slot.Amount,
	})
	if err != nil {
		return nil, err
	}

	res, err := c.resRepo.CreateTx(ctx, tx, models.Reservation{
		UserID:          userID,
		ClinicID:        slot.ClinicID,
		SlotID:          slotID,
		PaymentStatusID: pay.ID,
		Status:          models.ReservationHold,
	})
	if err != nil {
		return nil, err
	}

	if err := c.slotRepo.UpdateStateTx(ctx, tx, slotID, models.SlotHold); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &HoldResult{
		ReservationID: res.ID,
		PaymentID:     pay.ID,
		Amount:        slot.Amount,
		ExpiresAt:     res.CreatedAt.Add(5 * time.Minute),
	}, nil
}

// BookSlot finalizes a held slot: transitions slot from hold -> booked and updates payment status.
// Payment processing is orchestrated by the client/BFF with the payment provider prior to calling this.
func (c *BookingController) BookSlot(ctx context.Context, userID, slotID, paymentID string) (*models.Reservation, error) {
	res, err := c.resRepo.GetBySlotID(ctx, slotID)
	if err != nil {
		if err == repository.ErrNotFound {
			return nil, ErrReservationNotFound
		}
		return nil, err
	}

	if res.UserID != userID {
		return nil, ErrNotYourReservation
	}

	tx, err := db.BeginTx(ctx, c.db)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = c.slotRepo.LockForUpdate(ctx, tx, slotID, models.SlotHold)
	if err != nil {
		return nil, ErrInvalidSlotState
	}

	if err := c.paymentRepo.UpdateStatusTx(ctx, tx, paymentID, models.PaymentSuccess); err != nil {
		return nil, err
	}

	if err := c.slotRepo.UpdateStateTx(ctx, tx, slotID, models.SlotBooked); err != nil {
		return nil, err
	}

	if err := c.resRepo.UpdateStatusTx(ctx, tx, res.ID, models.ReservationBooked); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	res.Status = models.ReservationBooked
	return res, nil
}
