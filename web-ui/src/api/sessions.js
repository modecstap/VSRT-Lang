import { apiClient } from './client';
import {
  clearSessionStorage,
  readSessionStorage,
  writeSessionStorage,
} from './storage';

export const SESSIONS_SWR_KEY = ['sessions'];

export function sessionRecordsKey(sessionId) {
  return ['session-records', sessionId || 'active'];
}

export function getActiveSessionId() {
  return readSessionStorage()?.activeSessionId || '';
}

export function setActiveSessionId(sessionId) {
  if (!sessionId) {
    clearSessionStorage();
    return;
  }

  writeSessionStorage({ activeSessionId: sessionId });
}

export function clearActiveSessionId() {
  clearSessionStorage();
}

export async function fetchUserSessions() {
  const response = await apiClient.get('/users/sessions');
  return response.data || [];
}

export async function createSession(sessionName = 'Session') {
  const response = await apiClient.post('/sessions', { name: sessionName });
  const sessionId = response.data?.id;

  if (!sessionId) {
    throw new Error('Session was not created');
  }

  setActiveSessionId(sessionId);

  return String(sessionId);
}

export async function deleteSession(sessionId) {
  await apiClient.delete(`/sessions/${sessionId}`);

  if (getActiveSessionId() === String(sessionId)) {
    clearActiveSessionId();
  }
}

export async function getOrCreateSessionId() {
  const storedSessionId = getActiveSessionId();

  if (storedSessionId) {
    return storedSessionId;
  }

  return createSession();
}

export async function fetchSessionRecords(sessionId) {
  const response = await apiClient.get(`/sessions/${sessionId}/records`);
  return response.data?.records || [];
}

export async function saveSessionRecord(sessionId, payload) {
  const response = await apiClient.post(`/sessions/${sessionId}/records`, {
    phrase: payload.phrase,
    context: payload.context,
  });

  return response.data;
}
