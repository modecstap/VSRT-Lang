import { apiClient } from '../../../../../api/client';

export async function fetchSessionCards(sessionId) {
  const response = await apiClient.get(`/sessions/${sessionId}/cards`);
  return response.data?.cards || [];
}

export async function saveEstimation(sessionId, phrase, estimation) {
  await apiClient.post(`/sessions/${sessionId}/cards`, { phrase, estimation });
}
