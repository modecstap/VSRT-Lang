import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import SessionTab from './SessionTab';

describe('SessionTab', () => {
  it('adds a new saved word and shows its details after writing', async () => {
    render(<SessionTab />);

    const wordInput = screen.getByLabelText(/word/i);
    const contextInput = screen.getByLabelText(/context/i);

    await userEvent.type(wordInput, 'sunrise');
    await userEvent.type(contextInput, 'The sunrise was beautiful.');
    await userEvent.click(screen.getByRole('button', { name: /write/i }));

    await waitFor(() => {
      expect(screen.getByText('sunrise')).toBeInTheDocument();
    });

    expect(screen.getByText(/the sunrise was beautiful/i)).toBeInTheDocument();
  });
});
