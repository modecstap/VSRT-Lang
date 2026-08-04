const storedWords = [
  {
    id: 1,
    word: 'apple',
    translation: 'яблоко',
    synonyms: ['fruit'],
    antonyms: ['orange'],
    meaning: 'A round fruit often eaten fresh.',
    contexts: ['An apple a day keeps the doctor away.'],
  },
  {
    id: 2,
    word: 'river',
    translation: 'река',
    synonyms: ['stream'],
    antonyms: ['desert'],
    meaning: 'A natural watercourse flowing toward a sea or lake.',
    contexts: ['The river was calm at sunrise.'],
  },
];

const normalizeWord = (word = '') => word.trim().toLowerCase();

export async function loadSavedWords() {
  return storedWords.map((item) => ({ ...item }));
}

export async function fetchWordEntry(word) {
  const normalized = normalizeWord(word);
  const entry = storedWords.find((item) => normalizeWord(item.word) === normalized);

  return entry ? { ...entry } : null;
}

export async function saveWordEntry(payload) {
  const normalized = normalizeWord(payload.word);
  const existingIndex = storedWords.findIndex(
    (item) => normalizeWord(item.word) === normalized
  );

  const entry = {
    id: existingIndex >= 0 ? storedWords[existingIndex].id : Date.now(),
    word: payload.word.trim(),
    translation: payload.translation || '—',
    synonyms: payload.synonyms || [],
    antonyms: payload.antonyms || [],
    meaning: payload.meaning || 'Saved from the session form.',
    contexts: payload.context ? [payload.context] : [],
  };

  if (existingIndex >= 0) {
    storedWords[existingIndex] = entry;
  } else {
    storedWords.unshift(entry);
  }

  return storedWords.map((item) => ({ ...item }));
}
