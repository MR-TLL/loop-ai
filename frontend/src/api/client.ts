import type { Clinic, Slot, HoldResponse, BookResponse } from '../types';

const BASE = (import.meta.env.VITE_API_BASE_URL as string) || 'http://localhost:8080/api';

async function apiFetch<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: `HTTP ${res.status}` }));
    throw new Error((body as { error?: string }).error ?? `HTTP ${res.status}`);
  }
  return res.json() as Promise<T>;
}

export const api = {
  listClinics: () =>
    apiFetch<Clinic[]>('/v1/clinics'),

  getClinic: (clinicId: string) =>
    apiFetch<Clinic>(`/v1/clinics/${clinicId}`),

  getSlots: (clinicId: string) =>
    apiFetch<Slot[]>(`/v1/clinics/${clinicId}/slots`),

  holdSlot: (slotId: string, userId: string) =>
    apiFetch<HoldResponse>(`/v1/slots/${slotId}/hold`, {
      method: 'POST',
      body: JSON.stringify({ user_id: userId }),
    }),

  bookSlot: (slotId: string, userId: string, paymentId: string) =>
    apiFetch<BookResponse>(`/v1/slots/${slotId}/book`, {
      method: 'POST',
      body: JSON.stringify({ user_id: userId, payment_id: paymentId }),
    }),
};
