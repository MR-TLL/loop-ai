package models

import "time"

type SlotState string

const (
	SlotInvalid   SlotState = "invalid" // zero-value sentinel - MUST BE FIRST
	SlotAvailable SlotState = "available"
	SlotHold      SlotState = "hold"
	SlotBooked    SlotState = "booked"
	SlotCancelled SlotState = "cancelled"
)

type Slot struct {
	ID        string    `json:"id"`
	ClinicID  string    `json:"clinic_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	State     SlotState `json:"state"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
