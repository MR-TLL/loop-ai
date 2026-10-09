import { useState, useEffect } from 'react';
import './App.css';
import { ClinicList } from './components/ClinicList';
import { SlotGrid } from './components/SlotGrid';
import { BookingModal } from './components/BookingModal';
import { api } from './api/client';
import type { Clinic, Slot, HoldResponse } from './types';

const USER_ID = 'a0000000-0000-0000-0000-000000000001';

function App() {
  const [view, setView] = useState<'clinics' | 'detail'>('clinics');
  const [selectedClinic, setSelectedClinic] = useState<Clinic | null>(null);
  const [slots, setSlots] = useState<Slot[]>([]);
  const [loadingSlots, setLoadingSlots] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [holdData, setHoldData] = useState<HoldResponse | null>(null);
  const [holdSlot, setHoldSlot] = useState<Slot | null>(null);
  const [modalLoading, setModalLoading] = useState(false);
  const [modalError, setModalError] = useState<string | null>(null);

  const loadSlots = async (clinicId: string) => {
    setLoadingSlots(true);
    try {
      const data = await api.getSlots(clinicId);
      setSlots(data);
      setError(null);
    } catch (e: any) {
      setError(e.message || 'Failed to load slots');
    } finally {
      setLoadingSlots(false);
    }
  };

  useEffect(() => {
    let interval: number;
    if (view === 'detail' && selectedClinic) {
      loadSlots(selectedClinic.id);
      interval = window.setInterval(() => {
        loadSlots(selectedClinic.id);
      }, 30000);
    }
    return () => clearInterval(interval);
  }, [view, selectedClinic]);

  const handleSelectClinic = (clinic: Clinic) => {
    setSelectedClinic(clinic);
    setView('detail');
  };

  const handleHold = async (slot: Slot) => {
    try {
      const resp = await api.holdSlot(slot.id, USER_ID);
      setHoldData(resp);
      setHoldSlot(slot);
      if (selectedClinic) loadSlots(selectedClinic.id);
    } catch (e: any) {
      alert(e.message || 'Failed to hold slot');
    }
  };

  const handleBook = async () => {
    if (!holdSlot || !holdData) return;
    setModalLoading(true);
    setModalError(null);
    try {
      await api.bookSlot(holdSlot.id, USER_ID, holdData.payment_id);
      setHoldData(null);
      setHoldSlot(null);
      if (selectedClinic) loadSlots(selectedClinic.id);
    } catch (e: any) {
      setModalError(e.message || 'Failed to book slot');
    } finally {
      setModalLoading(false);
    }
  };

  const handleModalClose = () => {
    setHoldData(null);
    setHoldSlot(null);
    setModalError(null);
  };

  if (view === 'clinics') {
    return (
      <div className="container">
        <h1>🏥 Clinic Slot Booking</h1>
        <ClinicList onSelect={handleSelectClinic} />
      </div>
    );
  }

  return (
    <div className="container">
      <div className="detail-header">
        <button className="back-btn" onClick={() => setView('clinics')}>← Back to Clinics</button>
        <div className="user-tag">Booking as user: {USER_ID}</div>
      </div>
      
      {selectedClinic && (
        <div className="clinic-info">
          <h2>{selectedClinic.name}</h2>
          <p>{selectedClinic.address}</p>
        </div>
      )}

      {error && <div className="error banner">{error}</div>}

      <div className="slots-section">
        <h3>Available Slots {loadingSlots && <small>(Refreshing...)</small>}</h3>
        <SlotGrid 
          slots={slots} 
          onHold={handleHold} 
          currentUserId={USER_ID} 
        />
      </div>

      {holdData && holdSlot && (
        <BookingModal
          slot={holdSlot}
          holdData={holdData}
          onConfirm={handleBook}
          onClose={handleModalClose}
          loading={modalLoading}
          error={modalError}
        />
      )}
    </div>
  );
}

export default App;
