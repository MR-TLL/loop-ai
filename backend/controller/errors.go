package controller

import "errors"

var (
	ErrSlotNotAvailable    = errors.New("slot is not available")
	ErrInvalidSlotState    = errors.New("slot is not in the required state")
	ErrNotYourReservation  = errors.New("reservation does not belong to this user")
	ErrReservationNotFound = errors.New("reservation not found")
)
