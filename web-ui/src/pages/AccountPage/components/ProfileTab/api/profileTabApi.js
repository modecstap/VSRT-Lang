import axios from 'axios';
import { getAuthHeaders } from '../../../../../api/auth';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';
export const SESSION_STORAGE_KEY = 'session_tab_session_id';

export async function fetchUserSessions() {
  const response = await axios.get(
    `${API_BASE_URL}/users/sessions`,
    getAuthHeaders()
  );

  return response.data || [];
}

export async function deleteSession(sessionId) {
  await axios.delete(
    `${API_BASE_URL}/sessions/${sessionId}`,
    getAuthHeaders()
  );
}
