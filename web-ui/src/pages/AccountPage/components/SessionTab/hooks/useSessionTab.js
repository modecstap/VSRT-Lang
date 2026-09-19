import { useState } from 'react';
import useSWR from 'swr';
import { getActiveSessionId } from '../../../../../api/sessions';
import { loadSavedWords, saveWordEntry, sessionRecordsKey } from '../api/sessionTabApi';
import {
  buildWordDetails,
  buildWordMap,
  mapRecordToWord,
  mergeSavedWord,
  normalizeWord,
} from '../model/sessionTabModel';

const createInitialForm = () => ({
  word: '',
  context: '',
});

function useSessionTab() {
  const sessionId = getActiveSessionId();
  const {
    data: words = [],
    error,
    isLoading,
    mutate,
  } = useSWR(sessionRecordsKey(sessionId), async () => {
    const records = await loadSavedWords();
    return records.map((record, index) => mapRecordToWord(record, index));
  });
  const [form, setForm] = useState(createInitialForm);
  const [selectedWord, setSelectedWord] = useState(null);
  const [isWriting, setIsWriting] = useState(false);
  const [writeError, setWriteError] = useState('');

  const wordsByNormalized = buildWordMap(words);
  const activeWord = selectedWord
    || (form.word ? wordsByNormalized.get(normalizeWord(form.word)) : words[0])
    || null;

  const handleWordChange = (event) => {
    const nextWord = event.target.value;
    setForm((current) => ({ ...current, word: nextWord }));
    setSelectedWord(wordsByNormalized.get(normalizeWord(nextWord)) || null);
  };

  const handleContextChange = (event) => {
    setForm((current) => ({ ...current, context: event.target.value }));
  };

  const handleWrite = async (event) => {
    event.preventDefault();

    const word = form.word.trim();
    const context = form.context.trim();

    if (!word) {
      return;
    }

    setIsWriting(true);
    setWriteError('');

    try {
      const savedRecord = await saveWordEntry({
        word,
        context,
      });
      const savedWord = mapRecordToWord(savedRecord, words.length);

      const nextWords = await mutate(
        (currentWords = []) => mergeSavedWord(currentWords, savedWord),
        { revalidate: false }
      );

      mutate();

      const nextEntry = buildWordMap(nextWords || []).get(normalizeWord(word)) || savedWord;
      setSelectedWord(nextEntry);
      setForm(createInitialForm());
    } catch (submitError) {
      setWriteError('Unable to save word');
    } finally {
      setIsWriting(false);
    }
  };

  const handleSelectWord = (word) => {
    setSelectedWord(word);
    setForm({ word: word.word, context: '' });
  };

  return {
    details: buildWordDetails(activeWord),
    error: error ? 'Unable to load words' : '',
    form,
    handleContextChange,
    handleSelectWord,
    handleWordChange,
    handleWrite,
    isWriting,
    loading: isLoading,
    selectedWord: activeWord,
    words,
    writeError,
  };
}

export default useSessionTab;
