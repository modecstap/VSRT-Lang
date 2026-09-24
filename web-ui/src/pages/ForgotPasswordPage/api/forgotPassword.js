import { apiClient } from '../../../api/client';

export const requestPasswordReset = async ({ email }) => {
  await apiClient.post('/forgot-password', { email });
};
