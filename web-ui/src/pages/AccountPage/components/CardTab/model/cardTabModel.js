export const ESTIMATIONS = ['MOMENTAL', 'EASY', 'DIFFICULT', 'REPEAT'];

export function cardsKey(sessionId, visitId) {
  return ['session-cards', sessionId, visitId];
}

export function mapCard(raw) {
  return {
    phrase: raw?.front?.phrase || '',
    baseForm: raw?.front?.base_form || '',
    contexts: raw?.front?.contexts || [],
    synonyms: raw?.front?.synonyms || [],
    translations: raw?.back?.translations || [],
    contextTranslations: raw?.back?.context_translations || [],
  };
}

export function createRun(cards = []) {
  return { queue: cards, passed: 0, revealed: false, estimation: null };
}

export function revealCard(run) {
  return { ...run, revealed: true };
}

export function estimateCard(run, estimation) {
  return { ...run, estimation };
}

export function advanceRun(run) {
  if (run.queue.length === 0) return run;
  const [current, ...rest] = run.queue;
  const passed = Boolean(run.estimation) && run.estimation !== 'REPEAT';
  return {
    queue: passed ? rest : [...rest, current],
    passed: run.passed + (passed ? 1 : 0),
    revealed: false,
    estimation: null,
  };
}

const joinInline = (items) => items.join(', ') || '—';

const joinLines = (items) => items.map((item) => item || '—').join('\n') || '—';

export function buildCardView(card) {
  if (!card) {
    return null;
  }

  return {
    phrase: card.phrase || '—',
    baseForm: card.baseForm || '—',
    synonyms: joinInline(card.synonyms),
    translations: joinInline(card.translations),
    contexts: joinLines(card.contexts),
    contextTranslations: joinLines(card.contextTranslations),
  };
}
