import axios from 'axios';
import { getAuthHeaders } from '../../../../../api/auth';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export async function fetchUserSessions() {
  const response = await axios.get(
    `${API_BASE_URL}/users/sessions`,
    getAuthHeaders()
  );

  return response.data || [];
}
