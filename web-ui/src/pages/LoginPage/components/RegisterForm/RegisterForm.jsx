import styles from '../LoginForm/LoginForm.module.css';
import Button from '../../../../shared/ui/Button/Button';
import FormField from '../LoginForm/FormField';
import { useRegisterForm } from '../../hooks/useRegisterForm';

function RegisterForm({ onSwitchToLogin }) {
  const { form, message, isSubmitting, handleChange, handleSubmit } = useRegisterForm({
    onSuccess: onSwitchToLogin,
  });

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
