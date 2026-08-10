import axios from 'axios';
import { getAuthHeaders } from '../../../../../api/auth';
import { normalizeWord } from '../model/sessionTabModel';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
const SESSION_STORAGE_KEY = 'session_tab_session_id';

export async function createSession(sessionName = 'Session') {
  const response = await axios.post(
    `${API_BASE_URL}/sessions`,
    { name: sessionName },
    getAuthHeaders()
  );

  const sessionId = response.data?.id;

  if (!sessionId) {
    throw new Error('Session was not created');
  }

  localStorage.setItem(SESSION_STORAGE_KEY, String(sessionId));

  return sessionId;
}

const getOrCreateSessionId = async () => {
  const storedSessionId = localStorage.getItem(SESSION_STORAGE_KEY);

  if (storedSessionId) {
    return storedSessionId;
  }

  return createSession();
};

const mapRecordToWord = (record, index = 0) => {
  const phrase = record?.phrase || record?.Phrase || '';
  const translations = record?.translations || record?.Translations || [];
  const synonyms = record?.synonyms || record?.Synonyms || [];
  const antonyms = record?.antonyms || record?.Antonyms || [];
  const baseForm = record?.baseForm || record?.BaseForm || phrase;
  const contexts = record?.contexts || record?.Contexts || [];

  return {
    id: record?.id || `${normalizeWord(phrase) || 'word'}-${index}`,
    word: phrase,
    translation: translations[0] || '—',
    synonyms,
    antonyms,
    meaning: baseForm || 'No meaning provided yet.',
    contexts: contexts.map((item) => item?.phrase || item?.Phrase || item?.translation || item?.Translation || ''),
  };
};

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

  return records.map((record, index) => mapRecordToWord(record, index));
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

  const records = recordsResponse.data?.records || [];

  return records.map((record, index) => mapRecordToWord(record, index));
}
