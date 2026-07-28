import styles from '../LoginForm/LoginForm.module.css';
import Button from '../../../../shared/ui/Button/Button';
import FormField from '../LoginForm/FormField';
import { useState } from 'react';

function RegisterForm({ onSwitchToLogin }) {
  const [form, setForm] = useState({
    username: '',
    email: '',
    password: '',
  });
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

      setMessage('Registration successful. You can sign in now.');
      setForm({ username: '', email: '', password: '' });
    } catch (error) {
      setMessage(error.message || 'Registration failed');
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form className={styles.form} onSubmit={handleSubmit}>
      <div className={styles.grid}>
        <FormField
          label=" "
          id="register-username"
          name="username"
          type="text"
          placeholder="Username"
          value={form.username}
          onChange={handleChange}
        />

        <FormField
          label=" "
          id="register-email"
          name="email"
          type="email"
          placeholder="Email"
          value={form.email}
          onChange={handleChange}
        />

        <FormField
          label=" "
          id="register-password"
          name="password"
          type="password"
          placeholder="Password"
          value={form.password}
          onChange={handleChange}
          autoComplete="new-password"
        />
      </div>

      <div className={styles.actions}>
        <Button type="submit" disabled={isSubmitting}>
          {isSubmitting ? 'Creating account...' : 'Create account'}
        </Button>
        <div>
            <button type="button" className={styles.link} onClick={onSwitchToLogin}>
            BACK TO LOGIN
            </button>
        </div>
      </div>

      {message ? <p className={styles.link}>{message}</p> : null}
    </form>
  );
}

export default RegisterForm;
