import { useState } from 'react';
import { login } from '../api/login';
import {
  createInitialLoginForm,
} from '../model/loginFormModel';

export function useLoginForm() {
  const [form, setForm] = useState(createInitialLoginForm());
  const [error, setError] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

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
    setIsSubmitting(true);

    try {
      const data = await login(form);

      if (data?.access_token) {
        localStorage.setItem('access_token', data.access_token);
        localStorage.setItem('refresh_token', data.refresh_token || '');
      }

      window.location.assign('/');
    } catch (submitError) {
      setError(submitError.response?.data?.error?.message || submitError.message || 'Unable to sign in');
    } finally {
      setIsSubmitting(false);
    }
  };

  return {
    form,
    error,
    isSubmitting,
    handleChange,
    handleSubmit,
  };
}
