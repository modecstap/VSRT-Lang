import { useState } from 'react';
import { login } from '../api/login';
import {
  buildProfilePath,
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
      const username = data?.username ?? form.username;

      window.location.assign(buildProfilePath(username));
    } catch (submitError) {
      setError(submitError.message || 'Unable to sign in');
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
