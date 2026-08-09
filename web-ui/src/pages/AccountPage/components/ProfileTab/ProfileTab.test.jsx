import { render, screen } from '@testing-library/react';
import axios from 'axios';
import ProfileTab from './ProfileTab';

jest.mock('axios');

describe('ProfileTab', () => {
  beforeEach(() => {
    localStorage.clear();
    axios.get.mockReset();
  });

  it('loads sessions from the server', async () => {
    axios.get.mockResolvedValueOnce({
      data: [
        { id: 7, name: 'Server session' },
        { id: 8, name: 'Another session' },
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
});
