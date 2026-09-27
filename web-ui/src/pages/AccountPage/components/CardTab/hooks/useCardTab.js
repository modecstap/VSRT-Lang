import { useState } from 'react';
import useSWR from 'swr';
import { getActiveSessionId } from '../../../../../api/sessions';
import { fetchSessionCards, saveEstimation } from '../api/cardTabApi';
import {
  advanceRun,
  buildCardView,
  cardsKey,
  createRun,
  estimateCard,
  mapCard,
  revealCard,
} from '../model/cardTabModel';

function useCardTab() {
  const sessionId = getActiveSessionId();
  const [visitId] = useState(() => Date.now());
  const { data: cards, error } = useSWR(
    sessionId ? cardsKey(sessionId, visitId) : null,
    async () => (await fetchSessionCards(sessionId)).map(mapCard),
    { revalidateOnFocus: false, revalidateOnReconnect: false, revalidateIfStale: false }
  );
  const [run, setRun] = useState(null);
  const [isEstimating, setIsEstimating] = useState(false);
  const [estimateError, setEstimateError] = useState('');

  const current = run ?? (cards ? createRun(cards) : null);

  let status = '';
  if (!sessionId) {
    status = 'Select a session in PROFILE.';
  } else if (error) {
    status = 'Unable to load cards';
  } else if (!current) {
    status = 'Loading cards...';
  } else if (current.queue.length === 0) {
    status = 'No cards to repeat.';
  }

  const card = status ? null : buildCardView(current.queue[0]);
  const revealed = current?.revealed ?? false;
  const estimation = current?.estimation ?? null;
  const canEstimate = Boolean(card) && revealed && !estimation && !isEstimating;

  const handleEstimate = async (nextEstimation) => {
    if (!canEstimate) {
      return;
    }

    setIsEstimating(true);
    setEstimateError('');

    try {
      await saveEstimation(sessionId, current.queue[0].phrase, nextEstimation);
      setRun(estimateCard(current, nextEstimation));
    } catch (submitError) {
      setEstimateError('Unable to save estimation');
    } finally {
      setIsEstimating(false);
    }
  };

  const handleMain = () => {
    if (!card || isEstimating) {
      return;
    }

    if (!revealed) {
      setRun(revealCard(current));
      return;
    }

    setRun(advanceRun(current));
    setEstimateError('');
  };

  return {
    canEstimate,
    card,
    estimation,
    handleEstimate,
    handleMain,
    mainDisabled: !card || isEstimating,
    mainLabel: estimateError || (isEstimating ? '...' : revealed ? 'next' : 'open'),
    passed: current?.passed ?? 0,
    remaining: current?.queue.length ?? 0,
    revealed,
    status,
  };
}

export default useCardTab;
