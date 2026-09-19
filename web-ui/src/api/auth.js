import {
  clearAuthStorage,
  readAuthStorage,
  writeAuthStorage,
} from './storage';

export function getAccessToken() {
  return readAuthStorage()?.accessToken || '';
}

export function hasAccessToken() {
  return Boolean(getAccessToken());
}

export function saveAuthTokens(tokens = {}) {
  if (!tokens.AccessToken) {
    return;
  }

  writeAuthStorage({
    accessToken: tokens.AccessToken,
    refreshToken: tokens.RefreshToken || '',
  });
}

export function clearAuthTokens() {
  clearAuthStorage();
}
