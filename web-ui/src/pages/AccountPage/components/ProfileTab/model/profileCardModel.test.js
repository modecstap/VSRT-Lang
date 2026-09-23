import { profileCardView } from './profileCardModel';

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
