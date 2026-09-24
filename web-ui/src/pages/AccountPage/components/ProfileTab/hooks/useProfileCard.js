import { useEffect, useRef, useState } from 'react';
import useSWR from 'swr';
import { saveAvatar } from '../api/saveAvatar';
import { saveProfile } from '../api/saveProfile';
import { CURRENT_USER_SWR_KEY, fetchCurrentUser } from '../api/getCurrentUser';
import { avatarUploadError, profileCardView, profileSaveError } from '../model/profileCardModel';

const UPLOAD_ERROR_MS = 5000;
const RESULT_MS = 3000;

function draftField(value) {
  return typeof value === 'string' ? value : '';
}

function useProfileCard() {
  const { data, error, mutate } = useSWR(CURRENT_USER_SWR_KEY, fetchCurrentUser);
  const [uploading, setUploading] = useState(false);
  const [uploadError, setUploadError] = useState('');
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState({ username: '', email: '' });
  const [phase, setPhase] = useState('idle');
  const [saveError, setSaveError] = useState('');
  const errorTimer = useRef(null);

  useEffect(() => () => clearTimeout(errorTimer.current), []);

  useEffect(() => {
    if (phase !== 'success' && phase !== 'error') {
      return undefined;
    }

    const timeoutId = window.setTimeout(() => {
      setSaveError('');
      setPhase('idle');
    }, RESULT_MS);

    return () => window.clearTimeout(timeoutId);
  }, [phase]);

  const showUploadError = (message) => {
    clearTimeout(errorTimer.current);
    setUploadError(message);
    errorTimer.current = setTimeout(() => setUploadError(''), UPLOAD_ERROR_MS);
  };

  const handleAvatarChange = async (event) => {
    const file = event.target.files?.[0];
    event.target.value = '';
    if (!file || uploading) return;

    clearTimeout(errorTimer.current);
    setUploading(true);
    setUploadError('');

    try {
      await saveAvatar(file);
      await mutate();
    } catch (err) {
      showUploadError(avatarUploadError(err));
    } finally {
      setUploading(false);
    }
  };

  const onDraftChange = (event) => {
    const { name, value } = event.target;
    if (name !== 'username' && name !== 'email') return;
    setDraft((current) => ({ ...current, [name]: value }));
  };

  const onProfileButton = async () => {
    if (editing) {
      if (phase !== 'idle') return;
      setPhase('submitting');
      try {
        await saveProfile({ username: draft.username, email: draft.email });
        setEditing(false);
        setPhase('success');
        try {
          await mutate();
        } catch {
          // A failed refresh stays on the card. It is not a save error.
        }
      } catch (err) {
        setEditing(false);
        setSaveError(profileSaveError(err));
        setPhase('error');
      }
      return;
    }

    if (phase !== 'idle' || data == null) return;
    setDraft({
      username: draftField(data.username),
      email: draftField(data.email),
    });
    setEditing(true);
  };

  let buttonLabel = 'CHANGE';
  if (phase === 'submitting') {
    buttonLabel = 'Sending...';
  } else if (phase === 'success') {
    buttonLabel = 'SUCCESS';
  } else if (phase === 'error') {
    buttonLabel = saveError;
  } else if (editing) {
    buttonLabel = 'SAVE';
  }

  const buttonDisabled =
    data == null || phase === 'submitting' || phase === 'success' || phase === 'error';

  return {
    ...profileCardView({ data, error }),
    uploading,
    uploadError,
    handleAvatarChange,
    editing,
    draft,
    onDraftChange,
    buttonLabel,
    buttonDisabled,
    onProfileButton,
  };
}

export default useProfileCard;
