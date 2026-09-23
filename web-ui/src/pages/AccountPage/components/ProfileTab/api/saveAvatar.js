import { apiClient } from '../../../../../api/client';

export async function saveAvatar(file) {
  const body = new FormData();
  body.append('avatar', file);
  await apiClient.post('/users/avatar', body);
}
