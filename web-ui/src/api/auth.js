import axios from 'axios';

const ACCESS_TOKEN_KEY = 'access_token';
const REFRESH_TOKEN_KEY = 'refresh_token';

axios.interceptors.request.use((config) => {
  const token = localStorage.getItem(ACCESS_TOKEN_KEY);

  if (token) {
    config.headers = config.headers || {};
    config.headers.Authorization = `Bearer ${token}`;
  }

  return config;
}, (error) => Promise.reject(error));

export function saveAuthTokens(tokens = {}) {
  if (tokens.AccessToken) {
    localStorage.setItem(ACCESS_TOKEN_KEY, tokens.AccessToken);
  }

  if (tokens.RefreshToken !== undefined) {
    localStorage.setItem(REFRESH_TOKEN_KEY, tokens.RefreshToken || '');
  }
}

export function clearAuthTokens() {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
}

export function getAuthHeaders() {
  const token = localStorage.getItem(ACCESS_TOKEN_KEY);

  return token
    ? { headers: { Authorization: `Bearer ${token}` } }
    : {};
}
