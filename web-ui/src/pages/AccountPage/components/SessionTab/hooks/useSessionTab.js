import { useEffect, useState } from 'react';
import { fetchWordEntry, loadSavedWords, saveWordEntry } from '../api/sessionTabApi';
import { buildWordDetails, normalizeWord } from '../model/sessionTabModel';

const initialForm = {
  word: '',
  context: '',
};

function useSessionTab() {
  const [words, setWords] = useState([]);
  const [details, setDetails] = useState(null);
  const [selectedWord, setSelectedWord] = useState(null);
  const [form, setForm] = useState(initialForm);

  useEffect(() => {
    let ignore = false;

    loadSavedWords().then((savedWords) => {
      if (!ignore) {
        setWords(savedWords);
        setDetails(buildWordDetails(savedWords[0] || null));
        setSelectedWord(savedWords[0] || null);
      }
    });

    return () => {
      ignore = true;
    };
  }, []);

  const handleWordChange = async (event) => {
    const nextWord = event.target.value;
    setForm((current) => ({ ...current, word: nextWord }));

    if (!nextWord.trim()) {
      setDetails(null);
      return;
    }

    const entry = await fetchWordEntry(nextWord);
    setDetails(buildWordDetails(entry));
    setSelectedWord(entry);
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

    const savedWords = await saveWordEntry({
      word,
      translation: '—',
      synonyms: [],
      antonyms: [],
      meaning: 'Saved from the session form.',
      context,
    });

    setWords(savedWords);
    const nextEntry = savedWords.find(
      (entry) => normalizeWord(entry.word) === normalizeWord(word)
    );

    setSelectedWord(nextEntry || null);
    setDetails(buildWordDetails(nextEntry || null));
    setForm(initialForm);
  };

  const handleSelectWord = (word) => {
    setSelectedWord(word);
    setDetails(buildWordDetails(word));
    setForm({ word: word.word, context: '' });
  };

  return {
    details,
    form,
    handleContextChange,
    handleSelectWord,
    handleWordChange,
    handleWrite,
    selectedWord,
    words,
  };
}

export default useSessionTab;
