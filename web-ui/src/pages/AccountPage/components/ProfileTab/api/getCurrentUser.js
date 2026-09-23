import { apiClient } from '../../../../../api/client';

export const CURRENT_USER_SWR_KEY = ['current-user'];

export async function fetchCurrentUser() {
  const response = await apiClient.get('/users/me');
  return response.data;
}
