import {
  fetchSessionRecords,
  getOrCreateSessionId,
  saveSessionRecord,
  sessionRecordsKey,
} from '../../../../../api/sessions';

export { sessionRecordsKey };

export async function loadSavedWords() {
  const sessionId = await getOrCreateSessionId();

  if (!sessionId) {
    return [];
  }

  return fetchSessionRecords(sessionId);
}

export function fetchWordEntry(word, words = []) {
  return Array.isArray(words)
    ? words.find((item) => item.word.trim().toLowerCase() === word.trim().toLowerCase()) || null
    : null;
}

export async function saveWordEntry(payload) {
  const sessionId = await getOrCreateSessionId();

  if (!sessionId) {
    throw new Error('Session was not created');
  }

  return saveSessionRecord(sessionId, {
    phrase: payload.word.trim(),
    context: payload.context,
  });
}
