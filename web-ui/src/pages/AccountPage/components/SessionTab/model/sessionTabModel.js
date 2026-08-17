export function normalizeWord(word = '') {
  return word.trim().toLowerCase();
}

export function buildWordDetails(entry) {
  if (!entry) {
    return null;
  }

  return {
    id: entry.id,
    word: entry.word,
    translation: entry.translation.join("\t|\t") || '—',
    synonyms: entry.synonyms || [],
    antonyms: entry.antonyms || [],
    meaning: entry.meaning || 'No meaning provided yet.',
    contexts: entry.contexts || [],
  };
}
