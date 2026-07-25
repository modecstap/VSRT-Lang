export const createInitialLoginForm = () => ({
  username: '',
  password: '',
});

export const buildProfilePath = (username) => `/profile/${username}`;
