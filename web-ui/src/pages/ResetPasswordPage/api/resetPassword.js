import { apiClient } from '../../../api/client';

export const resetPassword = async ({ token, password }) => {
  await apiClient.post('/reset-password', { token, password });
};
