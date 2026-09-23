const PLACEHOLDER = 'https://i.pravatar.cc/150?img=12';

const loadingCard = {
  username: '...',
  email: '...',
  avatar: PLACEHOLDER,
};

const loadErrorCard = {
  username: 'Unable to load profile',
  email: '',
  avatar: PLACEHOLDER,
};

function stringField(value) {
  return typeof value === 'string' ? value : '';
}

function avatarSrc(avatar) {
  if (typeof avatar === 'string' && avatar.length > 0) {
    return avatar;
  }
  return PLACEHOLDER;
}

function cardFromUser(data) {
  return {
    username: stringField(data.username),
    email: stringField(data.email),
    avatar: avatarSrc(data.avatar),
  };
}

function withRefreshMark(username) {
  const mark = '(Unable to refresh profile)';
  if (username.length > 0) {
    return username + ' ' + mark;
  }
  return mark;
}

export function profileCardView({ data, error }) {
  if (data == null) {
    return error ? loadErrorCard : loadingCard;
  }
  const card = cardFromUser(data);
  if (!error) return card;
  return { ...card, username: withRefreshMark(card.username) };
}
