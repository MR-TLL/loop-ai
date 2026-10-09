package repository

import "errors"

var (
	ErrNotFound         = errors.New("record not found")
	ErrSlotNotAvailable = errors.New("slot not available or already locked")
)
