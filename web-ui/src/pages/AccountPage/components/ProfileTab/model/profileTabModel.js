export function mapSessionToViewModel(session = {}) {
    return {
        id: session.ID,
        date: session.CreatedAt || '—',
        name: session.User || 'Untitled session',
        saved: session.Records ? session.Records.lenght(): 0,
    };
}

export function buildSessionsViewModel(sessions = []) {
  return sessions.map(mapSessionToViewModel);
}
