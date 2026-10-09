import type { SlotState } from '../types';

const STATE_COLORS: Record<SlotState, string> = {
  invalid: '#e53e3e',
  available: '#38a169',
  hold: '#d69e2e',
  booked: '#3182ce',
  cancelled: '#718096',
};

export function StatusBadge({ state }: { state: SlotState }) {
  return (
    <span style={{
      backgroundColor: STATE_COLORS[state],
      color: '#fff',
      padding: '2px 10px',
      borderRadius: '12px',
      fontSize: '12px',
      fontWeight: 600,
      textTransform: 'uppercase',
    }}>
      {state}
    </span>
  );
}
