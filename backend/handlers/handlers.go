package handlers

type Handlers struct {
	Clinic  *ClinicHandler
	Slot    *SlotHandler
	Booking *BookingHandler
}

func NewHandlers(clinic *ClinicHandler, slot *SlotHandler, booking *BookingHandler) *Handlers {
	return &Handlers{
		Clinic:  clinic,
		Slot:    slot,
		Booking: booking,
	}
}
