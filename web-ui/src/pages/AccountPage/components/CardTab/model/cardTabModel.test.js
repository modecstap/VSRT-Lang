import { advanceRun, buildCardView, mapCard } from './cardTabModel';

const run = (estimation) => ({ queue: ['a', 'b'], passed: 1, revealed: true, estimation });

test('advanceRun without estimation moves the card to the end', () => {
  expect(advanceRun(run(null))).toEqual({ queue: ['b', 'a'], passed: 1, revealed: false, estimation: null });
});

test('advanceRun after REPEAT moves the card to the end', () => {
  expect(advanceRun(run('REPEAT'))).toEqual({ queue: ['b', 'a'], passed: 1, revealed: false, estimation: null });
});

test('advanceRun after EASY drops the card and leaves the input unchanged', () => {
  const input = run('EASY');

  expect(advanceRun(input)).toEqual({ queue: ['b'], passed: 2, revealed: false, estimation: null });
  expect(input).toEqual(run('EASY'));
});

test('advanceRun keeps an empty queue as is', () => {
  const input = { queue: [], passed: 3, revealed: false, estimation: null };

  expect(advanceRun(input)).toBe(input);
});

test('buildCardView shows a dash for empty fields', () => {
  const view = buildCardView(mapCard({ front: { phrase: 'go', contexts: ['go home', ''] } }));

  expect(view.baseForm).toBe('—');
  expect(view.synonyms).toBe('—');
  expect(view.contextTranslations).toBe('—');
  expect(view.contexts).toBe('go home\n—');
});
