import { avatarUploadError, profileCardView, profileSaveError } from './profileCardModel';

const placeholder = 'https://i.pravatar.cc/150?img=12';
const body = {
  username: 'Modecstap',
  email: 'a@b.c',
  avatar: 'data:image/png;base64,qq',
};
const marked = 'Modecstap (Unable to refresh profile)';

test('profileCardView covers load, error, success, and refresh failure', () => {
  expect(profileCardView({ data: undefined, error: undefined })).toEqual({
    username: '...',
    email: '...',
    avatar: placeholder,
  });

  expect(profileCardView({ data: undefined, error: new Error('down') })).toEqual({
    username: 'Unable to load profile',
    email: '',
    avatar: placeholder,
  });

  expect(profileCardView({ data: body, error: undefined })).toEqual({
    username: 'Modecstap',
    email: 'a@b.c',
    avatar: 'data:image/png;base64,qq',
  });

  expect(profileCardView({ data: { ...body, avatar: null }, error: undefined })).toEqual({
    username: 'Modecstap',
    email: 'a@b.c',
    avatar: placeholder,
  });

  expect(profileCardView({ data: { ...body, avatar: '' }, error: undefined })).toEqual({
    username: 'Modecstap',
    email: 'a@b.c',
    avatar: placeholder,
  });

  expect(profileCardView({ data: body, error: new Error('down') })).toEqual({
    username: marked,
    email: 'a@b.c',
    avatar: 'data:image/png;base64,qq',
  });

  expect(profileCardView({ data: { ...body, avatar: null }, error: new Error('down') })).toEqual({
    username: marked,
    email: 'a@b.c',
    avatar: placeholder,
  });
});

function uploadError(code) {
  return { response: { data: { error: { code } } } };
}

test('profileSaveError maps server codes and falls back', () => {
  expect(profileSaveError(uploadError('username_required'))).toBe('Username required');
  expect(profileSaveError(uploadError('email_required'))).toBe('Email required');
  expect(profileSaveError(uploadError('invalid_email'))).toBe('Invalid email');
  expect(profileSaveError(uploadError('username_taken'))).toBe('Username taken');
  expect(profileSaveError(uploadError('email_taken'))).toBe('Email taken');
  expect(profileSaveError(uploadError('profile_save_failed'))).toBe('Unable to save');
  expect(profileSaveError(new Error('down'))).toBe('Unable to save');
});

test('avatarUploadError maps server codes and falls back', () => {
  expect(avatarUploadError(uploadError('avatar_too_large'))).toBe('Image is too large');
  expect(avatarUploadError(uploadError('invalid_avatar_type'))).toBe(
    'Use a JPEG, PNG, or WebP image'
  );
  expect(avatarUploadError(uploadError('avatar_dimensions_invalid'))).toBe(
    'Image is too wide or too tall'
  );
  expect(avatarUploadError(uploadError('avatar_save_failed'))).toBe('Unable to save avatar');
  expect(avatarUploadError(new Error('down'))).toBe('Unable to save avatar');
});
