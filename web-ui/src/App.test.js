import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
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

test('returns the button label to Sign In after an error is shown for 5 seconds', async () => {
  jest.useFakeTimers();
  axios.post.mockRejectedValue({
    response: { data: { error: { message: 'Invalid credentials' } } },
  });

  render(<App />);

  fireEvent.change(screen.getByPlaceholderText(/login/i), { target: { value: 'alice' } });
  fireEvent.change(screen.getByPlaceholderText(/Password/i), { target: { value: 'secret' } });
  fireEvent.click(screen.getByRole('button', { name: /sign in/i }));

  await waitFor(() => {
    expect(screen.getByRole('button', { name: /invalid credentials/i })).toBeInTheDocument();
  });

  act(() => {
    jest.advanceTimersByTime(2500);
  });

  expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument();

  jest.useRealTimers();
});
