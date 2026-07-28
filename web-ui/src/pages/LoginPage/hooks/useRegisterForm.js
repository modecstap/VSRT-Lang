import { useState } from 'react';

const createInitialRegisterForm = () => ({
  username: '',
  email: '',
  password: '',
});

export function useRegisterForm({ onSuccess } = {}) {
  const [form, setForm] = useState(createInitialRegisterForm());
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
      const response = await fetch('http://localhost:8080/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form),
      });

      if (!response.ok) {
        throw new Error('Registration failed');
      }

      setForm(createInitialRegisterForm());
      onSuccess?.();
    } catch (error) {
      setMessage(error.message || 'Registration failed');
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
