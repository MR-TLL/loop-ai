import type { Slot } from '../types';
import { StatusBadge } from './StatusBadge';

interface Props {
  slot: Slot;
  onHold: (slot: Slot) => void;
  currentUserId: string;
}

export function SlotCard({ slot, onHold }: Props) {
  const start = new Date(slot.start_time);
  const end = new Date(slot.end_time);

  const timeString = `${start.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })} - ${end.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}`;
  const dateString = start.toLocaleDateString(undefined, { weekday: 'short', month: 'short', day: 'numeric' });

  return (
    <div className="slot-card">
      <div className="slot-header">
        <strong>{timeString}</strong>
      </div>
      <div className="slot-date">{dateString}</div>
      <div className="slot-price" style={{ fontWeight: 600, color: '#2c3e50', margin: '4px 0' }}>
        ₹{slot.amount}
      </div>
      <div className="slot-status">
        <StatusBadge state={slot.state} />
      </div>
      <div className="slot-actions" style={{ marginTop: '8px' }}>
        {slot.state === 'available' && (
          <button className="btn btn-green" onClick={() => onHold(slot)}>Hold Slot</button>
        )}
        {slot.state === 'booked' && (
          <span style={{ fontSize: '13px', color: '#4a5568' }}>Confirmed</span>
        )}
        {slot.state === 'hold' && (
          <button className="btn btn-disabled" disabled>On Hold</button>
        )}
      </div>
    </div>
  );
}
