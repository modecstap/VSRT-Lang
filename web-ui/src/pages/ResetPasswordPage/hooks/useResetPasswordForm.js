import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useSWRConfig } from 'swr';
import { clearAuthTokens } from '../../../api/auth';
import { clearActiveSessionId } from '../../../api/sessions';
import { resetPassword } from '../api/resetPassword';
import { createInitialResetPasswordForm } from '../model/resetPasswordFormModel';

const ERROR_MS = 2500;
const SUCCESS_MS = 3000;

export function useResetPasswordForm() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { mutate } = useSWRConfig();
  const token = searchParams.get('token') || '';

  const [form, setForm] = useState(createInitialResetPasswordForm);
  const [phase, setPhase] = useState('idle');
  const [error, setError] = useState('');

  useEffect(() => {
    if (phase !== 'error') {
      return undefined;
    }

    const timeoutId = window.setTimeout(() => {
      setError('');
      setPhase('idle');
    }, ERROR_MS);

    return () => window.clearTimeout(timeoutId);
  }, [phase]);

  useEffect(() => {
    if (phase !== 'success') {
      return undefined;
    }

    const timeoutId = window.setTimeout(() => {
      navigate('/');
    }, SUCCESS_MS);

    return () => window.clearTimeout(timeoutId);
  }, [phase, navigate]);

  const handleChange = (event) => {
    const { name, value } = event.target;
    setForm((current) => ({
      ...current,
      [name]: value,
    }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError('');

    if (form.password !== form.repeatPassword) {
      setError('Passwords do not match');
      setPhase('error');
      return;
    }

    const weakMessage =
      'password must be at least 8 characters and include uppercase, lowercase and digit';
    if (
      form.password.length < 8 ||
      !/[A-Z]/.test(form.password) ||
      !/[a-z]/.test(form.password) ||
      !/[0-9]/.test(form.password)
    ) {
      setError(weakMessage);
      setPhase('error');
      return;
    }

    if (!token) {
      setError('This link is invalid or expired');
      setPhase('error');
      return;
    }

    setPhase('submitting');

    try {
      await resetPassword({ token, password: form.password });
      clearAuthTokens();
      clearActiveSessionId();
      mutate(() => true, undefined, { revalidate: false });
      setPhase('success');
    } catch (submitError) {
      setError(submitError.response?.data?.error?.message || submitError.message || 'Unable to reset password');
      setPhase('error');
    }
  };

  let buttonLabel = 'Update';
  if (phase === 'submitting') {
    buttonLabel = 'Updating...';
  } else if (phase === 'error') {
    buttonLabel = error;
  } else if (phase === 'success') {
    buttonLabel = 'Success';
  }

  return {
    form,
    phase,
    buttonLabel,
    handleChange,
    handleSubmit,
  };
}
