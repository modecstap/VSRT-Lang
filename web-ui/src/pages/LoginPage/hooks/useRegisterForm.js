import { useState } from 'react';
import { registerUser } from '../api/register';
import { createInitialRegisterForm } from '../model/registerFormModel';


export function useRegisterForm({ onSuccess } = {}) {
  const [form, setForm] = useState(createInitialRegisterForm);
  const [message, setMessage] = useState('');
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
    setMessage('');
    setIsSubmitting(true);

    try {
      await registerUser(form);

      setForm(createInitialRegisterForm());
      onSuccess?.();
    } catch (error) {
      setMessage(error.response?.data?.error?.message || error.message || 'Registration failed');
    } finally {
      setIsSubmitting(false);
    }
  };

  return {
    form,
    message,
    isSubmitting,
    handleChange,
    handleSubmit,
  };
}

export default useRegisterForm;
