import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import axios from 'axios';
import App from './App';

jest.mock('axios');

test('submits login form to the backend endpoint', async () => {
  axios.post.mockResolvedValue({
    data: { access_token: 'token', refresh_token: 'refresh-token' },
  });

  render(<App />);

  await userEvent.type(screen.getByPlaceholderText(/login/i), 'alice');
  await userEvent.type(screen.getByPlaceholderText(/Password/i), 'secret');
  await userEvent.click(screen.getByRole('button', { name: /sign in/i }));

  expect(axios.post).toHaveBeenCalledWith('http://localhost:8080/login', {
    login: 'alice',
    password: 'secret',
  });
});
