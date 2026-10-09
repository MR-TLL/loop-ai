import type { Slot, HoldResponse } from '../types';

interface Props {
  slot: Slot;
  holdData: HoldResponse;
  onConfirm: () => void;
  onClose: () => void;
  loading: boolean;
  error: string | null;
}

export function BookingModal({ slot, holdData, onConfirm, onClose, loading, error }: Props) {
  const start = new Date(slot.start_time);
  const timeString = `${start.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
  const dateString = start.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' });
  const expiresAt = new Date(holdData.expires_at).toLocaleTimeString();

  return (
    <div className="modal-overlay">
      <div className="modal-box">
        <h2>Confirm Booking</h2>
        <p><strong>Date & Time:</strong> {dateString} at {timeString}</p>
        <p><strong>Amount:</strong> ₹{slot.amount}</p>
        <p><strong>Hold expires at:</strong> {expiresAt}</p>
        <p style={{ fontSize: '12px', color: '#666' }}>Payment ID: {holdData.payment_id}</p>
        
        {error && <p className="error">{error}</p>}
        
        <div className="modal-actions" style={{ marginTop: '20px', display: 'flex', gap: '10px', justifyContent: 'flex-end' }}>
          <button className="btn btn-secondary" onClick={onClose} disabled={loading}>Cancel</button>
          <button className="btn btn-primary" onClick={onConfirm} disabled={loading}>Confirm & Pay (Mocked)</button>
        </div>
      </div>
    </div>
  );
}
