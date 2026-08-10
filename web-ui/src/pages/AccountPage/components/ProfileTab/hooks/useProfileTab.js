import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { fetchUserSessions } from '../api/profileTabApi';
import { createSession } from '../../SessionTab/api/sessionTabApi';
import { buildSessionsViewModel } from '../model/profileTabModel';

function useProfileTab() {
  const navigate = useNavigate();
  const [sessions, setSessions] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [creating, setCreating] = useState(false);
  const [sessionName, setSessionName] = useState('');

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

  const handleCreateSession = async () => {
    setCreating(true);
    setError('');

    try {
      const name = sessionName.trim() || 'Session';
      await createSession(name);
      navigate('/account/session');
    } catch (err) {
      setError('Unable to create new session');
    } finally {
      setCreating(false);
    }
  };

  const handleSessionNameChange = (event) => {
    setSessionName(event.target.value);
  };

  return {
    sessions,
    loading,
    error,
    creating,
    sessionName,
    handleSessionNameChange,
    handleCreateSession,
  };
}

export default useProfileTab;
