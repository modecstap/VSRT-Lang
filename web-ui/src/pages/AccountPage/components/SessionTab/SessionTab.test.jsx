import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import axios from 'axios';
import SessionTab from './SessionTab';

jest.mock('axios');

describe('SessionTab', () => {
  beforeEach(() => {
    localStorage.clear();
    axios.post.mockReset();
    axios.get.mockReset();
  });

  it('loads records from the server and saves a new one', async () => {
    axios.post
      .mockResolvedValueOnce({ data: { id: 1, name: 'Session' } })
      .mockResolvedValueOnce({ data: { message: 'record saved' } });

    axios.get
      .mockResolvedValueOnce({
        data: {
          records: [
            {
              phrase: 'server word',
              translations: ['слово'],
              synonyms: [],
              antonyms: [],
              baseForm: 'server word',
              contexts: [{ phrase: 'The server context', translation: 'Серверный контекст' }],
            },
          ],
        },
      })
      .mockResolvedValueOnce({
        data: {
          records: [
            {
              phrase: 'server word',
              translations: ['слово'],
              synonyms: [],
              antonyms: [],
              baseForm: 'server word',
              contexts: [{ phrase: 'The server context', translation: 'Серверный контекст' }],
            },
            {
              phrase: 'new phrase',
              translations: ['новое слово'],
              synonyms: [],
              antonyms: [],
              baseForm: 'new phrase',
              contexts: [{ phrase: 'The new context', translation: 'Новый контекст' }],
            },
          ],
        },
      });

    render(<SessionTab />);

    expect(await screen.findByText('server word')).toBeInTheDocument();

    const wordInput = screen.getByLabelText(/word/i);
    const contextInput = screen.getByLabelText(/context/i);

    await userEvent.type(wordInput, 'new phrase');
    await userEvent.type(contextInput, 'The new context');
    await userEvent.click(screen.getByRole('button', { name: /write/i }));

    await waitFor(() => {
      expect(screen.getByText('new phrase')).toBeInTheDocument();
    });

    expect(axios.post).toHaveBeenCalledWith(
      expect.stringContaining('/sessions'),
      expect.anything(),
      expect.anything()
    );
    expect(axios.post).toHaveBeenCalledWith(
      expect.stringContaining('/sessions/1/records'),
      expect.objectContaining({ phrase: 'new phrase', context: 'The new context' }),
      expect.anything()
    );
  });
});
