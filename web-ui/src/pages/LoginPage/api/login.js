import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export const login = async (credentials) => {
  const payload = {
    login: credentials.login,
    password: credentials.password,
  };

  const response = await axios.post(`${API_BASE_URL}/login`, payload);

  return response.data;
};
