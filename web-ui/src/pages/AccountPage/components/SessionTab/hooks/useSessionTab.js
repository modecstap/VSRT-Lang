import { useEffect, useState } from 'react';
import { loadSavedWords, saveWordEntry } from '../api/sessionTabApi';
import { buildWordDetails, normalizeWord } from '../model/sessionTabModel';

const initialForm = {
  word: '',
  context: '',
};

function useSessionTab() {
  const [words, setWords] = useState([]);
  const [selectedWord, setSelectedWord] = useState(null);
  const [form, setForm] = useState(initialForm);

  useEffect(() => {
    let ignore = false;

    const loadData = async () => {
      try {
        const savedWords = await loadSavedWords();

        if (!ignore) {
          const nextWords = savedWords || [];
          setWords(nextWords);
          setSelectedWord(nextWords[0] || null);
        }
      } catch (error) {
        if (!ignore) {
          setWords([]);
          setSelectedWord(null);
        }
      }
    };

    loadData();

    return () => {
      ignore = true;
    };
  }, []);

  const handleWordChange = (event) => {
    const nextWord = event.target.value;
    setForm((current) => ({ ...current, word: nextWord }));

    const entry = words.find((item) => normalizeWord(item.word) === normalizeWord(nextWord));
    setSelectedWord(entry || null);
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
      context,
    });

    setWords(savedWords);
    const nextEntry = savedWords.find(
      (entry) => normalizeWord(entry.word) === normalizeWord(word)
    );

    setSelectedWord(nextEntry || null);
    setForm(initialForm);
  };

  const handleSelectWord = (word) => {
    setSelectedWord(word);
    setForm({ word: word.word, context: '' });
  };

  return {
    details: buildWordDetails(selectedWord),
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
