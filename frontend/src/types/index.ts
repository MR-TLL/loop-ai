export type SlotState = 'invalid' | 'available' | 'hold' | 'booked' | 'cancelled';
export type PaymentState = 'invalid' | 'pending' | 'success' | 'failed';
export type ReservationState = 'invalid' | 'hold' | 'booked' | 'cancelled';

export interface Clinic {
  id: string;
  name: string;
  address: string;
  created_at: string;
  updated_at: string;
}

export interface Slot {
  id: string;
  clinic_id: string;
  start_time: string;
  end_time: string;
  state: SlotState;
  amount: number;
  created_at: string;
  updated_at: string;
}

export interface HoldResponse {
  reservation_id: string;
  payment_id: string;
  amount: number;
  expires_at: string;
}

export interface BookResponse {
  id: string;
  user_id: string;
  clinic_id: string;
  slot_id: string;
  payment_status_id: string;
  status: ReservationState;
  created_at: string;
  updated_at: string;
}
