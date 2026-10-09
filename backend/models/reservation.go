package models

import "time"

type ReservationState string

const (
	ReservationInvalid   ReservationState = "invalid" // zero-value sentinel - MUST BE FIRST
	ReservationHold      ReservationState = "hold"
	ReservationBooked    ReservationState = "booked"
	ReservationCancelled ReservationState = "cancelled"
)

type Reservation struct {
	ID              string           `json:"id"`
	UserID          string           `json:"user_id"`
	ClinicID        string           `json:"clinic_id"`
	SlotID          string           `json:"slot_id"`
	PaymentStatusID string           `json:"payment_status_id"`
	Status          ReservationState `json:"status"`
	CreatedAt       time.Time        `json:"created_at"` // expiry anchor
	UpdatedAt       time.Time        `json:"updated_at"`
}
