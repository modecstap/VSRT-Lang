import { removeSavedWord } from './sessionTabModel';

test('removeSavedWord drops the matching phrase and leaves the input unchanged', () => {
  const words = [{ word: 'Hello' }, { word: 'world' }];

  expect(removeSavedWord(words, 'hello')).toEqual([{ word: 'world' }]);
  expect(words).toEqual([{ word: 'Hello' }, { word: 'world' }]);
});
