import { useEffect, useRef, useState } from 'react';
import useSWR from 'swr';
import { saveAvatar } from '../api/saveAvatar';
import { CURRENT_USER_SWR_KEY, fetchCurrentUser } from '../api/getCurrentUser';
import { avatarUploadError, profileCardView } from '../model/profileCardModel';

const UPLOAD_ERROR_MS = 5000;

function useProfileCard() {
  const { data, error, mutate } = useSWR(CURRENT_USER_SWR_KEY, fetchCurrentUser);
  const [uploading, setUploading] = useState(false);
  const [uploadError, setUploadError] = useState('');
  const errorTimer = useRef(null);

  useEffect(() => () => clearTimeout(errorTimer.current), []);

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

  return {
    ...profileCardView({ data, error }),
    uploading,
    uploadError,
    handleAvatarChange,
  };
}

export default useProfileCard;
