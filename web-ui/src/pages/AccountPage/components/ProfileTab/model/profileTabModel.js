export function mapSessionToViewModel(session = {}, index = 0) {
  return {
    id: session.id || session.ID || `session-${index}`,
    date: session.CreatedAt || session.createdAt || '—',
    name: session.Name || 'Untitled session',
    saved: Array.isArray(session.Records) ? session.Records.length : 0,
  };
}

export function buildSessionsViewModel(sessions = []) {
  return sessions.map((session, index) => mapSessionToViewModel(session, index));
}
