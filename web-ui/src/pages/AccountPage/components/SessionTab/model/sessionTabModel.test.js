import { mapRecordToWord, mergeSavedWord, removeSavedWord } from './sessionTabModel';

test('removeSavedWord drops the matching phrase and leaves the input unchanged', () => {
  const words = [{ word: 'Hello' }, { word: 'world' }];

  expect(removeSavedWord(words, 'hello')).toEqual([{ word: 'world' }]);
  expect(words).toEqual([{ word: 'Hello' }, { word: 'world' }]);
});

test('mapRecordToWord reads count and Count and falls back to 0', () => {
  expect(mapRecordToWord({ phrase: 'hello', count: 1 }).count).toBe(1);
  expect(mapRecordToWord({ phrase: 'hello', Count: 4 }).count).toBe(4);
  expect(mapRecordToWord({ phrase: 'hello' }).count).toBe(0);
});

test('mergeSavedWord replaces count and leaves the input unchanged', () => {
  const words = [{ word: 'hello', count: 1, contexts: [] }];

  expect(mergeSavedWord(words, { word: 'hello', count: 2, contexts: [] })[0].count).toBe(2);
  expect(words[0].count).toBe(1);
});
