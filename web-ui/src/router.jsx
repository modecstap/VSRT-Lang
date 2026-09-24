import { lazy, Suspense } from 'react';
import { createBrowserRouter, Navigate, Outlet } from 'react-router-dom';
import { hasAccessToken } from './api/auth';

const LoginPage = lazy(() => import('./pages/LoginPage/LoginPage'));
const ForgotPasswordPage = lazy(() => import('./pages/ForgotPasswordPage/ForgotPasswordPage'));
const ResetPasswordPage = lazy(() => import('./pages/ResetPasswordPage/ResetPasswordPage'));
const AccountPage = lazy(() => import('./pages/AccountPage/AccountPage'));
const ProfileTab = lazy(() => import('./pages/AccountPage/components/ProfileTab/ProfileTab'));
const SessionTab = lazy(() => import('./pages/AccountPage/components/SessionTab/SessionTab'));

function RouteFallback() {
  return <div>Loading...</div>;
}

function withSuspense(element) {
  return <Suspense fallback={<RouteFallback />}>{element}</Suspense>;
}

function GuestOnly() {
  if (hasAccessToken()) {
    return <Navigate to="/account" replace />;
  }

  return withSuspense(<LoginPage />);
}

function RequireAuth() {
  if (!hasAccessToken()) {
    return <Navigate to="/" replace />;
  }

  return <Outlet />;
}

export const router = createBrowserRouter([
  {
    path: '/',
    element: <GuestOnly />,
  },
  {
    path: '/forgot-password',
    element: withSuspense(<ForgotPasswordPage />),
  },
  {
    path: '/reset-password',
    element: withSuspense(<ResetPasswordPage />),
  },
  {
    path: '/account',
    element: <RequireAuth />,
    children: [
      {
        element: withSuspense(<AccountPage />),
        children: [
          {
            index: true,
            element: <Navigate to="profile" replace />,
          },
          {
            path: 'profile',
            element: withSuspense(<ProfileTab />),
          },
          {
            path: 'session',
            element: withSuspense(<SessionTab />),
          },
        ],
      },
    ],
  },
]);
