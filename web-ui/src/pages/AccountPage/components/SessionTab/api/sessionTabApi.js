import axios from 'axios';
import { getAuthHeaders } from '../../../../../api/auth';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
const SESSION_STORAGE_KEY = 'session_tab_session_id';

const getOrCreateSessionId = async () => {
  const storedSessionId = localStorage.getItem(SESSION_STORAGE_KEY);

  if (storedSessionId) {
    return storedSessionId;
  }

  const sessionResponse = await axios.post(
    `${API_BASE_URL}/sessions`,
    { name: 'Session' },
    getAuthHeaders()
  );

  const sessionId = sessionResponse.data?.id;

  localStorage.setItem(SESSION_STORAGE_KEY, String(sessionId));

  return sessionId;
};

const mapRecordToWord = (record) => ({
  word: record.Phrase,
  translation: record.Translations?.[0] || '—',
  synonyms: record.Synonyms || [],
  antonyms: record.Antonyms || [],
  meaning: record.BaseForm || record.Phrase,
  contexts: (record.Contexts || []).map((item) => item.Phrase || item.Translation || ''),
});

export async function loadSavedWords() {
  const sessionId = await getOrCreateSessionId();

  if (!sessionId) {
    return [];
  }

  const recordsResponse = await axios.get(
    `${API_BASE_URL}/sessions/${sessionId}/records`,
    getAuthHeaders()
  );

  const records = recordsResponse.data?.records || [];

  return records.map(mapRecordToWord);
}

export async function fetchWordEntry(word, words = null) {
  if (Array.isArray(words)) {
    const normalized = word.trim().toLowerCase();
    return words.find((item) => item.word.trim().toLowerCase() === normalized) || null;
  }

  const allWords = await loadSavedWords();
  const normalized = word.trim().toLowerCase();

  return allWords.find((item) => item.word.trim().toLowerCase() === normalized) || null;
}

export async function saveWordEntry(payload) {
  const sessionId = await getOrCreateSessionId();

  if (!sessionId) {
    throw new Error('Session was not created');
  }

  await axios.post(
    `${API_BASE_URL}/sessions/${sessionId}/records`,
    {
      phrase: payload.word.trim(),
      context: payload.context,
    },
    getAuthHeaders()
  );

  const recordsResponse = await axios.get(
    `${API_BASE_URL}/sessions/${sessionId}/records`,
    getAuthHeaders()
  );

  console.log('Records response:', recordsResponse.data);
  const records = recordsResponse.data?.records || [];

  return records.map(mapRecordToWord);
}
