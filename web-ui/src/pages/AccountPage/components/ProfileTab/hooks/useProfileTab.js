import { useEffect, useState } from 'react';
import { fetchUserSessions } from '../api/profileTabApi';
import { buildSessionsViewModel } from '../model/profileTabModel';

function useProfileTab() {
  const [sessions, setSessions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    let ignore = false;

    const loadSessions = async () => {
      try {
        const response = await fetchUserSessions();

        if (!ignore) {
          setSessions(buildSessionsViewModel(response));
          setError('');
        }
      } catch (err) {
        if (!ignore) {
          setSessions([]);
          setError('Unable to load sessions');
        }
      } finally {
        if (!ignore) {
          setLoading(false);
        }
      }
    };

    loadSessions();

    return () => {
      ignore = true;
    };
  }, []);

  return {
    sessions,
    loading,
    error,
  };
}

export default useProfileTab;
