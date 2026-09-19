const AUTH_STORAGE_KEY = 'vsrt.auth.v1';
const SESSION_STORAGE_KEY = 'vsrt.session.v1';

const LEGACY_ACCESS_TOKEN_KEY = 'access_token';
const LEGACY_REFRESH_TOKEN_KEY = 'refresh_token';
const LEGACY_SESSION_ID_KEY = 'session_tab_session_id';

function readJson(key) {
  try {
    const raw = window.localStorage.getItem(key);
    if (!raw) {
      return null;
    }

    return JSON.parse(raw);
  } catch {
    return null;
  }
}

function writeJson(key, value) {
  window.localStorage.setItem(key, JSON.stringify(value));
}

function migrateLegacyAuth() {
  const existing = readJson(AUTH_STORAGE_KEY);
  if (existing?.v === 1) {
    return existing;
  }

  const accessToken = window.localStorage.getItem(LEGACY_ACCESS_TOKEN_KEY);
  const refreshToken = window.localStorage.getItem(LEGACY_REFRESH_TOKEN_KEY);

  if (!accessToken && !refreshToken) {
    return null;
  }

  const next = {
    v: 1,
    accessToken: accessToken || '',
    refreshToken: refreshToken || '',
  };

  writeJson(AUTH_STORAGE_KEY, next);
  window.localStorage.removeItem(LEGACY_ACCESS_TOKEN_KEY);
  window.localStorage.removeItem(LEGACY_REFRESH_TOKEN_KEY);

  return next;
}

function migrateLegacySession() {
  const existing = readJson(SESSION_STORAGE_KEY);
  if (existing?.v === 1) {
    return existing;
  }

  const activeSessionId = window.localStorage.getItem(LEGACY_SESSION_ID_KEY);
  if (!activeSessionId) {
    return null;
  }

  const next = {
    v: 1,
    activeSessionId: String(activeSessionId),
  };

  writeJson(SESSION_STORAGE_KEY, next);
  window.localStorage.removeItem(LEGACY_SESSION_ID_KEY);

  return next;
}

export function readAuthStorage() {
  return migrateLegacyAuth() || readJson(AUTH_STORAGE_KEY);
}

export function writeAuthStorage(tokens) {
  writeJson(AUTH_STORAGE_KEY, {
    v: 1,
    accessToken: tokens.accessToken || '',
    refreshToken: tokens.refreshToken || '',
  });
}

export function clearAuthStorage() {
  window.localStorage.removeItem(AUTH_STORAGE_KEY);
  window.localStorage.removeItem(LEGACY_ACCESS_TOKEN_KEY);
  window.localStorage.removeItem(LEGACY_REFRESH_TOKEN_KEY);
}

export function readSessionStorage() {
  return migrateLegacySession() || readJson(SESSION_STORAGE_KEY);
}

export function writeSessionStorage(session) {
  writeJson(SESSION_STORAGE_KEY, {
    v: 1,
    activeSessionId: session.activeSessionId ? String(session.activeSessionId) : '',
  });
}

export function clearSessionStorage() {
  window.localStorage.removeItem(SESSION_STORAGE_KEY);
  window.localStorage.removeItem(LEGACY_SESSION_ID_KEY);
}
