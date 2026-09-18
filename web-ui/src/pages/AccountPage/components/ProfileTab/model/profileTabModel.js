export function mapSessionToViewModel(session = {}, index = 0) {
  return {
    id: session.id || session.ID || `session-${index}`,
    date: session.CreatedAt || session.createdAt || '—',
    name: session.Name || 'Untitled session',
    saved: Object.keys(session.Records).length,
  };
}

export function buildSessionsViewModel(sessions = []) {
  return sessions.map((session, index) => mapSessionToViewModel(session, index));
}
