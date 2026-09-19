import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { saveAuthTokens } from '../../../api/auth';
import { prefetchAccount } from '../../../routes/prefetch';
import { login } from '../api/login';
import { createInitialLoginForm } from '../model/loginFormModel';

export function useLoginForm() {
  const navigate = useNavigate();
  const [form, setForm] = useState(createInitialLoginForm);
  const [error, setError] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (!error) {
      return undefined;
    }

    const timeoutId = window.setTimeout(() => {
      setError('');
    }, 2500);

    return () => window.clearTimeout(timeoutId);
  }, [error]);

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

      if (data?.AccessToken) {
        saveAuthTokens(data);
        prefetchAccount();
        navigate('/account');
      }
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
