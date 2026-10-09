package models

import "time"

type PaymentState string

const (
	PaymentInvalid PaymentState = "invalid" // zero-value sentinel - MUST BE FIRST
	PaymentPending PaymentState = "pending"
	PaymentSuccess PaymentState = "success"
	PaymentFailed  PaymentState = "failed"
)

type Payment struct {
	ID        string       `json:"id"`
	Status    PaymentState `json:"status"`
	Amount    float64      `json:"amount"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}
