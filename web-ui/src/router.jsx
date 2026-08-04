import { createBrowserRouter, Navigate } from 'react-router-dom';

import LoginPage from './pages/LoginPage/LoginPage';
import AccountPage from './pages/AccountPage/AccountPage';
import ProfileTab from './pages/AccountPage/components/ProfileTab/ProfileTab';
import SessionTab from './pages/AccountPage/components/SessionTab/SessionTab';

export const router = createBrowserRouter([
  {
    path: '/',
    element: <LoginPage />,
  },
  {
    path: '/account',
    element: <AccountPage />,
    children: [
      {
        index: true,
        element: <Navigate to="profile" replace />,
      },
      {
        path: 'profile',
        element: <ProfileTab />,
      },
      {
        path: 'session',
        element: <SessionTab />,
      },
    ],
  },
]);