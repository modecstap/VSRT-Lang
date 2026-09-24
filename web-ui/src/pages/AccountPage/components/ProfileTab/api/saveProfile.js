import { apiClient } from '../../../../../api/client';

export function saveProfile({ username, email }) {
  return apiClient.post('/users/me', { username, email });
}
