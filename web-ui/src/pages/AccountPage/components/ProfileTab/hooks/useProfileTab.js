import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import useSWR, { useSWRConfig } from 'swr';
import {
  createSession,
  deleteSession,
  fetchUserSessions,
  setActiveSessionId,
  SESSIONS_SWR_KEY,
  sessionRecordsKey,
} from '../api/profileTabApi';
import { buildSessionsViewModel } from '../model/profileTabModel';

function useProfileTab() {
  const navigate = useNavigate();
  const { mutate } = useSWRConfig();
  const {
    data: sessions = [],
    error: loadError,
    isLoading,
  } = useSWR(SESSIONS_SWR_KEY, async () => {
    const response = await fetchUserSessions();
    return buildSessionsViewModel(response);
  });
  const [creating, setCreating] = useState(false);
  const [deletingSessionId, setDeletingSessionId] = useState(null);
  const [sessionName, setSessionName] = useState('');
  const [actionError, setActionError] = useState('');

  const handleCreateSession = async () => {
    setCreating(true);
    setActionError('');

    try {
      const name = sessionName.trim() || 'Session';
      await createSession(name);
      await mutate(SESSIONS_SWR_KEY);
      navigate('/account/session');
    } catch (err) {
      setActionError('Unable to create new session');
    } finally {
      setCreating(false);
    }
  };

  const handleSessionNameChange = (event) => {
    setSessionName(event.target.value);
  };

  const handleSelectSession = (session) => {
    setActiveSessionId(session.id);
    navigate('/account/session');
  };

  const handleDeleteSession = async (session) => {
    setDeletingSessionId(session.id);
    setActionError('');

    try {
      await deleteSession(session.id);
      await mutate(
        SESSIONS_SWR_KEY,
        (currentSessions = []) => currentSessions.filter((item) => item.id !== session.id),
        { revalidate: false }
      );
      await mutate(sessionRecordsKey(session.id), undefined, { revalidate: false });
    } catch (err) {
      setActionError('Unable to delete session');
    } finally {
      setDeletingSessionId(null);
    }
  };

  return {
    sessions,
    loading: isLoading,
    error: actionError || (loadError ? 'Unable to load sessions' : ''),
    creating,
    deletingSessionId,
    sessionName,
    handleSessionNameChange,
    handleCreateSession,
    handleSelectSession,
    handleDeleteSession,
  };
}

export default useProfileTab;
