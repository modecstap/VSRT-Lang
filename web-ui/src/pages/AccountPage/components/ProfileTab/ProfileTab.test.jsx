import { render, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import ProfileTab from './ProfileTab';

const mockNavigate = jest.fn();

jest.mock('axios');
jest.mock('react-router-dom', () => ({
  useNavigate: () => mockNavigate,
}));

describe('ProfileTab', () => {
  beforeEach(() => {
    localStorage.clear();
    axios.get.mockReset();
    axios.post.mockReset();
    axios.delete.mockReset();
    mockNavigate.mockClear();
  });

  it('loads sessions from the server', async () => {
    axios.get.mockResolvedValueOnce({
      data: [
        { ID: 7, CreatedAt: '2024-01-01', User: 'Server session', Records: [] },
        { ID: 8, CreatedAt: '2024-01-02', User: 'Another session', Records: [] },
      ],
    });

    render(<ProfileTab />);

    expect(await screen.findByText('Server session')).toBeInTheDocument();
    expect(screen.getByText('Another session')).toBeInTheDocument();
    expect(axios.get).toHaveBeenCalledWith(
      expect.stringContaining('/users/sessions'),
      expect.anything()
    );
  });

  it('creates a new session and navigates to session tab', async () => {
    axios.get.mockResolvedValueOnce({
      data: [
        { ID: 7, CreatedAt: '2024-01-01', User: 'Server session', Records: [] },
      ],
    });

    axios.post.mockResolvedValueOnce({
      data: { id: 99 },
    });

    render(<ProfileTab />);

    expect(await screen.findByText('Server session')).toBeInTheDocument();

    await screen.findByRole('button', { name: /new/i });
    await screen.getByRole('button', { name: /new/i }).click();

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith('/account/session');
    });

    expect(axios.post).toHaveBeenCalledWith(
      expect.stringContaining('/sessions'),
      { name: 'Session' },
      expect.anything()
    );
    expect(localStorage.getItem('session_tab_session_id')).toBe('99');
  });

  it('deletes a session without selecting its row', async () => {
    axios.get.mockResolvedValueOnce({
      data: [
        { ID: 7, CreatedAt: '2024-01-01', Name: 'Server session', Records: [] },
        { ID: 8, CreatedAt: '2024-01-02', Name: 'Another session', Records: [] },
      ],
    });
    axios.delete.mockResolvedValueOnce({});

    render(<ProfileTab />);

    expect(await screen.findByText('Server session')).toBeInTheDocument();
    await screen.getAllByRole('button', { name: /delete/i })[0].click();

    await waitFor(() => {
      expect(axios.delete).toHaveBeenCalledWith(
        expect.stringContaining('/sessions/7'),
        expect.anything()
      );
    });
    expect(screen.queryByText('Server session')).not.toBeInTheDocument();
    expect(screen.getByText('Another session')).toBeInTheDocument();
    expect(mockNavigate).not.toHaveBeenCalled();
  });
});
