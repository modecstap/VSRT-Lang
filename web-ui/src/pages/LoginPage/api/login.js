import { apiClient } from '../../../api/client';

export const login = async (credentials) => {
  const payload = {
    login: credentials.login,
    password: credentials.password,
  };

  const response = await apiClient.post('/login', payload);

  return response.data;
};
