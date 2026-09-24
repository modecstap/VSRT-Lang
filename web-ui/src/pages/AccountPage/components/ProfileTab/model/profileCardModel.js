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

export function profileSaveError(error) {
  const code = error?.response?.data?.error?.code;
  switch (code) {
    case 'username_required':
      return 'Username required';
    case 'email_required':
      return 'Email required';
    case 'invalid_email':
      return 'Invalid email';
    case 'username_taken':
      return 'Username taken';
    case 'email_taken':
      return 'Email taken';
    default:
      return 'Unable to save';
  }
}

export function avatarUploadError(error) {
  const code = error?.response?.data?.error?.code;
  switch (code) {
    case 'avatar_too_large':
      return 'Image is too large';
    case 'invalid_avatar_type':
      return 'Use a JPEG, PNG, or WebP image';
    case 'avatar_dimensions_invalid':
      return 'Image is too wide or too tall';
    default:
      return 'Unable to save avatar';
  }
}

export function profileCardView({ data, error }) {
  if (data == null) {
    return error ? loadErrorCard : loadingCard;
  }
  const card = cardFromUser(data);
  if (!error) return card;
  return { ...card, username: withRefreshMark(card.username) };
}
