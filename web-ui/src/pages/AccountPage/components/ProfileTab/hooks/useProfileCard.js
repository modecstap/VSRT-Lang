import useSWR from 'swr';
import { CURRENT_USER_SWR_KEY, fetchCurrentUser } from '../api/getCurrentUser';
import { profileCardView } from '../model/profileCardModel';

function useProfileCard() {
  const { data, error } = useSWR(CURRENT_USER_SWR_KEY, fetchCurrentUser);
  return profileCardView({ data, error });
}

export default useProfileCard;
