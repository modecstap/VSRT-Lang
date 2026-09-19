export { SESSIONS_SWR_KEY } from '../../../../../api/sessions';

export function mapSessionToViewModel(session = {}, index = 0) {
  const records = session.Records || session.records || {};
  const saved = Array.isArray(records)
    ? records.length
    : Object.keys(records).length;

  return {
    id: session.id || session.ID || `session-${index}`,
    date: session.CreatedAt || session.createdAt || '—',
    name: session.Name || session.name || 'Untitled session',
    saved,
  };
}

export function buildSessionsViewModel(sessions = []) {
  return sessions.map((session, index) => mapSessionToViewModel(session, index));
}
