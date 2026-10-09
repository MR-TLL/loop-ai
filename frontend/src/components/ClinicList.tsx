import { useEffect, useState } from 'react';
import { api } from '../api/client';
import type { Clinic } from '../types';

interface Props {
  onSelect: (clinic: Clinic) => void;
}

export function ClinicList({ onSelect }: Props) {
  const [clinics, setClinics] = useState<Clinic[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    api.listClinics()
      .then(setClinics)
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <p>Loading clinics...</p>;
  if (error) return <p className="error">Error: {error}</p>;
  if (clinics.length === 0) return <p>No clinics found.</p>;

  return (
    <div className="clinic-list">
      {clinics.map(clinic => (
        <div key={clinic.id} className="clinic-card" onClick={() => onSelect(clinic)}>
          <h3>{clinic.name}</h3>
          <p>{clinic.address}</p>
          <button type="button">View Slots →</button>
        </div>
      ))}
    </div>
  );
}
