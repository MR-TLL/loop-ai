import type { Slot } from '../types';
import { SlotCard } from './SlotCard';

interface Props {
  slots: Slot[];
  onHold: (slot: Slot) => void;
  currentUserId: string;
}

export function SlotGrid({ slots, onHold, currentUserId }: Props) {
  if (slots.length === 0) {
    return <p>No slots found.</p>;
  }

  // Group slots by date
  const grouped: { [key: string]: Slot[] } = {};
  slots.forEach((s) => {
    const d = s.start_time.slice(0, 10);
    if (!grouped[d]) grouped[d] = [];
    grouped[d].push(s);
  });

  const dates = Object.keys(grouped).sort();

  return (
    <div className="slot-grid">
      {dates.map((dateKey) => {
        const dateSlots = grouped[dateKey].sort(
          (a, b) => new Date(a.start_time).getTime() - new Date(b.start_time).getTime()
        );
        const headingDate = new Date(dateSlots[0].start_time).toLocaleDateString(undefined, {
          weekday: 'long',
          month: 'short',
          day: 'numeric',
          year: 'numeric',
        });

        return (
          <div key={dateKey} className="slot-date-group" style={{ marginBottom: '24px' }}>
            <h4 className="slot-date-heading">{headingDate}</h4>
            <div className="slot-cards" style={{ display: 'flex', flexWrap: 'wrap', gap: '12px' }}>
              {dateSlots.map((slot) => (
                <SlotCard
                  key={slot.id}
                  slot={slot}
                  onHold={onHold}
                  currentUserId={currentUserId}
                />
              ))}
            </div>
          </div>
        );
      })}
    </div>
  );
}
