import { useEffect, useState } from 'react';
import { requestPasswordReset } from '../api/forgotPassword';
import { createInitialForgotPasswordForm } from '../model/forgotPasswordFormModel';

const ERROR_MS = 2500;
const SENT_MS = 3000;

export function useForgotPasswordForm() {
  const [form, setForm] = useState(createInitialForgotPasswordForm);
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
    if (phase !== 'sent') {
      return undefined;
    }

    const timeoutId = window.setTimeout(() => {
      setPhase('idle');
    }, SENT_MS);

    return () => window.clearTimeout(timeoutId);
  }, [phase]);

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
    setPhase('submitting');

    try {
      await requestPasswordReset(form);
      setPhase('sent');
    } catch (submitError) {
      setError(submitError.response?.data?.error?.message || submitError.message || 'Unable to send email');
      setPhase('error');
    }
  };

  let buttonLabel = 'Send';
  if (phase === 'submitting') {
    buttonLabel = 'Sending...';
  } else if (phase === 'error') {
    buttonLabel = error;
  } else if (phase === 'sent') {
    buttonLabel = 'Sent';
  }

  return {
    form,
    phase,
    buttonLabel,
    handleChange,
    handleSubmit,
  };
}
