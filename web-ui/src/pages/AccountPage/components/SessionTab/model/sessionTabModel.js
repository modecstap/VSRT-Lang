export { sessionRecordsKey } from '../../../../../api/sessions';

export function normalizeWord(word = '') {
  return word.trim().toLowerCase();
}

export function mapRecordToWord(record, index = 0) {
  const phrase = record?.phrase || record?.Phrase || '';
  const translations = record?.translations || record?.Translations || [];
  const synonyms = record?.synonyms || record?.Synonyms || [];
  const antonyms = record?.antonyms || record?.Antonyms || [];
  const baseForm = record?.baseForm || record?.BaseForm || record?.base_form || phrase;
  const contexts = record?.contexts || record?.Contexts || [];

  return {
    id: record?.id || `${normalizeWord(phrase) || 'word'}-${index}`,
    word: phrase,
    translation: translations,
    synonyms,
    antonyms,
    meaning: baseForm || 'No meaning provided yet.',
    contexts: contexts.map((item) => item?.phrase || item?.Phrase || item?.translation || item?.Translation || ''),
  };
}

export function buildWordMap(words = []) {
  const wordsByNormalized = new Map();

  for (const item of words) {
    wordsByNormalized.set(normalizeWord(item.word), item);
  }

  return wordsByNormalized;
}

export function findWordEntry(word, words = []) {
  return buildWordMap(words).get(normalizeWord(word)) || null;
}

export function mergeSavedWord(words = [], savedWord) {
  const normalized = normalizeWord(savedWord.word);
  const nextWords = words.slice();
  const existingIndex = nextWords.findIndex((item) => normalizeWord(item.word) === normalized);

  if (existingIndex === -1) {
    return [...nextWords, savedWord];
  }

  const current = nextWords[existingIndex];
  const contexts = Array.from(new Set([
    ...(current.contexts || []),
    ...(savedWord.contexts || []),
  ].filter(Boolean)));

  nextWords[existingIndex] = {
    ...current,
    ...savedWord,
    id: current.id,
    translation: savedWord.translation?.length ? savedWord.translation : current.translation,
    synonyms: savedWord.synonyms?.length ? savedWord.synonyms : current.synonyms,
    antonyms: savedWord.antonyms?.length ? savedWord.antonyms : current.antonyms,
    contexts,
  };

  return nextWords;
}

export function buildWordDetails(entry) {
  if (!entry) {
    return null;
  }

  const translation = Array.isArray(entry.translation)
    ? entry.translation.join('\t|\t')
    : entry.translation;

  return {
    id: entry.id,
    word: entry.word,
    translation: translation || '—',
    synonyms: entry.synonyms || [],
    antonyms: entry.antonyms || [],
    meaning: entry.meaning || 'No meaning provided yet.',
    contexts: entry.contexts || [],
  };
}
